package main

import (
	"fmt"
	"math/rand"
	"pf2eEngine/items"
	"time"

	"pf2eEngine/combat"
	"pf2eEngine/entity"
	"pf2eEngine/game"
)

func main() {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator

	// Create combatants
	warrior := entity.NewEntity("Warrior", 30, 15, 5, 3)
	goblin := entity.NewEntity("Goblin", 20, 13, 3, 2)
	goblin2 := entity.NewEntity("Goblin2", 20, 13, 3, 2)
	game.RegisterTrigger(items.ShieldBlock{Owner: warrior}, "BEFORE_DAMAGE")
	combatants := []*entity.Entity{warrior, goblin, goblin2}

	// Roll initiative and determine turn order
	combat.RollInitiative(combatants)

	fmt.Println("Combat begins!")

	// Initialize the game state and start the combat loop
	gameState := game.NewGameState(combatants)
	for !gameState.IsCombatOver() {
		currentEntity := gameState.GetCurrentTurnEntity()
		target := gameState.GetNextLivingTarget(currentEntity)
		if target != nil {
			combat.PerformAttack(currentEntity, target)
		}
		gameState.EndTurn()
	}

	// Announce the winner
	winner := gameState.GetWinner()
	if winner != nil {
		fmt.Printf("%s wins the combat!\n", winner.Name)
	} else {
		fmt.Println("Combat ends with no winner.")
	}
}
