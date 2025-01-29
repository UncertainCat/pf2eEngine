package game

import (
	"fmt"
	dice "pf2eEngine/util"
)

type Damage struct {
	Source  *Entity
	Target  *Entity
	Amount  map[DamageType]DamageAmount
	Blocked int
	Taken   int
}

func (d Damage) Double() Damage {
	for k, v := range d.Amount {
		d.Amount[k] = DamageAmount{Amount: v.Amount * 2, Type: v.Type}
	}
	return d
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
			StepType: BeforeDamage,
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
			StepType: AfterDamage,
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
	totalDamage := 0
	for _, dr := range damage.Amount {
		dealDamage := dr.Amount - damage.Blocked
		if dealDamage < 0 {
			dealDamage = 0
		}
		totalDamage += dealDamage
		applyDamage(damage, dealDamage)
	}
	damage.Taken = totalDamage

	executeStep(gs, NewAfterDamageStep(&damage), fmt.Sprintf("%s dealt %d damage to %s.", damage.Source.Name, damage.Taken, damage.Target.Name))
}

func applyDamage(damage Damage, totalDamage int) {
	damage.Target.TakeDamage(totalDamage)
}

type DamageType string

const (
	Bludgeoning DamageType = "BLUDGEONING"
	Piercing    DamageType = "PIERCING"
	Slashing    DamageType = "SLASHING"
)

type DamageRoll struct {
	Die   int
	Count int
	Bonus int
	Type  DamageType
}

func (dr DamageRoll) Roll() DamageAmount {
	amount := dr.Bonus
	for i := 0; i < dr.Count; i++ {
		amount += dice.Roll(dr.Die)
	}
	return DamageAmount{Amount: amount, Type: dr.Type}
}

type DamageAmount struct {
	Amount int
	Type   DamageType
}

type BaseAttack struct {
	Damage []DamageRoll
	Bonus  int
}

func (ba BaseAttack) RollDamage() map[DamageType]DamageAmount {
	damage := map[DamageType]DamageAmount{}
	for _, dr := range ba.Damage {
		amount := dr.Roll()
		damage[amount.Type] = amount
	}
	return damage
}
