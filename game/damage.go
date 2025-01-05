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
}

// CreateDamage function
func CreateDamage(ctx Context, source *entity.Entity, target *entity.Entity, amount int) Damage {
	damage := Damage{
		Source:  source,
		Target:  target,
		Amount:  amount,
		Blocked: 0,
	}
	ctx["damage"] = damage
	return damage
}

// Deal function
func Deal(ctx Context, damage Damage) {
	before(damage, ctx)

	totalDamage := damage.Amount - damage.Blocked
	if totalDamage < 0 {
		totalDamage = 0
	}

	fmt.Printf("Total damage after block: %d (original: %d, blocked: %d)\n", totalDamage, damage.Amount, damage.Blocked)
	applyDamage(damage, totalDamage)

	after(damage, ctx)
}

// applyDamage helper
func applyDamage(damage Damage, totalDamage int) {
	damage.Target.TakeDamage(totalDamage)
	fmt.Printf("%s takes %d damage! Remaining HP: %d\n", damage.Target.Name, totalDamage, damage.Target.HP)
}

// before and after steps
func before(damage Damage, ctx Context) {
	trigger(ctx, "beforeDamage")
}

func after(damage Damage, ctx Context) {
	trigger(ctx, "afterDamage")
}
