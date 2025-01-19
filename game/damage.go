package game

import (
	"fmt"
)

type Damage struct {
	Source  *Entity
	Target  *Entity
	Amount  int
	Blocked int
	Taken   int
}

type BeforeDamageStep struct {
	BaseStep
	Damage *Damage
}

type AfterDamageStep struct {
	BaseStep
	Damage *Damage
}

func NewBeforeDamageStep(damage *Damage) BeforeDamageStep {
	return BeforeDamageStep{
		BaseStep: BaseStep{
			stepType: BeforeDamage,
			metadata: map[string]interface{}{
				"Source": damage.Source.Name,
				"Target": damage.Target.Name,
				"Amount": damage.Amount,
			},
		},
		Damage: damage,
	}
}

func NewAfterDamageStep(damage *Damage) AfterDamageStep {
	return AfterDamageStep{
		BaseStep: BaseStep{
			stepType: AfterDamage,
			metadata: map[string]interface{}{
				"Source":  damage.Source.Name,
				"Target":  damage.Target.Name,
				"Amount":  damage.Amount,
				"Blocked": damage.Blocked,
				"Taken":   damage.Taken,
			},
		},
		Damage: damage,
	}
}

func Deal(gs *GameState, damage Damage) {
	executeStep(gs, NewBeforeDamageStep(&damage), fmt.Sprintf("%s is about to deal damage to %s.", damage.Source.Name, damage.Target.Name))
	totalDamage := damage.Amount - damage.Blocked
	if totalDamage < 0 {
		totalDamage = 0
	}
	applyDamage(damage, totalDamage)
	damage.Taken = totalDamage

	executeStep(gs, NewAfterDamageStep(&damage), fmt.Sprintf("%s dealt %d damage to %s.", damage.Source.Name, damage.Taken, damage.Target.Name))
}

func applyDamage(damage Damage, totalDamage int) {
	damage.Target.TakeDamage(totalDamage)
}
