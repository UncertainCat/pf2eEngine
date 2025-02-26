package controllerhttp

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

// ServeStaticFiles configures the server to serve the React frontend
func (cs *ControllerServer) ServeStaticFiles() {
	// Serve static files from the frontend/build directory
	frontendDir := "frontend/build"

	// Check if directory exists
	if _, err := os.Stat(frontendDir); os.IsNotExist(err) {
		fmt.Printf("Warning: Frontend directory %s does not exist. Static files will not be served.\n", frontendDir)
		return
	}

	// Create a file server to serve static assets
	fs := http.FileServer(http.Dir(frontendDir))

	// Handle routes for static files
	http.Handle("/static/", fs)
	http.Handle("/assets/", fs)
	http.Handle("/favicon.ico", fs)

	// Special handling for index.html
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// API endpoints should still go to their handlers
		if isAPIPath(r.URL.Path) {
			return
		}

		// For all other paths, serve index.html
		indexPath := filepath.Join(frontendDir, "index.html")
		http.ServeFile(w, r, indexPath)
	})

	fmt.Println("Static file server configured. Frontend available at http://localhost:8080")
}

// isAPIPath determines if a path should be handled by API controllers
func isAPIPath(path string) bool {
	apiPaths := []string{
		"/action",
		"/steps",
		"/ws",
	}

	for _, apiPath := range apiPaths {
		if path == apiPath {
			return true
		}
	}

	return false
}
