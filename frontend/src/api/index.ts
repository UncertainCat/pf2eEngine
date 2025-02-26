// Export all API types and client
export * from './types';
export * from './client';

// Create and export a singleton instance of the API client
import { ApiClient } from './client';

// Singleton instance
export const apiClient = new ApiClient();

// Make sure we connect to WebSocket
apiClient.connect().catch(error => {
  console.error('Failed to connect to WebSocket:', error);
});