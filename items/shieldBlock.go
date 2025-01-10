package items

import (
	"fmt"
	"pf2eEngine/entity"
	"pf2eEngine/game"
)

type ShieldBlock struct {
	Owner *entity.Entity
}

// Enforce ShieldBlock implements the Trigger interface
var _ game.Trigger = ShieldBlock{}

func (trigger ShieldBlock) Priority() int {
	return 10 // Example priority: higher values execute earlier
}

func (trigger ShieldBlock) Condition(step game.Step) bool {
	if damageStep, ok := step.(game.BeforeDamageStep); ok {
		return damageStep.Damage.Target != nil && damageStep.Damage.Target == trigger.Owner
	}
	return false
}

func (trigger ShieldBlock) Execute(step game.Step) {
	if !trigger.Owner.UseReaction() {
		fmt.Printf("%s has no reactions remaining to block!\n", trigger.Owner.Name)
		return
	}

	if damageStep, ok := step.(game.BeforeDamageStep); ok {
		damageStep.Damage.Blocked += 5
		fmt.Printf("%s uses Shield Block to block 5 damage.\n", trigger.Owner.Name)
	}
}
