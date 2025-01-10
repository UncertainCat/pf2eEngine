package game

import (
	"fmt"
	"pf2eEngine/entity"
	dice "pf2eEngine/util"
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
	if !attacker.UseAction(1) {
		fmt.Printf("%s does not have enough actions to attack!\n", attacker.Name)
		return
	}

	roll := dice.Roll(20)
	attackBonus := attacker.AttackBonus - attacker.MapCounter*5
	degree := calculateDegreeOfSuccess(roll, attackBonus, defender.AC)
	fmt.Printf("%s rolls a %d (total %d) to attack %s (AC %d): %v\n", attacker.Name, roll, roll+attackBonus, defender.Name, defender.AC, degree)

	switch degree {
	case CriticalSuccess:
		damage := (dice.Roll(8) + attacker.DamageBonus) * 2
		fmt.Printf("Critical hit! %s deals %d damage to %s.\n", attacker.Name, damage, defender.Name)
		Deal(Damage{Source: attacker, Target: defender, Amount: damage})
	case Success:
		damage := dice.Roll(8) + attacker.DamageBonus
		fmt.Printf("Hit! %s deals %d damage to %s.\n", attacker.Name, damage, defender.Name)
		Deal(Damage{Source: attacker, Target: defender, Amount: damage})
	case CriticalFailure:
		fmt.Printf("Critical miss! %s fumbles the attack.\n", attacker.Name)
	case Failure:
		fmt.Printf("Miss! %s fails to hit %s.\n", attacker.Name, defender.Name)
	}
	attacker.MapCounter++
}
