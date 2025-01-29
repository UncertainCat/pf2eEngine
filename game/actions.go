package game

import (
	"fmt"
	dice "pf2eEngine/util"
)

type Attack struct {
	Attacker *Entity
	Defender *Entity
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
			StepType: BeforeAttack,
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
			StepType: AfterAttack,
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
	FreeAction   ActionType = "FREE_ACTION"
	Activity     ActionType = "ACTIVITY"
	Reaction     ActionType = "REACTION"
	EndOfTurn    ActionType = "END_OF_TURN"
)

type Action struct {
	Name        string
	Type        ActionType
	Cost        int
	Description string
	perform     func(gs *GameState, actor *Entity)
}

type StartActionStep struct {
	BaseStep
	Action Action
	Actor  *Entity
}

type EndActionStep struct {
	BaseStep
	Action Action
	Actor  *Entity
}

func EndTurnAction(gs *GameState, actor *Entity) Action {
	return Action{
		Name:    "End Turn",
		Type:    EndOfTurn,
		perform: func(gs *GameState, actor *Entity) {},
	}
}

func ExecuteAction(gs *GameState, actor *Entity, action Action) {
	if actor.ActionsRemaining < action.Cost {
		fmt.Printf("%s does not have enough actions to perform %s. Actions remaining: %d, action cost: %d.\n",
			actor.Name, action.Name, actor.ActionsRemaining, action.Cost)
		return
	}

	actor.SpendAction(action.Cost)
	fmt.Printf("%s used %d actions. Actions remaining: %d.\n",
		actor.Name, action.Cost, actor.ActionsRemaining)

	executeStep(gs, StartActionStep{
		BaseStep: BaseStep{StepType: StartTurn},
		Action:   action,
		Actor:    actor,
	}, fmt.Sprintf("%s starts the action: %s.", actor.Name, action.Name))

	action.perform(gs, actor)

	executeStep(gs, EndActionStep{
		BaseStep: BaseStep{StepType: EndTurn},
		Action:   action,
		Actor:    actor,
	}, fmt.Sprintf("%s completed the action: %s.", actor.Name, action.Name))
}

// PerformAttack encapsulates the full attack logic
func PerformAttack(gs *GameState, baseAttack BaseAttack, attacker *Entity, defender *Entity) {
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
		Bonus:    baseAttack.Bonus - attacker.MapCounter*5,
		Result:   roll + baseAttack.Bonus - attacker.MapCounter*5,
	}
	attack.Degree = calculateDegreeOfSuccess(roll, attack.Result, defender.AC)

	details := fmt.Sprintf(
		"Attack Details:\n\tAttacker: %s\n\tDefender: %s\n\tRoll: %d\n\tBonus: %d\n\tResult: %d\n\tDefender AC: %d\n\tDegree: %v",
		attacker.Name, defender.Name, roll, attack.Bonus, attack.Result, defender.AC, attack.Degree.String(),
	)

	damageRoll := baseAttack.RollDamage()

	damage := Damage{Source: attacker, Target: defender, Amount: damageRoll}
	switch attack.Degree {
	case CriticalSuccess:
		executeStep(gs, NewBeforeAttackStep(attack), fmt.Sprintf("%s has critically hit %s! Details:\n%s", attacker.Name, defender.Name, details))
		Deal(gs, damage.Double())
	case Success:
		executeStep(gs, NewBeforeAttackStep(attack), fmt.Sprintf("%s has hit %s. Details:\n%s", attacker.Name, defender.Name, details))
		Deal(gs, damage)
	case Failure:
		executeStep(gs, NewBeforeAttackStep(attack), fmt.Sprintf("%s has missed %s. Details:\n%s", attacker.Name, defender.Name, details))
	case CriticalFailure:
		executeStep(gs, NewBeforeAttackStep(attack), fmt.Sprintf("%s has critically missed %s! Details:\n%s", attacker.Name, defender.Name, details))
	}

	executeStep(gs, NewAfterAttackStep(attack), fmt.Sprintf("%s has finished attacking %s.", attacker.Name, defender.Name))
	attacker.MapCounter++
}
