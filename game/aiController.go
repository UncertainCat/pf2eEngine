package game

import (
	"math"
)

// AIController implements basic AI logic for entities
type AIController struct{}

func NewAIController() AIController {
	return AIController{}
}

func (a AIController) NextAction(gs *GameState, e *Entity) Action {
	target := findNearestEntity(gs, e)
	if target == nil || e.ActionsRemaining == 0 {
		return EndTurnAction(gs, e)
	}
	return Strike(target)
}

// findNearestEntity locates the closest living entity to the given entity
func findNearestEntity(gs *GameState, e *Entity) *Entity {
	currentPos := gs.Grid.GetEntityPosition(e)
	var closest *Entity
	minDistance := math.MaxInt

	for _, other := range gs.Entities {
		if other == e || !other.IsAlive() {
			continue
		}
		otherPos := gs.Grid.GetEntityPosition(other)
		distance := gs.Grid.CalculateDistance(currentPos, otherPos)
		if distance < minDistance {
			closest = other
			minDistance = distance
		}
	}
	return closest
}
