// API Client implementation
import { EventEmitter } from 'events';

// Event types
export const GameEvent = {
    GAME_STATE: 'GAME_STATE',
    TURN_START: 'TURN_START',
    DAMAGE: 'DAMAGE',
    DAMAGE_RESULT: 'DAMAGE_RESULT'
};

// Entity states
export const EntityState = {
    IDLE: 'IDLE',
    ACTING: 'ACTING',
    DEAD: 'DEAD'
};

class ApiClient extends EventEmitter {
    constructor() {
        super();
        this.socket = null;
        this.baseUrl = process.env.NODE_ENV === 'production' 
            ? window.location.origin
            : 'http://localhost:8080';
    }

    // Connect to WebSocket with auto-reconnect
    connect() {
        if (this.socket && (this.socket.readyState === WebSocket.CONNECTING || 
                           this.socket.readyState === WebSocket.OPEN)) {
            console.log('WebSocket already connected or connecting');
            return;
        }
        
        const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${wsProtocol}//${window.location.host}/ws`;
        
        console.log(`Connecting to WebSocket at ${wsUrl}`);
        this.socket = new WebSocket(wsUrl);
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 5;
        this.reconnectDelay = 1000; // Start with 1 second delay

        this.socket.onopen = () => {
            console.log('WebSocket connection established');
            this.reconnectAttempts = 0;
            this.reconnectDelay = 1000; // Reset reconnect delay
            this.emit('connect');
        };

        this.socket.onclose = (event) => {
            // Code 1000 = normal closure, 1001 = going away
            const isNormalClosure = event.code === 1000 || event.code === 1001;
            console.log(`WebSocket connection closed: ${event.code} - ${event.reason || 'No reason provided'}`);
            this.emit('disconnect');
            
            // Don't auto-reconnect if this was a normal closure or we're manually disconnecting
            if (!this.manualDisconnect && !isNormalClosure) {
                this.scheduleReconnect();
            }
        };

        this.socket.onerror = (error) => {
            console.error('WebSocket error:', error);
            this.emit('error', error);
            
            // Only attempt reconnect if we haven't reached our max attempts
            if (this.reconnectAttempts < this.maxReconnectAttempts) {
                this.scheduleReconnect();
            }
        };

        this.socket.onmessage = (message) => {
            try {
                const event = JSON.parse(message.data);
                // Log only type to reduce noise
                console.log(`Received WebSocket event of type: ${event.type}`);
                
                // Emit specific event
                if (event.type) {
                    const eventType = event.type.toLowerCase();
                    this.emit(eventType, event);
                }
                
                // Emit catchall event
                this.emit('all', event);
            } catch (error) {
                console.error('Error parsing message:', error);
            }
        };
    }
    
    // Schedule a reconnect with exponential backoff
    scheduleReconnect() {
        if (this.reconnectAttempts >= this.maxReconnectAttempts) {
            console.log('Maximum reconnect attempts reached');
            return;
        }
        
        this.reconnectAttempts++;
        const delay = Math.min(30000, this.reconnectDelay * Math.pow(1.5, this.reconnectAttempts - 1));
        console.log(`Scheduling reconnect attempt ${this.reconnectAttempts} in ${delay}ms`);
        
        setTimeout(() => {
            console.log(`Attempting to reconnect (attempt ${this.reconnectAttempts})`);
            this.connect();
        }, delay);
    }

    // Disconnect WebSocket
    disconnect() {
        this.manualDisconnect = true;
        if (this.socket) {
            this.socket.close(1000, "Normal closure");
            this.socket = null;
        }
        this.manualDisconnect = false;
    }

    // Get current game state
    async getGameState() {
        const response = await fetch(`${this.baseUrl}/api/v1/state`);
        if (!response.ok) {
            throw new Error(`Failed to get game state: ${response.statusText}`);
        }
        return response.json();
    }
    
    // Get step history
    async getStepHistory(index = 0, limit = 100) {
        const response = await fetch(`${this.baseUrl}/api/v1/steps?index=${index}&limit=${limit}`);
        if (!response.ok) {
            throw new Error(`Failed to get step history: ${response.statusText}`);
        }
        return response.json();
    }

    // Send a command to the game
    async sendCommand(command) {
        const response = await fetch(`${this.baseUrl}/api/v1/action`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(command)
        });
        
        if (!response.ok) {
            throw new Error(`Failed to send command: ${response.statusText}`);
        }
        
        return response.json();
    }
}

// Export a singleton instance
export const apiClient = new ApiClient();