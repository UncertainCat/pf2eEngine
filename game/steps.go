package game

import (
	"sort"
)

type StepType string

const (
	BeforeDamage StepType = "BEFORE_DAMAGE"
	AfterDamage  StepType = "AFTER_DAMAGE"
	Movement     StepType = "MOVEMENT"
	Spellcast    StepType = "SPELLCAST"
)

type Step interface {
	Type() StepType
	Metadata() map[string]interface{}
}

type BaseStep struct {
	stepType StepType
	metadata map[string]interface{}
}

func (s BaseStep) Type() StepType {
	return s.stepType
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

func executeStep(step Step) {
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
