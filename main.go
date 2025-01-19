// main.go
package main

import (
	"math/rand"
	"pf2eEngine/game"
	"pf2eEngine/items"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator

	// Create combatants
	warrior := game.NewEntity("Warrior", 30, 15, 5, 3)
	goblin := game.NewEntity("Goblin", 20, 13, 3, 2)
	game.RegisterTrigger(items.ShieldBlock{Owner: warrior}, "BEFORE_DAMAGE")
	spawns := []game.Spawn{
		{Unit: goblin, Coordinates: [2]int{0, 0}},
		{Unit: warrior, Coordinates: [2]int{1, 0}},
	}

	// Start the combat loop
	game.StartCombat(spawns, 10, 10)
}
