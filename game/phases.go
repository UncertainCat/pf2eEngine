package game

import (
	"fmt"
	"pf2eEngine/combat"
	"pf2eEngine/entity"
)

// SetupPhase handles any pre-combat setup like rolling initiative
func SetupPhase(entities []*entity.Entity) {
	fmt.Println("Starting setup phase...")
	combat.RollInitiative(entities)
}

// CombatPhase runs the main combat loop
func CombatPhase(gs *GameState) {
	fmt.Println("Starting combat phase...")

	for !gs.IsCombatOver() {
		currentEntity := gs.GetCurrentTurnEntity()
		target := gs.GetNextLivingTarget(currentEntity)
		if target != nil {
			combat.PerformAttack(currentEntity, target)
		}
		gs.EndTurn()
	}
}

// ResolutionPhase determines the outcome of the combat
func ResolutionPhase(gs *GameState) {
	fmt.Println("Starting resolution phase...")
	winner := gs.GetWinner()
	if winner != nil {
		fmt.Printf("%s wins the combat!\n", winner.Name)
	} else {
		fmt.Println("Combat ends with no clear winner.")
	}
}
