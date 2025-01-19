package controllerhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pf2eEngine/game"
)

// ControllerServer handles HTTP requests for player-controlled entities
type ControllerServer struct {
	Port       int
	Controller *game.PlayerController
}

// NewControllerServer initializes a ControllerServer
func NewControllerServer(port int, controller *game.PlayerController) *ControllerServer {
	return &ControllerServer{
		Port:       port,
		Controller: controller,
	}
}

// HTTPHandler processes incoming HTTP requests and queues actions
func (cs *ControllerServer) HTTPHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		Action string
		Target string // Optional target entity name
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	action := game.Action{
		Name:        input.Action,
		Type:        game.SingleAction,
		Cost:        1,
		Description: fmt.Sprintf("%s action issued via HTTP", input.Action),
		Perform: func(gs *game.GameState, actor *game.Entity) {
			if input.Target != "" {
				for _, target := range gs.Entities {
					if target.Name == input.Target {
						game.PerformAttack(gs, actor, target)
						return
					}
				}
				fmt.Printf("Target '%s' not found.\n", input.Target)
			}
		},
	}

	command := game.PlayerCommand{
		Action: action,
	}
	cs.Controller.AddAction(action)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Action '%s' queued for entity.", input.Action)
}

// Start starts the HTTP server
func (cs *ControllerServer) Start() {
	http.HandleFunc("/action", cs.HTTPHandler)
	fmt.Printf("HTTP server running on port %d\n", cs.Port)
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", cs.Port), nil); err != nil {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()
}
