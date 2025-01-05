package game

type Context map[string]interface{}

var triggers = make(map[string][]func(Context) bool)

func RegisterTrigger(stepName string, trigger func(Context) bool) {
	triggers[stepName] = append(triggers[stepName], trigger)
}

func trigger(ctx Context, stepName string) {
	if stepTriggers, found := triggers[stepName]; found {
		for _, t := range stepTriggers {
			if t(ctx) {
				// Trigger logic is executed
			}
		}
	}
}
