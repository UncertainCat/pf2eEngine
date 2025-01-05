package game

import (
	"pf2eEngine/entity"
)

// GameState represents the state of the game, including turn order and current turn
type GameState struct {
	Entities     []*entity.Entity
	CurrentIndex int
}

// NewGameState initializes a new game state with the given entities
func NewGameState(entities []*entity.Entity) *GameState {
	return &GameState{
		Entities:     entities,
		CurrentIndex: 0,
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

// GetNextLivingTarget finds the next living target for the current entity
func (gs *GameState) GetNextLivingTarget(current *entity.Entity) *entity.Entity {
	for _, e := range gs.Entities {
		if e != current && e.IsAlive() {
			return e
		}
	}
	return nil
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
