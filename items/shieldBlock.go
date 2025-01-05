package items

import (
	"fmt"
	"pf2eEngine/game"
)

func ShieldBlock(ctx game.Context) bool {
	damage, ok := ctx["damage"].(game.Damage)
	if !ok {
		return false
	}

	source := damage.Source
	if ctx["currentTurn"] != source {
		return false
	}

	damage.Blocked += 5
	ctx["damage"] = damage

	fmt.Printf("%s uses Shield Block! Blocking %d damage.\n", damage.Target.Name, 5)
	return true
}
