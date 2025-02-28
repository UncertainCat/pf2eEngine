package main

import (
	"fmt"
	"math/rand"
	"pf2eEngine/controllerhttp"
	"pf2eEngine/game"
	"pf2eEngine/items"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator

	d8Plus3Slashing := game.DamageRoll{
		Die:   8,
		Count: 1,
		Bonus: 3,
		Type:  game.Slashing,
	}

	// Create combatants
	warrior := game.NewEntity("Warrior", 30, 15, game.GoodGuys)
	warriorAttack := game.BaseAttack{
		Damage: []game.DamageRoll{d8Plus3Slashing},
		Bonus:  5,
	}
	warrior.AddActionCard(game.NewStrikeCard(warriorAttack))
	warrior.AddActionCard(game.NewStrideCard())
	
	goblin1 := makeAGoblin("Goblin 1")
	goblin2 := makeAGoblin("Goblin 2")
	goblin3 := makeAGoblin("Goblin 3")
	goblin4 := makeAGoblin("Goblin 4")
	
	game.RegisterTrigger(items.ShieldBlock{Owner: warrior}, "BEFORE_DAMAGE")
	
	// Position entities with more spacing to demonstrate grid-based movement
	spawns := []game.Spawn{
		{Unit: goblin1, Coordinates: [2]int{0, 0}},    // Top left
		{Unit: goblin2, Coordinates: [2]int{6, 0}},    // Top right
		{Unit: warrior, Coordinates: [2]int{3, 5}},    // Middle bottom - player is farther away
		{Unit: goblin3, Coordinates: [2]int{0, 8}},    // Bottom left
		{Unit: goblin4, Coordinates: [2]int{6, 8}},    // Bottom right
	}

	// Initialize game state
	gameState := game.NewGameState(spawns, 10, 10)

	// Initialize player controller
	playerController := game.NewPlayerController(gameState)

	// Initialize the HTTP server
	server := controllerhttp.NewControllerServer(8080, playerController)

	// Connect game state to server
	gameState.StepCallback = server.BroadcastGameStep
	server.GameState = gameState

	// Start the server
	server.Start()

	// Don't automatically start combat - wait for player to step through manually
	fmt.Println("Server is ready. Use the frontend to step through combat.")

	// Prevent the main function from exiting
	select {}
}

func makeAGoblin(name string) *game.Entity {
	goblin := game.NewEntity(name, 20, 13, game.BadGuys)
	goblinAttack := game.BaseAttack{
		Damage: []game.DamageRoll{{Die: 6, Count: 1, Bonus: 1, Type: game.Piercing}},
		Bonus:  3,
	}
	goblin.AddActionCard(game.NewStrikeCard(goblinAttack))
	goblin.AddActionCard(game.NewStrideCard()) // Add Stride action card
	return goblin
}
