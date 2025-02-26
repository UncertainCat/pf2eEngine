# PF2E Engine - Developer Guide

## Build & Run Commands
- Build backend: `go build -o pf2e-server main.go`
- Run backend: `./pf2e-server` or `go run main.go`
- Build frontend: `cd frontend && npm install && npm run build`
- Run frontend (dev): `cd frontend && npm start`
- Full build & run: `make all`
- Run tests: `go test ./...`
- Run single test: `go test -v ./path/to/package -run TestName`

## Code Style Guidelines
- **Imports**: Standard library first, then external packages, then local packages
- **Naming**: CamelCase for exported names, camelCase for unexported
- **Error Handling**: Return errors instead of panicking; use descriptive error messages
- **Comments**: Document exported functions, types, and constants with Go-style comments
- **Types**: Use strong typing with structs and interfaces; avoid interface{}
- **Functions**: Keep functions small and focused on a single responsibility
- **Formatting**: Always run `go fmt` before committing code
- **Package Structure**: Organize by domain concepts (game, combat, entity)
- **Tests**: Write tests for all public functions; aim for >80% coverage

Remember to seed random generators for deterministic tests.