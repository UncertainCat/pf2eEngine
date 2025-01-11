package game

type DegreeOfSuccess int

const (
	CriticalFailure DegreeOfSuccess = iota
	Failure
	Success
	CriticalSuccess
)

func (d DegreeOfSuccess) String() string {
	switch d {
	case CriticalFailure:
		return "Critical Failure"
	case Failure:
		return "Failure"
	case Success:
		return "Success"
	case CriticalSuccess:
		return "Critical Success"
	default:
		return "Unknown"
	}
}

func calculateDegreeOfSuccess(roll int, total int, dc int) DegreeOfSuccess {
	if roll == 20 {
		if total >= dc+10 {
			return CriticalSuccess
		}
		if total >= dc {
			return CriticalSuccess
		}
		return Success
	}
	if roll == 1 {
		if total >= dc+10 {
			return Success
		}
		return Failure
	}
	if total >= dc+10 {
		return CriticalSuccess
	}
	if total >= dc {
		return Success
	}
	if total <= dc-10 {
		return CriticalFailure
	}
	return Failure
}
