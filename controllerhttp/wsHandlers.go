package controllerhttp

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"pf2eEngine/game"
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

// BroadcastGameStep sends a game step to all WebSocket clients in a frontend-friendly format
func (cs *ControllerServer) BroadcastGameStep(step interface{}, message string) {
	// Convert the step to a frontend-friendly format
	stepData := map[string]interface{}{
		"message": message,
	}

	// Add specific data based on step type
	switch s := step.(type) {
	case game.BeforeAttackStep:
		stepData["type"] = "ATTACK"
		stepData["data"] = map[string]interface{}{
			"attacker": s.Attack.Attacker.Name,
			"defender": s.Attack.Defender.Name,
			"roll":     s.Attack.Roll,
			"result":   s.Attack.Result,
		}
	case game.AfterAttackStep:
		stepData["type"] = "ATTACK_RESULT"
		stepData["data"] = map[string]interface{}{
			"attacker": s.Attack.Attacker.Name,
			"defender": s.Attack.Defender.Name,
			"degree":   s.Attack.Degree.String(),
		}
	case game.BeforeDamageStep:
		stepData["type"] = "DAMAGE"
		stepData["data"] = map[string]interface{}{
			"source": s.Damage.Source.Name,
			"target": s.Damage.Target.Name,
			"amount": s.Damage.Amount,
		}
	case game.AfterDamageStep:
		stepData["type"] = "DAMAGE_RESULT"
		stepData["data"] = map[string]interface{}{
			"source":  s.Damage.Source.Name,
			"target":  s.Damage.Target.Name,
			"blocked": s.Damage.Blocked,
			"taken":   s.Damage.Taken,
		}
	case game.StartTurnStep:
		stepData["type"] = "TURN_START"
		if s.Entity != nil {
			stepData["data"] = map[string]interface{}{
				"entity": map[string]interface{}{
					"name": s.Entity.Name,
					"id":   s.Entity.Id.String(),
				},
			}
		}
	case game.EndTurnStep:
		stepData["type"] = "TURN_END"
		if s.Entity != nil {
			stepData["data"] = map[string]interface{}{
				"entity": map[string]interface{}{
					"name": s.Entity.Name,
					"id":   s.Entity.Id.String(),
				},
			}
		}
	default:
		stepData["type"] = "INFO"
	}

	// Convert to JSON and broadcast
	jsonData, err := json.Marshal(stepData)
	if err != nil {
		fmt.Printf("Error marshaling step data: %v\n", err)
		return
	}

	cs.wsBroadcast <- jsonData
}
