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

func (trigger ShieldBlock) Condition(step game.Step) bool {
	damage, ok := step.(game.BeforeDamageStep)
	if ok && damage.Target != nil && damage.Target == trigger.Owner {
		return true
	}
	return false
}

func (trigger ShieldBlock) Execute(step game.Step) {
	damage, ok := step.(game.BeforeDamageStep)
	if !ok {
		return
	}
	damage.Blocked += 5
	return
}
