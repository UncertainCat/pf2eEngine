package game

// Position represents a coordinate on the grid.
type Position struct {
	X, Y int
}

// Grid represents the game grid, managing entities and their positions.
type Grid struct {
	Width  int
	Height int
	Cells  map[Position]*Entity
}

// NewGrid initializes a new grid with the given dimensions.
func NewGrid(width, height int) *Grid {
	return &Grid{
		Width:  width,
		Height: height,
		Cells:  make(map[Position]*Entity),
	}
}

// AddEntity places an entity at a specific position on the grid.
func (g *Grid) AddEntity(pos Position, e *Entity) bool {
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
func (g *Grid) GetEntityAt(pos Position) *Entity {
	return g.Cells[pos]
}

// GetEntityPosition retrieves the position of a specific entity on the grid.
func (g *Grid) GetEntityPosition(e *Entity) Position {
	for pos, entity := range g.Cells {
		if entity == e {
			return pos
		}
	}
	return Position{-1, -1} // Invalid position if entity is not found
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

// AreAdjacent checks if two positions are adjacent on the grid.
func (g *Grid) AreAdjacent(pos1, pos2 Position) bool {
	// In grid coordinates, positions are adjacent if they are at most 1 square away in each direction
	dx := abs(pos2.X - pos1.X)
	dy := abs(pos2.Y - pos1.Y)
	return dx <= 1 && dy <= 1 && !(dx == 0 && dy == 0)
}

// CalculateDistance computes the distance between two positions using PF2E rules.
// In PF2E, diagonal movement costs 5ft for the first square, then alternates between 10ft and 5ft.
func (g *Grid) CalculateDistance(from, to Position) int {
	dx := abs(to.X - from.X)
	dy := abs(to.Y - from.Y)
	
	// Calculate diagonal movement count
	diagonals := min(dx, dy)
	
	// Calculate straight movement count
	straights := max(dx, dy) - diagonals
	
	// Apply PF2E diagonal rules: first diagonal is 5ft, second is 10ft, third is 5ft, etc.
	diagonalDistance := 0
	for i := 0; i < diagonals; i++ {
		if i%2 == 0 {
			diagonalDistance += 5
		} else {
			diagonalDistance += 10
		}
	}
	
	// Straight movements are always 5ft each
	straightDistance := straights * 5
	
	return diagonalDistance + straightDistance
}

// CalculateDistanceBetweenEntities computes the distance between two entities.
func (g *Grid) CalculateDistanceBetweenEntities(e1, e2 *Entity) int {
	pos1 := g.GetEntityPosition(e1)
	pos2 := g.GetEntityPosition(e2)
	return g.CalculateDistance(pos1, pos2)
}

// abs returns the absolute value of an integer.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// min returns the minimum of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max returns the maximum of two integers.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// FindBestMove finds the best position to move towards the target from a starting position,
// considering the maximum movement range and grid obstacles.
func (g *Grid) FindBestMove(from Position, toward Position, maxDistance int) Position {
	// Simple approach: try to move closer in both X and Y dimensions
	dx := toward.X - from.X
	dy := toward.Y - from.Y
	
	// Determine step direction for X and Y
	stepX := 0
	if dx > 0 {
		stepX = 1
	} else if dx < 0 {
		stepX = -1
	}
	
	stepY := 0
	if dy > 0 {
		stepY = 1
	} else if dy < 0 {
		stepY = -1
	}
	
	// Try to move along both dimensions
	bestPosition := from
	bestDistance := g.CalculateDistance(from, toward)
	
	// Try moving up to maxDistance steps
	for steps := 1; steps <= maxDistance; steps++ {
		// Try moving in various directions 
		candidates := []Position{
			{X: from.X + stepX*steps, Y: from.Y}, // Move horizontally
			{X: from.X, Y: from.Y + stepY*steps}, // Move vertically
			{X: from.X + stepX*min(steps, abs(dx)), Y: from.Y + stepY*min(steps, abs(dy))}, // Move diagonally
		}
		
		for _, candidate := range candidates {
			// Check if the position is valid and unoccupied
			if g.IsValidPosition(candidate) && !g.IsOccupied(candidate) {
				candidateDistance := g.CalculateDistance(candidate, toward)
				if candidateDistance < bestDistance {
					bestPosition = candidate
					bestDistance = candidateDistance
				}
			}
		}
	}
	
	return bestPosition
}
