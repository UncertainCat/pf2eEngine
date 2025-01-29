package game

import (
	"encoding/json"
	"fmt"
	dice "pf2eEngine/util"
	"sort"
)

// GameState represents the state of the game, including turn order and current turn
// Logs are added to track events in both human-readable and JSON format
type GameState struct {
	Grid        *Grid
	Initiative  []*Entity
	CurrentTurn int
	Logs        []LogEntry
	StepHistory *StepHistory
}

type StepHistory struct {
	Steps []Step
}

func (sh *StepHistory) AddStep(step Step) {
	sh.Steps = append(sh.Steps, step)
}

type LogEntry struct {
	Message  string
	Metadata map[string]interface{}
	JSON     string
}

// NewGameState initializes a new game state with the given entities and grid
func NewGameState(spawns []Spawn, gridWidth, gridHeight int) *GameState {
	entities := []*Entity{}
	for _, spawn := range spawns {
		entities = append(entities, spawn.Unit)
	}
	gs := &GameState{
		CurrentTurn: 0,
		Grid:        NewGrid(gridWidth, gridHeight),
		Initiative:  entities,
		StepHistory: &StepHistory{},
	}
	// Place entities on the grid at their initial positions (e.g., in a row)
	for _, e := range spawns {
		gs.Grid.AddEntity(Position{X: e.Coordinates[0], Y: e.Coordinates[1]}, e.Unit)
	}
	gs.RollInitiative()
	return gs
}

type StartTurnStep struct {
	BaseStep
	Entity *Entity
}

type EndTurnStep struct {
	BaseStep
	Entity *Entity
}

func (gs *GameState) IsEntityTurn(e *Entity) bool {
	return gs.GetCurrentTurnEntity() == e
}

// RollInitiative determines initiative for all combatants and sorts them in descending order
func (gs *GameState) RollInitiative() {
	for _, e := range gs.Initiative {
		roll := dice.Roll(20)
		e.RollInitiative(roll)
		executeStep(gs, StartTurnStep{
			BaseStep: BaseStep{StepType: StartTurn},
			Entity:   e,
		}, fmt.Sprintf("%s rolls initiative: %d", e.Name, e.Initiative))
	}

	sort.Slice(gs.Initiative, func(i, j int) bool {
		return gs.Initiative[i].Initiative > gs.Initiative[j].Initiative
	})

	executeStep(gs, StartTurnStep{
		BaseStep: BaseStep{StepType: StartTurn},
		Entity:   nil,
	}, "Initiative order determined")
}

func (gs *GameState) getInitiativeOrder() []string {
	order := []string{}
	for _, e := range gs.Initiative {
		order = append(order, fmt.Sprintf("%s (Initiative: %d)", e.Name, e.Initiative))
	}
	return order
}

// IsCombatOver checks if the combat has ended
func (gs *GameState) IsCombatOver() bool {
	var faction Faction
	firstLivingEntity := true

	for _, e := range gs.Initiative {
		if e.IsAlive() {
			if firstLivingEntity {
				faction = e.Faction
				firstLivingEntity = false
			} else if e.Faction != faction {
				return false
			}
		}
	}
	return true
}

// GetCurrentTurnEntity returns the entity whose turn it is
func (gs *GameState) GetCurrentTurnEntity() *Entity {
	return gs.Initiative[gs.CurrentTurn]
}

// LogEvent logs both human-readable and JSON-format logs
func (gs *GameState) LogEvent(message string, metadata map[string]interface{}) {
	logEntry := LogEntry{
		Message:  message,
		Metadata: metadata,
		JSON:     toJSON(metadata),
	}
	gs.Logs = append(gs.Logs, logEntry)
	fmt.Println(message)
}

func toJSON(metadata map[string]interface{}) string {
	jsonData, _ := json.Marshal(metadata)
	return string(jsonData)
}

// StartTurn triggers a StartTurnStep and logs the event
func (gs *GameState) StartTurn() {
	currentEntity := gs.GetCurrentTurnEntity()
	startStep := StartTurnStep{
		BaseStep: BaseStep{StepType: StartTurn},
		Entity:   currentEntity,
	}
	executeStep(gs, startStep, fmt.Sprintf("%s's turn begins", currentEntity.Name))

	currentEntity.ResetTurnResources()
}

// EndTurn triggers an EndTurnStep and logs the event
func (gs *GameState) EndTurn() {
	currentEntity := gs.GetCurrentTurnEntity()
	endStep := EndTurnStep{
		BaseStep: BaseStep{StepType: EndTurn},
		Entity:   currentEntity,
	}
	executeStep(gs, endStep, fmt.Sprintf("%s's turn ends", currentEntity.Name))

	for {
		gs.CurrentTurn = (gs.CurrentTurn + 1) % len(gs.Initiative)
		if gs.Initiative[gs.CurrentTurn].IsAlive() {
			break
		}
	}
}

// GetWinner returns the winning entity, if combat is over
func (gs *GameState) GetWinner() *Entity {
	if !gs.IsCombatOver() {
		return nil
	}
	for _, e := range gs.Initiative {
		if e.IsAlive() {
			return e
		}
	}
	return nil
}

type Spawn struct {
	Unit        *Entity
	Coordinates [2]int
}

func StartCombat(spawns []Spawn, gridWidth, gridHeight int) *GameState {
	gs := NewGameState(spawns, gridWidth, gridHeight)
	go func() {
		RunCombat(gs)
	}()
	return gs
}

func RunCombat(gs *GameState) {
	for !gs.IsCombatOver() {
		currentEntity := gs.GetCurrentTurnEntity()
		gs.StartTurn()
		action := currentEntity.Controller.NextAction(gs, currentEntity)
		for action.Type != EndOfTurn {
			ExecuteAction(gs, currentEntity, action)
			action = currentEntity.Controller.NextAction(gs, currentEntity)
		}
		gs.EndTurn()
	}

	winner := gs.GetWinner()
	if winner != nil {
		gs.LogEvent(fmt.Sprintf("%s wins the combat!", winner.Name), map[string]interface{}{
			"winner": winner.Name,
		})
	} else {
		gs.LogEvent("Combat ends with no winner.", map[string]interface{}{})
	}
}
