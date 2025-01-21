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

	d8Plus3Slashing := game.DamageRoll{
		Die:   8,
		Bonus: 3,
	}

	d6Plus1Piercing := game.DamageRoll{
		Die:   6,
		Bonus: 1,
	}

	// Create combatants
	warrior := game.NewEntity("Warrior", 30, 15)
	warriorAttack := game.BaseAttack{
		Damage: []game.DamageRoll{d8Plus3Slashing},
		Bonus:  5,
	}
	warrior.AddActionCard(game.NewStrikeCard(warriorAttack))
	goblin := game.NewEntity("Goblin", 20, 13)
	goblinAttack := game.BaseAttack{
		Damage: []game.DamageRoll{d6Plus1Piercing},
		Bonus:  3,
	}
	goblin.AddActionCard(game.NewStrikeCard(goblinAttack))
	game.RegisterTrigger(items.ShieldBlock{Owner: warrior}, "BEFORE_DAMAGE")
	spawns := []game.Spawn{
		{Unit: goblin, Coordinates: [2]int{0, 0}},
		{Unit: warrior, Coordinates: [2]int{1, 0}},
	}

	// Start the combat loop
	game.StartCombat(spawns, 10, 10)
}
