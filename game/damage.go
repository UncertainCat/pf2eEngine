package game

import (
	"fmt"
	"pf2eEngine/entity"
)

// Damage struct
type Damage struct {
	Source  *entity.Entity
	Target  *entity.Entity
	Amount  int
	Blocked int
	Taken   int
}

type BeforeDamageStep struct {
	*Damage
}

func (d BeforeDamageStep) Type() StepType {
	return "BEFORE_DAMAGE"
}

type AfterDamageStep struct {
	*Damage
}

func (d AfterDamageStep) Type() StepType {
	return "AFTER_DAMAGE"
}

// Confirm before damage implements the step interface
var _ Step = BeforeDamageStep{}

// Deal function
func Deal(damage Damage) {
	executeStep(BeforeDamageStep{&damage})
	totalDamage := damage.Amount - damage.Blocked
	if totalDamage < 0 {
		totalDamage = 0
	}

	fmt.Printf("Total damage after block: %d (original: %d, blocked: %d)\n", totalDamage, damage.Amount, damage.Blocked)
	applyDamage(damage, totalDamage)
	damage.Taken = totalDamage
	afterDamageStep := AfterDamageStep{&damage}
	executeStep(&afterDamageStep)
}

func applyDamage(damage Damage, totalDamage int) {
	damage.Target.TakeDamage(totalDamage)
}
