// main.go
package main

import (
	"math/rand"
	"pf2eEngine/entity"
	"pf2eEngine/game"
	"pf2eEngine/items"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator

	// Create combatants
	warrior := entity.NewEntity("Warrior", 30, 15, 5, 3)
	goblin := entity.NewEntity("Goblin", 20, 13, 3, 2)
	game.RegisterTrigger(items.ShieldBlock{Owner: warrior}, "BEFORE_DAMAGE")
	combatants := []*entity.Entity{warrior, goblin}

	// Start the combat loop
	game.StartCombat(combatants)
}
