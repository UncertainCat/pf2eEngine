// shieldBlock.go
package items

import (
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
	if damageStep, ok := step.(game.BeforeDamageStep); ok {
		damageStep.Damage.Blocked += 5
	}
}
