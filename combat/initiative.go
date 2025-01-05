package combat

import (
	"fmt"
	dice "pf2eEngine/util"
	"sort"

	"pf2eEngine/entity"
)

// RollInitiative determines initiative for all combatants and sorts them in descending order
func RollInitiative(entities []*entity.Entity) {
	for _, e := range entities {
		roll := dice.Roll(20) // Roll a d20 for initiative
		e.RollInitiative(roll)
		fmt.Printf("%s rolls initiative: %d\n", e.Name, e.Initiative)
	}

	// Sort entities by initiative (descending order)
	sort.Slice(entities, func(i, j int) bool {
		return entities[i].Initiative > entities[j].Initiative
	})

	fmt.Println("Initiative order:")
	for i, e := range entities {
		fmt.Printf("%d: %s (Initiative: %d)\n", i+1, e.Name, e.Initiative)
	}
}
