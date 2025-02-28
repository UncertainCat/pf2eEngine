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
	
	// Get current positions
	actorPos := gs.Grid.GetEntityPosition(e)
	targetPos := gs.Grid.GetEntityPosition(target)
	
	// Prepare common parameters
	params := map[string]interface{}{}
	params["targetID"] = target.Id
	
	// If not adjacent to target, try to stride towards them
	if !gs.Grid.AreAdjacent(actorPos, targetPos) {
		// Look for a stride card
		strideCard := getActionCardByName(e, "Stride")
		if strideCard != nil {
			action, err := strideCard.GenerateAction(gs, e, params)
			if err == nil {
				fmt.Printf("%s decides to stride towards %s.\n", e.Name, target.Name)
				return action
			}
		}
	}
	
	// If adjacent to target or no stride card available, try to strike
	strikes := getActionCardsByName(e, "Strike")
	if len(strikes) > 0 {
		action, err := strikes[0].GenerateAction(gs, e, params)
		if err == nil {
			return action
		}
	}
	
	// If no valid action could be generated, end turn
	fmt.Printf("%s has no valid actions remaining.\n", e.Name)
	return EndTurnAction(gs, e)
}

func getActionCardByName(e *Entity, name string) *ActionCard {
	for _, card := range e.ActionCards {
		if card.Name == name {
			return card
		}
	}
	return nil
}

func getActionCardsByName(e *Entity, name string) []*ActionCard {
	var cards []*ActionCard
	for _, card := range e.ActionCards {
		if card.Name == name {
			cards = append(cards, card)
		}
	}
	return cards
}

func getStrikeCards(e *Entity) []*ActionCard {
	return getActionCardsByName(e, "Strike")
}

// findNearestEntity locates the closest living entity to the given entity
func findNearestEnemy(gs *GameState, e *Entity) *Entity {
	currentPos := gs.Grid.GetEntityPosition(e)
	var closest *Entity
	minDistance := math.MaxInt

	for _, other := range gs.Initiative {
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
