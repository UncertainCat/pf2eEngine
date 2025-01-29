package controllerhttp

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"pf2eEngine/game"
	"strconv"
)

// ControllerServer handles HTTP requests for player-controlled entities
type ControllerServer struct {
	Port       int
	Controller *game.PlayerController
	GameState  *game.GameState
}

// NewControllerServer initializes a ControllerServer
func NewControllerServer(port int, controller *game.PlayerController) *ControllerServer {
	return &ControllerServer{
		Port:       port,
		Controller: controller,
	}
}

type CommandRequest struct {
	EntityID     uuid.UUID              `json:"entity_id"`
	ActionCardId uuid.UUID              `json:"action_card_id"`
	Params       map[string]interface{} `json:"params"`
}

// HTTPHandler processes incoming HTTP requests and queues actions
func (cs *ControllerServer) HTTPHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}
	command := CommandRequest{}
	if err := json.NewDecoder(r.Body).Decode(&command); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err := cs.Controller.AddActionWithCard(command.EntityID, command.ActionCardId, command.Params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// StepsHandler processes GET requests to return game steps starting from a given index
func (cs *ControllerServer) StepsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	indexStr := r.URL.Query().Get("index")
	if indexStr == "" {
		http.Error(w, "Index query parameter is required", http.StatusBadRequest)
		return
	}

	index, err := strconv.Atoi(indexStr)
	if err != nil || index < 0 {
		http.Error(w, "Invalid index", http.StatusBadRequest)
		return
	}

	if index >= len(cs.GameState.StepHistory.Steps) {
		http.Error(w, "Index out of range", http.StatusBadRequest)
		return
	}

	steps := cs.GameState.StepHistory.Steps[index:]
	response, err := json.Marshal(steps)
	if err != nil {
		fmt.Printf("Error marshaling steps: %v\n", err)
		http.Error(w, "Failed to marshal steps", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// Start starts the HTTP server
func (cs *ControllerServer) Start() {
	http.HandleFunc("/action", cs.HTTPHandler)
	http.HandleFunc("/steps", cs.StepsHandler)
	fmt.Printf("HTTP server running on port %d\n", cs.Port)
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", cs.Port), nil); err != nil {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()
}
