package game

import (
	"fmt"
	"pf2eEngine/entity"
	dice "pf2eEngine/util"
)

type Attack struct {
	Attacker *entity.Entity
	Defender *entity.Entity
	Roll     int
	Bonus    int
	Result   int
	Degree   DegreeOfSuccess
}

type BeforeAttackStep struct {
	BaseStep
	Attack *Attack
}

type AfterAttackStep struct {
	BaseStep
	Attack *Attack
}

func NewBeforeAttackStep(attack *Attack) BeforeAttackStep {
	return BeforeAttackStep{
		BaseStep: BaseStep{
			stepType: BeforeAttack,
			metadata: map[string]interface{}{
				"Attacker": attack.Attacker.Name,
				"Defender": attack.Defender.Name,
			},
		},
		Attack: attack,
	}
}

func NewAfterAttackStep(attack *Attack) AfterAttackStep {
	return AfterAttackStep{
		BaseStep: BaseStep{
			stepType: AfterAttack,
			metadata: map[string]interface{}{
				"Attacker": attack.Attacker.Name,
				"Defender": attack.Defender.Name,
				"Result":   attack.Result,
				"Degree":   attack.Degree,
			},
		},
		Attack: attack,
	}
}

// ActionType defines the category of an action
type ActionType string

const (
	SingleAction ActionType = "SINGLE_ACTION"
	Activity     ActionType = "ACTIVITY"
	Reaction     ActionType = "REACTION"
)

type Action struct {
	Name        string
	Type        ActionType
	Cost        int
	Description string
	Perform     func(gs *GameState, actor *entity.Entity)
}

type StartActionStep struct {
	BaseStep
	Action Action
	Actor  *entity.Entity
}

type EndActionStep struct {
	BaseStep
	Action Action
	Actor  *entity.Entity
}

func ExecuteAction(gs *GameState, actor *entity.Entity, action Action) {
	if actor.ActionsRemaining < action.Cost {
		fmt.Printf("%s does not have enough actions to perform %s.\n", actor.Name, action.Name)
		return
	}

	executeStep(gs, StartActionStep{
		BaseStep: BaseStep{stepType: StartTurn},
		Action:   action,
		Actor:    actor,
	}, fmt.Sprintf("%s starts the action: %s.", actor.Name, action.Name))

	actor.UseAction(action.Cost)
	action.Perform(gs, actor)

	executeStep(gs, EndActionStep{
		BaseStep: BaseStep{stepType: EndTurn},
		Action:   action,
		Actor:    actor,
	}, fmt.Sprintf("%s completed the action: %s.", actor.Name, action.Name))
}

type StrikeAction struct {
	Target *entity.Entity
	Action
}

func Strike(target *entity.Entity) Action {
	return Action{
		Name:        "Strike",
		Type:        SingleAction,
		Cost:        1,
		Description: "A basic melee attack.",
		Perform: func(gs *GameState, actor *entity.Entity) {
			PerformAttack(gs, actor, target)
		},
	}
}

// PerformAttack encapsulates the full attack logic
func PerformAttack(gs *GameState, attacker *entity.Entity, defender *entity.Entity) {
	attackerPos := gs.Grid.GetEntityPosition(attacker)
	defenderPos := gs.Grid.GetEntityPosition(defender)

	if !gs.Grid.AreAdjacent(attackerPos, defenderPos) {
		fmt.Printf("%s cannot attack %s; they are not adjacent.\n", attacker.Name, defender.Name)
		return
	}

	roll := dice.Roll(20)
	attack := &Attack{
		Attacker: attacker,
		Defender: defender,
		Roll:     roll,
		Bonus:    attacker.AttackBonus - attacker.MapCounter*5,
		Result:   roll + attacker.AttackBonus - attacker.MapCounter*5,
	}
	attack.Degree = calculateDegreeOfSuccess(roll, attack.Result, defender.AC)

	details := fmt.Sprintf(
		"Attack Details:\n\tAttacker: %s\n\tDefender: %s\n\tRoll: %d\n\tBonus: %d\n\tResult: %d\n\tDefender AC: %d\n\tDegree: %v",
		attacker.Name, defender.Name, roll, attack.Bonus, attack.Result, defender.AC, attack.Degree.String(),
	)

	switch attack.Degree {
	case CriticalSuccess:
		damage := (dice.Roll(8) + attacker.DamageBonus) * 2
		executeStep(gs, NewBeforeAttackStep(attack), fmt.Sprintf("%s has critically hit %s! Details:\n%s", attacker.Name, defender.Name, details))
		Deal(gs, Damage{Source: attacker, Target: defender, Amount: damage})
	case Success:
		damage := dice.Roll(8) + attacker.DamageBonus
		executeStep(gs, NewBeforeAttackStep(attack), fmt.Sprintf("%s has hit %s. Details:\n%s", attacker.Name, defender.Name, details))
		Deal(gs, Damage{Source: attacker, Target: defender, Amount: damage})
	case Failure:
		executeStep(gs, NewBeforeAttackStep(attack), fmt.Sprintf("%s has missed %s. Details:\n%s", attacker.Name, defender.Name, details))
	case CriticalFailure:
		executeStep(gs, NewBeforeAttackStep(attack), fmt.Sprintf("%s has critically missed %s! Details:\n%s", attacker.Name, defender.Name, details))
	}

	executeStep(gs, NewAfterAttackStep(attack), fmt.Sprintf("%s has finished attacking %s.", attacker.Name, defender.Name))
	attacker.MapCounter++
}
