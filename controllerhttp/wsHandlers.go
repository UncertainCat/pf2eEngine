package controllerhttp

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

// upgrader is used to upgrade HTTP connections to WebSocket connections.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// For simplicity, allow all origins. In production, you should validate this.
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WSHandler upgrades the HTTP connection to a WebSocket and listens for commands.
func (cs *ControllerServer) WSHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("WebSocket upgrade error: %v\n", err)
		return
	}
	// Register the new client.
	cs.wsClients[conn] = true
	fmt.Println("New WebSocket client connected")

	// Listen for incoming messages.
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("WebSocket read error: %v\n", err)
			delete(cs.wsClients, conn)
			break
		}

		// Only process text messages.
		if messageType != websocket.TextMessage {
			continue
		}

		var command CommandRequest
		if err := json.Unmarshal(message, &command); err != nil {
			errMsg := "Invalid JSON"
			conn.WriteMessage(websocket.TextMessage, []byte(errMsg))
			continue
		}

		err = cs.Controller.AddActionWithCard(command.EntityID, command.ActionCardId, command.Params)
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			continue
		}

		successMsg := "Action queued successfully"
		conn.WriteMessage(websocket.TextMessage, []byte(successMsg))
	}
}

// broadcastUpdates listens on wsBroadcast and sends incoming messages to all connected WebSocket clients.
func (cs *ControllerServer) broadcastUpdates() {
	for {
		msg := <-cs.wsBroadcast
		for client := range cs.wsClients {
			if err := client.WriteMessage(websocket.TextMessage, msg); err != nil {
				fmt.Printf("WebSocket write error: %v\n", err)
				client.Close()
				delete(cs.wsClients, client)
			}
		}
	}
}
