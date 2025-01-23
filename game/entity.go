package game

import (
	"fmt"
	"github.com/google/uuid"
)

type Entity struct {
	Id                 uuid.UUID
	Name               string
	HP                 int
	AC                 int
	Initiative         int
	ActionsRemaining   int
	ReactionsRemaining int
	MapCounter         int
	Controller         Controller
	ActionCards        []*ActionCard
	Faction            Faction
}

func (e *Entity) AddActionCard(card *ActionCard) {
	e.ActionCards = append(e.ActionCards, card)
}

type Faction int

const (
	GoodGuys Faction = iota
	BadGuys
)

// NewEntity creates a new Entity instance
func NewEntity(name string, hp, ac int, faction Faction) *Entity {
	return &Entity{
		Id:                 uuid.New(),
		Name:               name,
		HP:                 hp,
		AC:                 ac,
		ActionsRemaining:   3,
		ReactionsRemaining: 1,
		Controller:         NewAIController(),
		Faction:            faction,
	}
}

// IsAlive checks if the entity is still alive
func (e *Entity) IsAlive() bool {
	return e.HP > 0
}

// SetController assigns a controllerhttp to the entity
func (e *Entity) SetController(controller Controller) {
	e.Controller = controller
}

// TakeDamage reduces the entity's HP by a given amount
func (e *Entity) TakeDamage(damage int) {
	e.HP -= damage
	if e.HP < 0 {
		e.HP = 0
	}
	fmt.Printf("%s takes %d damage! Remaining HP: %d\n", e.Name, damage, e.HP)
}

// RollInitiative sets the entity's initiative value
func (e *Entity) RollInitiative(roll int) {
	e.Initiative = roll
}

// ResetTurnResources resets actions and reactions at the start of a turn
func (e *Entity) ResetTurnResources() {
	e.ActionsRemaining = 3
	e.ReactionsRemaining = 1
	e.MapCounter = 0
}

// SpendAction attempts to consume an action
func (e *Entity) SpendAction(cost int) bool {
	if e.ActionsRemaining >= cost {
		e.ActionsRemaining -= cost
		return true
	}
	return false
}

// UseReaction attempts to consume a reaction
func (e *Entity) UseReaction() bool {
	if e.ReactionsRemaining > 0 {
		e.ReactionsRemaining--
		return true
	}
	return false
}

// GetNextLivingTarget finds the next living target for the entity
func (e *Entity) GetNextLivingTarget(entities []*Entity) *Entity {
	for _, target := range entities {
		if target != e && target.IsAlive() {
			return target
		}
	}
	return nil
}

// String provides a readable string representation of the entity
func (e *Entity) String() string {
	return fmt.Sprintf("%s (HP: %d, AC: %d, Actions: %d, Reactions: %d)", e.Name, e.HP, e.AC, e.ActionsRemaining, e.ReactionsRemaining)
}

func findEntityByID(entities []*Entity, id uuid.UUID) *Entity {
	for _, e := range entities {
		if e.Id == id {
			return e
		}
	}
	return nil
}

func findActionCardByID(entity *Entity, id uuid.UUID) (*ActionCard, error) {
	for _, card := range entity.GetActionCards() {
		if card.ID == id {
			return card, nil
		}
	}
	return nil, fmt.Errorf("action card not found")
}

func (e *Entity) GetActionCards() []*ActionCard {
	return e.ActionCards
}
