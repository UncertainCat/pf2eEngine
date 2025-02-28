package controllerhttp

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// ServeStaticFiles configures the server to serve the React frontend
func (cs *ControllerServer) ServeStaticFiles() {
	// Serve static files from the frontend/build directory
	frontendDir := "frontend/build"

	// Check if directory exists
	if _, err := os.Stat(frontendDir); os.IsNotExist(err) {
		fmt.Printf("Warning: Frontend directory %s does not exist.\n", frontendDir)
		return
	}

	// Handle static files with proper MIME types
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		filePath := filepath.Join(frontendDir, r.URL.Path)
		fmt.Printf("Static file request: %s -> %s\n", r.URL.Path, filePath)

		// Set proper Content-Type based on file extension
		ext := strings.ToLower(filepath.Ext(filePath))
		switch ext {
		case ".css":
			w.Header().Set("Content-Type", "text/css")
		case ".js":
			w.Header().Set("Content-Type", "application/javascript")
		case ".html":
			w.Header().Set("Content-Type", "text/html")
		case ".jpg", ".jpeg":
			w.Header().Set("Content-Type", "image/jpeg")
		case ".png":
			w.Header().Set("Content-Type", "image/png")
		case ".svg":
			w.Header().Set("Content-Type", "image/svg+xml")
		case ".json":
			w.Header().Set("Content-Type", "application/json")
		}

		// Actually serve the file
		http.ServeFile(w, r, filePath)
	})

	// Handle root requests (serve index.html)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Skip API endpoints
		if isAPIPath(r.URL.Path) {
			return
		}

		fmt.Printf("UI request: %s\n", r.URL.Path)
		
		// Set content type explicitly
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		
		// Serve index.html for all non-API paths
		http.ServeFile(w, r, filepath.Join(frontendDir, "index.html"))
	})

	fmt.Println("Static file server configured. Frontend available at http://localhost:8080")
}

// isAPIPath determines if a path should be handled by API controllers
func isAPIPath(path string) bool {
	apiPaths := []string{
		"/action",
		"/steps",
		"/ws",
		"/api/v1/action",
		"/api/v1/steps",
		"/api/v1/state",
	}

	for _, apiPath := range apiPaths {
		if path == apiPath {
			return true
		}
	}

	// Also check if path starts with /api/
	if len(path) >= 5 && path[:5] == "/api/" {
		return true
	}

	return false
}
