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
	actionGenerator func(gs *GameState, actor *Entity, params map[string]interface{}) (Action, error)
}

func (ac ActionCard) GenerateAction(gs *GameState, actor *Entity, params map[string]interface{}) (Action, error) {
	action, err := ac.actionGenerator(gs, actor, params)
	if err != nil {
		fmt.Printf("Failed to generate action: %v\n Params: %v\n", err, params)
		return Action{}, err
	}
	return action, nil
}

func getSingleTarget(gs *GameState, actor *Entity, criteria []TargetCriterion, params map[string]interface{}) (*Entity, error) {
	targetID := params[TargetID].(uuid.UUID)
	target := findEntityByID(gs.Initiative, targetID)
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
	target := findEntityByID(gs.Initiative, targetID)
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
		target := findEntityByID(gs.Initiative, targetId)
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
		actionGenerator: func(gs *GameState, actor *Entity, params map[string]interface{}) (Action, error) {
			target, err := getSingleTarget(gs, actor, criteria, params)
			if err != nil {
				return Action{}, err
			}
			return Action{
				Name: name,
				Cost: actionType.ToCost(params),
				perform: func(gs *GameState, actor *Entity) {
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
		"Make a melee strike against a target within reach.",
		[]TargetCriterion{
			IsAlive(),
			Range(5),
		},
		func(gs *GameState, actor *Entity, target *Entity) {
			PerformAttack(gs, attack, actor, target)
		},
	)
}

// NewStrideCard creates a movement action card according to PF2E rules.
func NewStrideCard() *ActionCard {
	return &ActionCard{
		ID:          uuid.New(),
		Name:        "Stride",
		Type:        OneActionCard,
		Description: "Move up to your Speed (default: 5 squares).",
		actionGenerator: func(gs *GameState, actor *Entity, params map[string]interface{}) (Action, error) {
			targetID, ok := params["targetID"].(uuid.UUID)
			if !ok {
				return Action{}, errors.New("target not found in params")
			}
			
			target := findEntityByID(gs.Initiative, targetID)
			if target == nil {
				return Action{}, errors.New("target entity does not exist")
			}
			
			return Action{
				Name: "Stride",
				Cost: 1,
				perform: func(gs *GameState, actor *Entity) {
					actorPos := gs.Grid.GetEntityPosition(actor)
					targetPos := gs.Grid.GetEntityPosition(target)
					
					// Default movement speed in PF2E is 25 feet (5 squares)
					speed := 5
					
					// Find the best position to move to
					newPos := gs.Grid.FindBestMove(actorPos, targetPos, speed)
					
					// If we found a valid position to move to
					if newPos != actorPos {
						success := gs.Grid.MoveEntity(actorPos, newPos)
						if success {
							fmt.Printf("%s strides from (%d,%d) to (%d,%d).\n", 
								actor.Name, actorPos.X, actorPos.Y, newPos.X, newPos.Y)
						} else {
							fmt.Printf("%s attempted to stride but was blocked.\n", actor.Name)
						}
					} else {
						fmt.Printf("%s cannot stride any closer to %s.\n", actor.Name, target.Name)
					}
				},
			}, nil
		},
	}
}
