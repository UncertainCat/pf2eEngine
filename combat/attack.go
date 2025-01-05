package combat

import (
	"fmt"
	dice "pf2eEngine/util"

	"pf2eEngine/entity"
)

type DegreeOfSuccess int

const (
	CriticalFailure DegreeOfSuccess = iota
	Failure
	Success
	CriticalSuccess
)

// calculateDegreeOfSuccess determines the degree of success for an attack roll
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

// PerformAttack executes an attack from one entity to another
func PerformAttack(attacker *entity.Entity, defender *entity.Entity) {
	roll := dice.Roll(20)
	degree := calculateDegreeOfSuccess(roll, attacker.AttackBonus, defender.AC)
	fmt.Printf("%s rolls a %d (total %d) to attack %s (AC %d): %v\n", attacker.Name, roll, roll+attacker.AttackBonus, defender.Name, defender.AC, degree)

	switch degree {
	case CriticalSuccess:
		damage := (dice.Roll(8) + attacker.DamageBonus) * 2
		defender.TakeDamage(damage)
		fmt.Printf("Critical hit! %s deals %d damage to %s.\n", attacker.Name, damage, defender.Name)
	case Success:
		damage := dice.Roll(8) + attacker.DamageBonus
		defender.TakeDamage(damage)
		fmt.Printf("Hit! %s deals %d damage to %s.\n", attacker.Name, damage, defender.Name)
	case CriticalFailure:
		fmt.Printf("Critical miss! %s fumbles the attack.\n", attacker.Name)
	case Failure:
		fmt.Printf("Miss! %s fails to hit %s.\n", attacker.Name, defender.Name)
	}
}
