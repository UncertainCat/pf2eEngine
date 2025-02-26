package game

import (
	"fmt"
	"github.com/google/uuid"
)

// PlayerController listens for HTTP commands to control an entity
type PlayerController struct {
	GameState  *GameState
	ActionChan chan Action // Channel for receiving actions
}

// NewPlayerController initializes a PlayerController
func NewPlayerController(gs *GameState) *PlayerController {
	return &PlayerController{
		GameState:  gs,
		ActionChan: make(chan Action, 1), // Buffered channel to hold one action
	}
}

// NextAction retrieves the next action from the channel
func (p *PlayerController) NextAction(gs *GameState, e *Entity) Action {
	return <-p.ActionChan
}

type ErrNotEntityTurn struct {
	EntityId uuid.UUID
}

func (e ErrNotEntityTurn) Error() string {
	return "Not entity turn"
}

type PlayerCommand struct {
	Action   Action
	EntityId uuid.UUID
}

func (p *PlayerController) AddAction(c PlayerCommand) error {
	if p.GameState.GetCurrentTurnEntity().Id != c.EntityId {
		return ErrNotEntityTurn{
			EntityId: c.EntityId,
		}
	}
	p.ActionChan <- c.Action
	return nil
}

func (p *PlayerController) AddActionWithCard(entityID uuid.UUID, actionCardID uuid.UUID, params map[string]interface{}) error {
	entity := findEntityByID(p.GameState.Initiative, entityID)
	if entity == nil {
		return fmt.Errorf("entity not found")
	}
	if !p.GameState.IsEntityTurn(entity) {
		return fmt.Errorf("not entity's turn")
	}

	actionCard, err := findActionCardByID(entity, actionCardID)
	if err != nil {
		return fmt.Errorf("action card not found")
	}

	action, err := actionCard.GenerateAction(p.GameState, entity, params)
	if err != nil {
		return fmt.Errorf("failed to create action: %w", err)
	}
	p.ActionChan <- action
	return nil
}
