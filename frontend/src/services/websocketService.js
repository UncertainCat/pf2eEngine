import { apiClient } from '../api';
import { EventEmitter } from 'events';

// WebSocket event types
export const WsEvent = {
    CONNECT: 'connect',
    DISCONNECT: 'disconnect',
    ERROR: 'error',
    GAME_EVENT: 'game_event',
};

class WebsocketService extends EventEmitter {
    constructor() {
        super();
        this.connected = false;
        
        // Set up event handlers
        apiClient.on('connect', () => {
            this.connected = true;
            this.emit(WsEvent.CONNECT);
        });
        
        apiClient.on('disconnect', () => {
            this.connected = false;
            this.emit(WsEvent.DISCONNECT);
        });
        
        apiClient.on('error', (error) => {
            this.emit(WsEvent.ERROR, error);
        });
        
        // Forward all game events
        apiClient.on('all', (event) => {
            this.emit(WsEvent.GAME_EVENT, event);
        });
    }
    
    // Connect to the WebSocket
    connect() {
        if (this.connected) {
            console.log('WebSocket already connected');
            return;
        }
        
        apiClient.connect();
    }
    
    // Disconnect from the WebSocket
    disconnect() {
        if (!this.connected) {
            console.log('WebSocket already disconnected');
            return;
        }
        
        apiClient.disconnect();
    }
    
    // Check if currently connected
    isConnected() {
        return this.connected;
    }
}

// Export a singleton instance
export const websocketService = new WebsocketService();