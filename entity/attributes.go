package entity

// Attributes represent additional attributes for an entity
// This can be extended for more complex mechanics

type Attributes struct {
	Strength     int
	Dexterity    int
	Constitution int
	Intelligence int
	Wisdom       int
	Charisma     int
}

// NewAttributes creates a new Attributes instance with default values
func NewAttributes(strength, dexterity, constitution, intelligence, wisdom, charisma int) *Attributes {
	return &Attributes{
		Strength:     strength,
		Dexterity:    dexterity,
		Constitution: constitution,
		Intelligence: intelligence,
		Wisdom:       wisdom,
		Charisma:     charisma,
	}
}

// ApplyModifiers adjusts attributes dynamically based on external factors
func (a *Attributes) ApplyModifiers(modifiers map[string]int) {
	if mod, ok := modifiers["Strength"]; ok {
		a.Strength += mod
	}
	if mod, ok := modifiers["Dexterity"]; ok {
		a.Dexterity += mod
	}
	if mod, ok := modifiers["Constitution"]; ok {
		a.Constitution += mod
	}
	if mod, ok := modifiers["Intelligence"]; ok {
		a.Intelligence += mod
	}
	if mod, ok := modifiers["Wisdom"]; ok {
		a.Wisdom += mod
	}
	if mod, ok := modifiers["Charisma"]; ok {
		a.Charisma += mod
	}
}
