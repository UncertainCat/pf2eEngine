package dice

import "math/rand"

// Roll simulates a dice roll of n sides
func Roll(sides int) int {
	return rand.Intn(sides) + 1
}

// RollMultiple rolls multiple dice of the same type and returns the total
func RollMultiple(sides, count int) int {
	total := 0
	for i := 0; i < count; i++ {
		total += Roll(sides)
	}
	return total
}

// RollWithModifier rolls a dice and adds a modifier to the result
func RollWithModifier(sides, modifier int) int {
	return Roll(sides) + modifier
}
