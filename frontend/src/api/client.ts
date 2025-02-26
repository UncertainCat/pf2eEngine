import {
  GameEvent,
  GameState,
  CommandRequest,
  CommandResponse,
  EventType
} from './types';

// API Client configuration
interface ApiClientConfig {
  baseUrl: string;
  apiVersion: string;
  wsEndpoint: string;
}

// Default configuration
const defaultConfig: ApiClientConfig = {
  baseUrl: 'http://localhost:8080',
  apiVersion: 'v1',
  wsEndpoint: 'ws://localhost:8080/ws'
};

// Event handler type
type EventHandler<T = any> = (event: GameEvent<T>) => void;

/**
 * API Client for PF2E Engine
 * Handles HTTP and WebSocket communication with the backend
 */
export class ApiClient {
  private config: ApiClientConfig;
  private websocket: WebSocket | null = null;
  private eventHandlers: Map<string, EventHandler[]> = new Map();
  private reconnectTimer: NodeJS.Timeout | null = null;
  private isConnecting: boolean = false;

  constructor(config: Partial<ApiClientConfig> = {}) {
    this.config = { ...defaultConfig, ...config };
  }

  /**
   * Connect to the WebSocket server
   */
  public connect(): Promise<void> {
    if (this.websocket && this.websocket.readyState === WebSocket.OPEN) {
      return Promise.resolve();
    }

    if (this.isConnecting) {
      return new Promise((resolve) => {
        // Add a one-time handler for the connect event
        const handler = () => {
          this.off('connect', handler);
          resolve();
        };
        this.on('connect', handler);
      });
    }

    this.isConnecting = true;
    return new Promise((resolve, reject) => {
      try {
        this.websocket = new WebSocket(this.config.wsEndpoint);

        this.websocket.onopen = () => {
          console.log('WebSocket connected');
          this.isConnecting = false;
          this.triggerEvent('connect', { type: 'connect', version: this.config.apiVersion, timestamp: new Date().toISOString() });
          resolve();
        };

        this.websocket.onmessage = (event) => {
          try {
            const data = JSON.parse(event.data);
            this.handleEvent(data);
          } catch (error) {
            console.error('Error parsing WebSocket message:', error);
            // Try to handle as plain text
            this.triggerEvent('message', {
              type: EventType.INFO,
              version: this.config.apiVersion,
              timestamp: new Date().toISOString(),
              message: event.data
            });
          }
        };

        this.websocket.onclose = () => {
          console.log('WebSocket disconnected');
          this.triggerEvent('disconnect', { 
            type: 'disconnect', 
            version: this.config.apiVersion, 
            timestamp: new Date().toISOString() 
          });
          this.websocket = null;
          this.scheduleReconnect();
        };

        this.websocket.onerror = (error) => {
          console.error('WebSocket error:', error);
          this.triggerEvent('error', { 
            type: 'error', 
            version: this.config.apiVersion, 
            timestamp: new Date().toISOString(),
            message: 'WebSocket error'
          });
          reject(error);
        };
      } catch (error) {
        this.isConnecting = false;
        console.error('Error creating WebSocket:', error);
        reject(error);
      }
    });
  }

  /**
   * Disconnect from the WebSocket server
   */
  public disconnect(): void {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }

    if (this.websocket) {
      this.websocket.close();
      this.websocket = null;
    }
  }

  /**
   * Schedule a reconnect attempt
   */
  private scheduleReconnect(): void {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
    }

    this.reconnectTimer = setTimeout(() => {
      console.log('Attempting to reconnect WebSocket...');
      this.connect().catch(() => {
        // If reconnect fails, schedule another attempt
        this.scheduleReconnect();
      });
    }, 3000); // Try to reconnect after 3 seconds
  }

  /**
   * Register an event handler
   */
  public on<T>(eventType: string, handler: EventHandler<T>): void {
    if (!this.eventHandlers.has(eventType)) {
      this.eventHandlers.set(eventType, []);
    }
    this.eventHandlers.get(eventType)?.push(handler as EventHandler);
  }

  /**
   * Unregister an event handler
   */
  public off(eventType: string, handler?: EventHandler): void {
    if (!handler) {
      // Remove all handlers for this event type
      this.eventHandlers.delete(eventType);
      return;
    }

    const handlers = this.eventHandlers.get(eventType);
    if (handlers) {
      const index = handlers.indexOf(handler);
      if (index !== -1) {
        handlers.splice(index, 1);
      }
      if (handlers.length === 0) {
        this.eventHandlers.delete(eventType);
      }
    }
  }

  /**
   * Handle an event from the server
   */
  private handleEvent(event: GameEvent): void {
    // First, trigger handlers for the specific event type
    this.triggerEvent(event.type, event);

    // Then, trigger handlers for 'all' events
    this.triggerEvent('all', event);

    // Handle specific event types
    switch (event.type) {
      case EventType.GAME_STATE:
        this.triggerEvent('gameState', event);
        break;
      case EventType.ATTACK:
      case EventType.ATTACK_RESULT:
        this.triggerEvent('attack', event);
        break;
      case EventType.DAMAGE:
      case EventType.DAMAGE_RESULT:
        this.triggerEvent('damage', event);
        break;
      case EventType.TURN_START:
      case EventType.TURN_END:
        this.triggerEvent('turn', event);
        break;
    }
  }

  /**
   * Trigger event handlers for a specific event type
   */
  private triggerEvent(eventType: string, event: GameEvent): void {
    const handlers = this.eventHandlers.get(eventType);
    if (handlers) {
      handlers.forEach(handler => {
        try {
          handler(event);
        } catch (error) {
          console.error(`Error in event handler for ${eventType}:`, error);
        }
      });
    }
  }

  /**
   * Get the current game state
   */
  public async getGameState(): Promise<GameState> {
    const response = await fetch(`${this.config.baseUrl}/api/${this.config.apiVersion}/state`);
    if (!response.ok) {
      throw new Error(`HTTP error ${response.status}: ${await response.text()}`);
    }
    const event = await response.json() as GameEvent<GameState>;
    return event.data as GameState;
  }

  /**
   * Get game steps starting from a specific index
   */
  public async getSteps(index: number = 0, limit: number = 20): Promise<GameEvent[]> {
    const response = await fetch(`${this.config.baseUrl}/api/${this.config.apiVersion}/steps?index=${index}&limit=${limit}`);
    if (!response.ok) {
      throw new Error(`HTTP error ${response.status}: ${await response.text()}`);
    }
    return await response.json() as GameEvent[];
  }

  /**
   * Send a command to the server
   */
  public async sendCommand(command: CommandRequest): Promise<CommandResponse> {
    const response = await fetch(`${this.config.baseUrl}/api/${this.config.apiVersion}/action`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(command)
    });

    if (!response.ok) {
      throw new Error(`HTTP error ${response.status}: ${await response.text()}`);
    }

    return await response.json() as CommandResponse;
  }
}