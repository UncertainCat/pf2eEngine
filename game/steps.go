package game

import (
	"sort"
)

type StepType string

const (
	BeforeDamage StepType = "BEFORE_DAMAGE"
	AfterDamage  StepType = "AFTER_DAMAGE"
	BeforeAttack StepType = "BEFORE_ATTACK"
	AfterAttack  StepType = "AFTER_ATTACK"
	StartTurn    StepType = "START_TURN"
	EndTurn      StepType = "END_TURN"
)

type Step interface {
	Type() StepType
	Metadata() map[string]interface{}
}

type BaseStep struct {
	StepType StepType
	metadata map[string]interface{}
}

func (s BaseStep) Type() StepType {
	return s.StepType
}

func (s BaseStep) Metadata() map[string]interface{} {
	return s.metadata
}

type Trigger interface {
	Priority() int
	Condition(Step) bool
	Execute(Step)
}

type BaseTrigger struct {
	priority int
}

func (t BaseTrigger) Priority() int {
	return t.priority
}

var triggers map[StepType][]Trigger

func RegisterTrigger(trigger Trigger, t StepType) {
	if triggers == nil {
		triggers = make(map[StepType][]Trigger)
	}
	triggers[t] = append(triggers[t], trigger)
}

func executeStep(gs *GameState, step Step, logMessage string) {
	gs.LogEvent(logMessage, step.Metadata())
	gs.StepHistory.AddStep(step)

	if triggersForStep, ok := triggers[step.Type()]; ok {
		sort.Slice(triggersForStep, func(i, j int) bool {
			return triggersForStep[i].Priority() > triggersForStep[j].Priority()
		})

		for _, trigger := range triggersForStep {
			if trigger.Condition(step) {
				trigger.Execute(step)
			}
		}
	}
}
