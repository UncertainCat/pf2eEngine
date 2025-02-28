package game

import (
	"fmt"
	"time"
)

// Controller interface for entity controllers
type Controller interface {
	NextAction(gs *GameState, entity *Entity) Action
}

// RunCombat runs the combat simulation automatically
func RunCombat(gs *GameState) {
	fmt.Println("Starting combat simulation...")
	gs.StartCombat()
	
	// Run combat until someone wins or we stop manually
	for {
		// Get current entity
		entity := gs.GetCurrentTurnEntity()
		if entity == nil {
			fmt.Println("Combat is over, no more entities.")
			break
		}
		
		// Skip if entity is dead
		if !entity.IsAlive() {
			gs.NextTurn()
			continue
		}
		
		// Get controller decision for AI entities
		controller := entity.Controller
		if controller != nil {
			// Process all actions for this entity's turn
			for entity.ActionsRemaining > 0 {
				// Let the controller decide what to do
				action := controller.NextAction(gs, entity)
				
				// If no action was chosen or it's an end turn action, move on
				if action.Name == "" || action.Type == EndOfTurn {
					break
				}
				
				// Execute the action
				ExecuteAction(gs, entity, action)
				
				// Short pause between actions
				time.Sleep(500 * time.Millisecond)
			}
		}
		
		// Move to next entity
		gs.NextTurn()
		
		// Slight pause between turns
		time.Sleep(1 * time.Second)
	}
}