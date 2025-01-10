package game

import (
	"fmt"
	"pf2eEngine/entity"
	dice "pf2eEngine/util"
	"sort"
)

// GameState represents the state of the game, including turn order and current turn
type GameState struct {
	Entities     []*entity.Entity
	CurrentIndex int
}

func NewCombat(entities []*entity.Entity) *GameState {
	gs := &GameState{
		Entities:     entities,
		CurrentIndex: 0,
	}
	gs.RollInitiative()
	gs.RunCombatLoop()
	return gs
}

// RollInitiative determines initiative for all combatants and sorts them in descending order
func (gs *GameState) RollInitiative() {
	for _, e := range gs.Entities {
		roll := dice.Roll(20)
		e.RollInitiative(roll)
		fmt.Printf("%s rolls initiative: %d\n", e.Name, e.Initiative)
	}

	sort.Slice(gs.Entities, func(i, j int) bool {
		return gs.Entities[i].Initiative > gs.Entities[j].Initiative
	})

	fmt.Println("Initiative order:")
	for i, e := range gs.Entities {
		fmt.Printf("%d: %s (Initiative: %d)\n", i+1, e.Name, e.Initiative)
	}
}

// IsCombatOver checks if the combat has ended
func (gs *GameState) IsCombatOver() bool {
	livingEntities := 0
	for _, e := range gs.Entities {
		if e.IsAlive() {
			livingEntities++
		}
	}
	return livingEntities <= 1
}

// GetCurrentTurnEntity returns the entity whose turn it is
func (gs *GameState) GetCurrentTurnEntity() *entity.Entity {
	return gs.Entities[gs.CurrentIndex]
}

// StartTurn resets the current entity's actions and reactions
func (gs *GameState) StartTurn() {
	currentEntity := gs.GetCurrentTurnEntity()
	currentEntity.ResetTurnResources()
	fmt.Printf("%s's turn begins: %d actions, %d reactions available.\n", currentEntity.Name, currentEntity.ActionsRemaining, currentEntity.ReactionsRemaining)
}

// EndTurn moves the turn to the next living entity
func (gs *GameState) EndTurn() {
	for {
		gs.CurrentIndex = (gs.CurrentIndex + 1) % len(gs.Entities)
		if gs.Entities[gs.CurrentIndex].IsAlive() {
			break
		}
	}
}

// GetWinner returns the winning entity, if combat is over
func (gs *GameState) GetWinner() *entity.Entity {
	if !gs.IsCombatOver() {
		return nil
	}
	for _, e := range gs.Entities {
		if e.IsAlive() {
			return e
		}
	}
	return nil
}

// RunCombatLoop handles the main combat loop
func (gs *GameState) RunCombatLoop() {
	for !gs.IsCombatOver() {
		currentEntity := gs.GetCurrentTurnEntity()
		gs.StartTurn()

		// Example: Find and attack the next target
		target := currentEntity.GetNextLivingTarget(gs.Entities)
		if target != nil {
			fmt.Printf("%s targets %s for an attack.\n", currentEntity.Name, target.Name)
			PerformAttack(currentEntity, target)
			PerformAttack(currentEntity, target)
			PerformAttack(currentEntity, target)
		} else {
			fmt.Printf("%s has no valid targets to attack.\n", currentEntity.Name)
		}

		gs.EndTurn()
	}

	winner := gs.GetWinner()
	if winner != nil {
		fmt.Printf("%s wins the combat!\n", winner.Name)
	} else {
		fmt.Println("Combat ends with no winner.")
	}
}
