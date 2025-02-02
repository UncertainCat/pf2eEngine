package controllerhttp

import (
	"fmt"
	"net/http"
)

// Start starts the HTTP server and registers both HTTP and WebSocket handlers.
func (cs *ControllerServer) Start() {
	http.HandleFunc("/action", cs.HTTPHandler)
	http.HandleFunc("/steps", cs.StepsHandler)
	http.HandleFunc("/ws", cs.WSHandler)
	fmt.Printf("Server running on port %d\n", cs.Port)

	// Start the broadcast goroutine.
	go cs.broadcastUpdates()

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", cs.Port), nil); err != nil {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()
}
