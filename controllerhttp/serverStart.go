package controllerhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Start starts the HTTP server and registers both HTTP and WebSocket handlers.
func (cs *ControllerServer) Start() {
	// API endpoints with CORS handling
	http.HandleFunc("/action", cs.corsMiddleware(cs.HTTPHandler))
	http.HandleFunc("/steps", cs.corsMiddleware(cs.StepsHandler))
	http.HandleFunc("/ws", cs.WSHandler) // WebSocket doesn't need CORS

	// Serve static frontend files
	cs.ServeStaticFiles()

	fmt.Printf("Server running on port %d\n", cs.Port)
	fmt.Printf("Access the frontend at http://localhost:%d\n", cs.Port)

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

		// Create a simplified game state for the frontend
		type EntityData struct {
			ID               string `json:"id"`
			Name             string `json:"name"`
			HP               int    `json:"hp"`
			MaxHP            int    `json:"maxHp"`
			AC               int    `json:"ac"`
			ActionsRemaining int    `json:"actionsRemaining"`
			Faction          int    `json:"faction"`
			X                int    `json:"x"`
			Y                int    `json:"y"`
			IsAlive          bool   `json:"isAlive"`
		}

		entities := make([]EntityData, 0, len(cs.GameState.Initiative))
		for _, entity := range cs.GameState.Initiative {
			// Get entity position
			pos := cs.GameState.Grid.GetEntityPosition(entity)

			entityData := EntityData{
				ID:               entity.Id.String(),
				Name:             entity.Name,
				HP:               entity.HP,
				MaxHP:            entity.HP, // Use initial HP as max HP for now
				AC:               entity.AC,
				ActionsRemaining: entity.ActionsRemaining,
				Faction:          int(entity.Faction),
				X:                pos.X,
				Y:                pos.Y,
				IsAlive:          entity.IsAlive(),
			}
			entities = append(entities, entityData)
		}

		// Create and send update message
		updateMsg := map[string]interface{}{
			"type":        "gameUpdate",
			"entities":    entities,
			"currentTurn": cs.GameState.GetCurrentTurnEntity().Id.String(),
		}

		jsonData, err := json.Marshal(updateMsg)
		if err == nil {
			cs.wsBroadcast <- jsonData
		}
	}
}
