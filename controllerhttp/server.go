package controllerhttp

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"pf2eEngine/game"
)

// ControllerServer handles HTTP requests and WebSocket connections for player-controlled entities.
type ControllerServer struct {
	Port       int
	Controller *game.PlayerController
	GameState  *game.GameState

	// wsClients holds all active WebSocket connections.
	wsClients map[*websocket.Conn]bool
	// wsBroadcast is a channel for broadcasting messages to all WebSocket clients.
	wsBroadcast chan []byte
}

// NewControllerServer initializes a ControllerServer.
func NewControllerServer(port int, controller *game.PlayerController) *ControllerServer {
	return &ControllerServer{
		Port:        port,
		Controller:  controller,
		wsClients:   make(map[*websocket.Conn]bool),
		wsBroadcast: make(chan []byte),
	}
}

// CommandRequest represents a command sent by HTTP or WebSocket clients.
type CommandRequest struct {
	EntityID     uuid.UUID              `json:"entity_id"`
	ActionCardId uuid.UUID              `json:"action_card_id"`
	Params       map[string]interface{} `json:"params"`
}
