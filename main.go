package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type DegreeOfSuccess int

const (
	CriticalFailure DegreeOfSuccess = iota
	Failure
	Success
	CriticalSuccess
)

type Entity struct {
	Name        string
	HP          int
	AC          int
	AttackBonus int
	DamageBonus int
	Initiative  int
}

// rollDice simulates a dice roll of n sides
func rollDice(sides int) int {
	return rand.Intn(sides) + 1
}

// calculateDegreeOfSuccess determines the degree of success for a roll
func calculateDegreeOfSuccess(roll int, modifier int, dc int) DegreeOfSuccess {
	result := roll + modifier
	if roll == 20 { // Natural 20 steps up one degree
		if result >= dc+10 {
			return CriticalSuccess
		} else if result >= dc {
			return CriticalSuccess // Step up from Success
		} else if result <= dc-10 {
			return Failure // Step up from Critical Failure
		} else {
			return Success // Step up from Failure
		}
	} else if roll == 1 { // Natural 1 steps down one degree
		if result >= dc+10 {
			return Success // Step down from Critical Success
		} else if result >= dc {
			return Failure // Step down from Success
		} else if result <= dc-10 {
			return CriticalFailure
		} else {
			return CriticalFailure // Step down from Failure
		}
	} else if result >= dc+10 {
		return CriticalSuccess
	} else if result >= dc {
		return Success
	} else if result <= dc-10 {
		return CriticalFailure
	} else {
		return Failure
	}
}

// degreeOfSuccessToString converts DegreeOfSuccess to a human-readable string
func degreeOfSuccessToString(degree DegreeOfSuccess) string {
	switch degree {
	case CriticalFailure:
		return "Critical Failure"
	case Failure:
		return "Failure"
	case Success:
		return "Success"
	case CriticalSuccess:
		return "Critical Success"
	default:
		return "Unknown"
	}
}

// attack performs an attack from the attacker to the defender
func attack(attacker, defender *Entity) {
	roll := rollDice(20)
	degree := calculateDegreeOfSuccess(roll, attacker.AttackBonus, defender.AC)
	fmt.Printf("%s rolls a %d (total %d) to attack %s (AC %d): %s\n", attacker.Name, roll, roll+attacker.AttackBonus, defender.Name, defender.AC, degreeOfSuccessToString(degree))

	switch degree {
	case CriticalSuccess:
		damage := (rollDice(8) + attacker.DamageBonus) * 2
		defender.HP -= damage
		fmt.Printf("Critical hit! %s deals %d damage to %s. %s's HP is now %d\n", attacker.Name, damage, defender.Name, defender.Name, defender.HP)
	case Success:
		damage := rollDice(8) + attacker.DamageBonus
		defender.HP -= damage
		fmt.Printf("Hit! %s deals %d damage to %s. %s's HP is now %d\n", attacker.Name, damage, defender.Name, defender.Name, defender.HP)
	case CriticalFailure:
		fmt.Printf("Critical miss! %s fumbles the attack.\n", attacker.Name)
	case Failure:
		fmt.Printf("Miss! %s fails to hit %s.\n", attacker.Name, defender.Name)
	}
}

// rollInitiative rolls for initiative for all entities and sorts them by result
func rollInitiative(entities []*Entity) {
	for _, entity := range entities {
		entity.Initiative = rollDice(20)
		fmt.Printf("%s rolls initiative: %d\n", entity.Name, entity.Initiative)
	}
	// Sort entities by initiative in descending order
	sort.Slice(entities, func(i, j int) bool {
		return entities[i].Initiative > entities[j].Initiative
	})
	fmt.Println("Initiative order:")
	for _, entity := range entities {
		fmt.Printf("%s (Initiative: %d)\n", entity.Name, entity.Initiative)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Seed random number generator

	// Create combatants
	entity1 := &Entity{Name: "Warrior", HP: 30, AC: 15, AttackBonus: 5, DamageBonus: 3}
	entity2 := &Entity{Name: "Goblin", HP: 20, AC: 13, AttackBonus: 3, DamageBonus: 2}

	entities := []*Entity{entity1, entity2}

	// Roll initiative
	rollInitiative(entities)

	fmt.Println("Combat begins!")

	// Main combat loop
	for {
		for _, entity := range entities {
			if entity.HP <= 0 {
				continue
			}

			// Find a target
			target := entities[0]
			if entity == target {
				target = entities[1]
			}

			attack(entity, target)

			// Check if combat is over
			aliveEntities := 0
			for _, e := range entities {
				if e.HP > 0 {
					aliveEntities++
				}
			}
			if aliveEntities <= 1 {
				fmt.Println("Combat ends!")
				for _, e := range entities {
					if e.HP > 0 {
						fmt.Printf("%s wins!\n", e.Name)
					}
				}
				return
			}
		}
	}
}
