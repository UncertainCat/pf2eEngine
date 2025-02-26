package controllerhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pf2eEngine/controllerhttp/api"
	"strconv"
	"time"
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
	
	// Create API response
	response := api.CommandResponse{
		Success: true,
		Message: "Action queued successfully",
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// StepsHandler processes GET requests to return game steps starting from a given index.
func (cs *ControllerServer) StepsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get index parameter with default value of 0
	indexStr := r.URL.Query().Get("index")
	index := 0
	if indexStr != "" {
		var err error
		index, err = strconv.Atoi(indexStr)
		if err != nil || index < 0 {
			http.Error(w, "Invalid index", http.StatusBadRequest)
			return
		}
	}
	
	// Get limit parameter with default value of 20
	limitStr := r.URL.Query().Get("limit")
	limit := 20
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}
	}
	
	// Get steps from history
	history := cs.GameState.StepHistory.GetSteps()
	
	// Check if index is valid
	if index >= len(history) {
		// Return empty array if index is out of bounds
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
		return
	}
	
	// Calculate end index
	endIndex := index + limit
	if endIndex > len(history) {
		endIndex = len(history)
	}
	
	// Convert steps to API events
	events := make([]api.GameEvent, 0, endIndex-index)
	for i := index; i < endIndex; i++ {
		step := history[i]
		message := fmt.Sprintf("Step %d", i)
		event := api.StepToEvent(step, message)
		events = append(events, event)
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(events)
}

// GameStateHandler returns the current game state in the API format
func (cs *ControllerServer) GameStateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Convert game state to API format
	gameState := api.GameStateToAPIState(cs.GameState)
	
	// Create an event wrapper
	event := api.GameEvent{
		EventBase: api.EventBase{
			Type:      api.EventTypeGameState,
			Version:   api.CurrentVersion,
			Timestamp: time.Now(),
			Message:   "Current game state",
		},
		Data: gameState,
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(event)
}