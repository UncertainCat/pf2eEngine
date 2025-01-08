package game

type Trigger interface {
	Execute(Step)
	Condition(Step) bool
}

var triggers map[StepType][]Trigger

func RegisterTrigger(trigger Trigger, t StepType) {
	if triggers == nil {
		triggers = make(map[StepType][]Trigger)
	}
	triggers[t] = append(triggers[t], trigger)
}

func executeStep(step Step) {
	for _, t := range triggers[step.Type()] {
		if t.Condition(step) {
			t.Execute(step)
		}
	}
}

type StepType string

type Step interface {
	Type() StepType
}
