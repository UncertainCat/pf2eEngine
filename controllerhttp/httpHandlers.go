package controllerhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// HTTPHandler processes incoming HTTP POST requests and queues actions.
func (cs *ControllerServer) HTTPHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}
	var command CommandRequest
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

// StepsHandler processes GET requests to return game steps starting from a given index.
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
