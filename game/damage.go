package game

import (
	"fmt"
	"pf2eEngine/entity"
)

type Damage struct {
	Source  *entity.Entity
	Target  *entity.Entity
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
				"Source": damage.Source,
				"Target": damage.Target,
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
				"Source": damage.Source,
				"Target": damage.Target,
				"Amount": damage.Amount,
			},
		},
		Damage: damage,
	}
}

func Deal(damage Damage) {
	executeStep(NewBeforeDamageStep(&damage))
	totalDamage := damage.Amount - damage.Blocked
	if totalDamage < 0 {
		totalDamage = 0
	}

	fmt.Printf("Total damage after block: %d (original: %d, blocked: %d)\n", totalDamage, damage.Amount, damage.Blocked)
	applyDamage(damage, totalDamage)
	damage.Taken = totalDamage

	executeStep(NewAfterDamageStep(&damage))
}

func applyDamage(damage Damage, totalDamage int) {
	damage.Target.TakeDamage(totalDamage)
}
