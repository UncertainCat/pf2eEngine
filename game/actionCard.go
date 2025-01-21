package game

import (
	"errors"
	"fmt"
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

const (
	ActionCost = "action_cost"
	TargetID   = "targetID"
)

func (a ActionCardType) ToCost(params map[string]interface{}) int {
	switch a {
	case OneActionCard:
		return 1
	case TwoActionCard:
		return 2
	case ThreeActionCard:
		return 3
	case VariableActionCard:
		return params[ActionCost].(int)
	case FreeActionCard:
		return 0
	default:
		return 0
	}
}

type ActionCard struct {
	ID              uuid.UUID
	Name            string
	Type            ActionCardType
	Description     string
	ActionGenerator func(gs *GameState, actor *Entity, params map[string]interface{}) (Action, error)
}

func (ac ActionCard) GenerateAction(gs *GameState, actor *Entity, params map[string]interface{}) (Action, error) {
	action, err := ac.ActionGenerator(gs, actor, params)
	if err != nil {
		fmt.Printf("Failed to generate action: %v\n Params: %v\n", err, params)
		return Action{}, err
	}
	return action, nil
}

func getSingleTarget(gs *GameState, actor *Entity, criteria []TargetCriterion, params map[string]interface{}) (*Entity, error) {
	targetID := params[TargetID].(uuid.UUID)
	target := findEntityByID(gs.Entities, targetID)
	if target == nil {
		return nil, errors.New("target not found")
	}
	for _, criterion := range criteria {
		if err := criterion(gs, actor, params); err != nil {
			return nil, err
		}
	}
	return target, nil
}

func getTarget(gs *GameState, params map[string]interface{}) (*Entity, error) {
	targetID, err := params[TargetID].(uuid.UUID)
	if !err {
		return nil, errors.New("targetID not found in params")
	}
	target := findEntityByID(gs.Entities, targetID)
	if target == nil {
		return nil, errors.New("target entity does not exist")
	}
	return target, nil
}

type TargetCriterion func(gs *GameState, actor *Entity, params map[string]interface{}) error

func Range(maxDistance int) TargetCriterion {
	return func(gs *GameState, actor *Entity, params map[string]interface{}) error {
		target, err := getTarget(gs, params)
		if err != nil {
			return err
		}
		gs.Grid.GetEntityPosition(actor)
		if gs.Grid.CalculateDistanceBetweenEntities(target, actor) > maxDistance {
			return fmt.Errorf("target is out of range (max: %d)", maxDistance)
		}
		return nil
	}
}

func IsAlive() TargetCriterion {
	return func(gs *GameState, actor *Entity, params map[string]interface{}) error {
		targetId, ok := params[TargetID].(uuid.UUID)
		target := findEntityByID(gs.Entities, targetId)
		if !ok {
			return errors.New("target not found")
		}
		if !target.IsAlive() {
			return errors.New("target is not alive")
		}
		return nil
	}
}

func NewSingleTargetActionCard(
	name string,
	actionType ActionCardType,
	description string,
	criteria []TargetCriterion,
	actionFunc func(gs *GameState, actor *Entity, target *Entity),
) *ActionCard {
	return &ActionCard{
		ID:          uuid.New(),
		Name:        name,
		Type:        actionType,
		Description: description,
		ActionGenerator: func(gs *GameState, actor *Entity, params map[string]interface{}) (Action, error) {
			target, err := getSingleTarget(gs, actor, criteria, params)
			if err != nil {
				return Action{}, err
			}
			return Action{
				Name: name,
				Cost: actionType.ToCost(params),
				Perform: func(gs *GameState, actor *Entity) {
					actionFunc(gs, actor, target)
				},
			}, nil
		},
	}
}

func NewStrikeCard(attack BaseAttack) *ActionCard {
	return NewSingleTargetActionCard(
		"Strike",
		OneActionCard,
		"Make a melee strike against a target within range 5.",
		[]TargetCriterion{
			IsAlive(),
			Range(5),
		},
		func(gs *GameState, actor *Entity, target *Entity) {
			PerformAttack(gs, attack, actor, target)
		},
	)
}
