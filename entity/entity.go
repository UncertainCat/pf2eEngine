package entity

import "fmt"

// Entity represents a combatant in the game
type Entity struct {
	Name        string
	HP          int
	AC          int
	AttackBonus int
	DamageBonus int
	Initiative  int
}

// NewEntity creates a new Entity instance
func NewEntity(name string, hp, ac, attackBonus, damageBonus int) *Entity {
	return &Entity{
		Name:        name,
		HP:          hp,
		AC:          ac,
		AttackBonus: attackBonus,
		DamageBonus: damageBonus,
	}
}

// IsAlive checks if the entity is still alive
func (e *Entity) IsAlive() bool {
	return e.HP > 0
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

// String provides a readable string representation of the entity
func (e *Entity) String() string {
	return fmt.Sprintf("%s (HP: %d, AC: %d, Initiative: %d)", e.Name, e.HP, e.AC, e.Initiative)
}
