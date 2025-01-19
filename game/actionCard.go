package game

import (
	"errors"
	"github.com/google/uuid"
)

type ActionCardType string

const (
	OneActionCard      ActionCardType = "ONE_ACTION"
	TwoActionCard      ActionCardType = "TWO_ACTION"
	ThreeActionCard    ActionCardType = "THREE_ACTION"
	VariableActionCard ActionCardType = "VARIABLE_ACTION"
	FreeActionCard     ActionCardType = "FREE_ACTION"
)

// ActionCard defines a reusable template for actions
type ActionCard struct {
	ID              uuid.UUID
	Name            string
	Type            ActionCardType
	Description     string
	ActionGenerator func(owner *Entity, state GameState, params map[string]interface{}) (Action, error)
}

// NewActionCard initializes an action card with given parameters
func NewActionCard(name string, actionType ActionCardType, description string) *ActionCard {
	return &ActionCard{
		ID:          uuid.New(),
		Name:        name,
		Type:        actionType,
		Description: description,
	}
}

// CreateAction instantiates an action based on the card and additional parameters
func (ac *ActionCard) CreateAction(gs GameState, owner *Entity, params map[string]interface{}) (Action, error) {
	return ac.ActionGenerator(owner, gs, params)
}

var (
	ErrInvalidParams  = errors.New("invalid parameters")
	ErrEntityNotFound = errors.New("entity not found")
)

// Basic Action Cards

func NewStrikeCard() *ActionCard {
	return &ActionCard{
		ID:          uuid.New(),
		Name:        "Strike",
		Type:        OneActionCard,
		Description: "Make a melee Strike against a target.",
		ActionGenerator: func(owner *Entity, state GameState, params map[string]interface{}) (Action, error) {
			targetID, ok := params["target"].(uuid.UUID)
			if !ok {
				return Action{}, ErrInvalidParams
			}

			target := findEntityByID(state.Entities, targetID)
			if target == nil {
				return Action{}, ErrEntityNotFound
			}

			return Action{
				Name:        "Strike",
				Type:        SingleAction,
				Cost:        1,
				Description: "Melee attack action",
				Perform: func(gs *GameState, actor *Entity) {
					PerformAttack(gs, actor, target)
				},
			}, nil
		},
	}
}
