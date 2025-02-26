.PHONY: all build run clean frontend

# Default target
all: frontend build run

# Build Go backend
build:
	go build -o pf2e-server main.go

# Build React frontend
frontend:
	cd frontend && npm install && npm run build

# Run the server
run:
	./pf2e-server

# Clean build artifacts
clean:
	rm -f pf2e-server
	rm -rf frontend/node_modules
	rm -rf frontend/build

# Dev mode: run frontend dev server and Go server separately
dev-frontend:
	cd frontend && npm start

dev-backend:
	go run main.go