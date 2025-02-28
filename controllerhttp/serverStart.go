package controllerhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pf2eEngine/controllerhttp/api"
	"time"
)

// Start starts the HTTP server and registers both HTTP and WebSocket handlers.
func (cs *ControllerServer) Start() {
	// API endpoints with CORS handling
	http.HandleFunc("/api/v1/action", cs.corsMiddleware(cs.HTTPHandler))
	http.HandleFunc("/api/v1/steps", cs.corsMiddleware(cs.StepsHandler))
	http.HandleFunc("/api/v1/state", cs.corsMiddleware(cs.GameStateHandler))
	http.HandleFunc("/ws", cs.WSHandler) // WebSocket doesn't need CORS
	
	// Support legacy endpoints for backward compatibility
	http.HandleFunc("/action", cs.corsMiddleware(cs.HTTPHandler))
	http.HandleFunc("/steps", cs.corsMiddleware(cs.StepsHandler))

	// Serve static frontend files
	cs.ServeStaticFiles()

	fmt.Printf("Server running on port %d\n", cs.Port)
	fmt.Printf("Access the frontend at http://localhost:%d\n", cs.Port)
	fmt.Printf("API endpoints available at http://localhost:%d/api/v1/\n", cs.Port)

	// Start the broadcast goroutine.
	go cs.broadcastUpdates()

	// Start game state broadcast to WebSocket clients
	go cs.startGameStateUpdates()

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", cs.Port), nil); err != nil {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()
}

// corsMiddleware adds CORS headers to API responses
func (cs *ControllerServer) corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next(w, r)
	}
}

// startGameStateUpdates periodically sends game state updates to WebSocket clients
func (cs *ControllerServer) startGameStateUpdates() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if cs.GameState == nil {
			continue
		}

		// Convert game state to API format
		gameState := api.GameStateToAPIState(cs.GameState)
		
		// Create an event wrapper
		event := api.GameEvent{
			EventBase: api.EventBase{
				Type:      api.EventTypeGameState,
				Version:   api.CurrentVersion,
				Timestamp: time.Now(),
				Message:   "Game state update",
				Metadata: map[string]interface{}{
					"isInitial": false,
					"isUpdate": true,
				},
			},
			Data: gameState,
		}
		
		jsonData, err := json.Marshal(event)
		if err == nil {
			cs.wsBroadcast <- jsonData
		} else {
			fmt.Printf("Error marshaling game state update: %v\n", err)
		}
	}
}
