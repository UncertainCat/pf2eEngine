package game

import (
	"pf2eEngine/entity"
)

// Position represents a coordinate on the grid.
type Position struct {
	X, Y int
}

// Grid represents the game grid, managing entities and their positions.
type Grid struct {
	Width  int
	Height int
	Cells  map[Position]*entity.Entity
}

// NewGrid initializes a new grid with the given dimensions.
func NewGrid(width, height int) *Grid {
	return &Grid{
		Width:  width,
		Height: height,
		Cells:  make(map[Position]*entity.Entity),
	}
}

// AddEntity places an entity at a specific position on the grid.
func (g *Grid) AddEntity(pos Position, e *entity.Entity) bool {
	if !g.IsValidPosition(pos) || g.IsOccupied(pos) {
		return false
	}
	g.Cells[pos] = e
	return true
}

// MoveEntity moves an entity from one position to another.
func (g *Grid) MoveEntity(from, to Position) bool {
	if !g.IsValidPosition(to) || g.IsOccupied(to) {
		return false
	}
	if entity, exists := g.Cells[from]; exists {
		delete(g.Cells, from)
		g.Cells[to] = entity
		return true
	}
	return false
}

// RemoveEntity removes an entity from the grid.
func (g *Grid) RemoveEntity(pos Position) {
	delete(g.Cells, pos)
}

// IsOccupied checks if a position is occupied by an entity.
func (g *Grid) IsOccupied(pos Position) bool {
	_, exists := g.Cells[pos]
	return exists
}

// IsValidPosition checks if a position is within grid bounds.
func (g *Grid) IsValidPosition(pos Position) bool {
	return pos.X >= 0 && pos.Y >= 0 && pos.X < g.Width && pos.Y < g.Height
}

// GetEntityAt retrieves the entity at a specific position.
func (g *Grid) GetEntityAt(pos Position) *entity.Entity {
	return g.Cells[pos]
}

// GetAdjacentPositions returns all valid adjacent positions to a given position.
func (g *Grid) GetAdjacentPositions(pos Position) []Position {
	candidates := []Position{
		{X: pos.X, Y: pos.Y - 1}, // Up
		{X: pos.X, Y: pos.Y + 1}, // Down
		{X: pos.X - 1, Y: pos.Y}, // Left
		{X: pos.X + 1, Y: pos.Y}, // Right
	}
	valid := []Position{}
	for _, p := range candidates {
		if g.IsValidPosition(p) {
			valid = append(valid, p)
		}
	}
	return valid
}

// CalculateDistance computes the Manhattan distance between two positions.
func (g *Grid) CalculateDistance(from, to Position) int {
	dx := to.X - from.X
	dy := to.Y - from.Y
	return abs(dx) + abs(dy)
}

// abs returns the absolute value of an integer.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
