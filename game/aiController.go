package game

import (
	"fmt"
	"math"
)

// AIController implements basic AI logic for entities
type AIController struct{}

func NewAIController() AIController {
	return AIController{}
}

func (a AIController) NextAction(gs *GameState, e *Entity) Action {
	target := findNearestEnemy(gs, e)
	if target == nil || e.ActionsRemaining == 0 {
		return EndTurnAction(gs, e)
	}
	strikes := getStrikeCards(e)
	if len(strikes) == 0 {
		fmt.Println("No strike cards available for", e.Name)
		return EndTurnAction(gs, e)
	}
	params := map[string]interface{}{}
	params["targetID"] = target.Id
	action, err := strikes[0].GenerateAction(gs, e, params)
	if err != nil {
		return EndTurnAction(gs, e)
	}
	return action
}

func getStrikeCards(e *Entity) []*ActionCard {
	cards := e.ActionCards
	// Filter out all non-strike cards
	var strikeCards []*ActionCard
	for _, card := range cards {
		if card.Name == "Strike" {
			strikeCards = append(strikeCards, card)
		}
	}
	return strikeCards
}

// findNearestEntity locates the closest living entity to the given entity
func findNearestEnemy(gs *GameState, e *Entity) *Entity {
	currentPos := gs.Grid.GetEntityPosition(e)
	var closest *Entity
	minDistance := math.MaxInt

	for _, other := range gs.Entities {
		if other == e || !other.IsAlive() || other.Faction == e.Faction {
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
