package game

import (
	"encoding/json"
	"fmt"
	dice "pf2eEngine/util"
	"sort"
	"time"
)

// StepCallback is a function type for callbacks when steps occur
type StepCallback func(step interface{}, message string)

// GameState represents the state of the game, including turn order and current turn
// Logs are added to track events in both human-readable and JSON format
type GameState struct {
	Grid                *Grid
	Initiative          []*Entity
	CurrentTurn         int
	Logs                []LogEntry
	StepHistory         *StepHistory
	StepCallback        StepCallback
	InitialEntities     []*Entity    // Copy of initial entities for resetting
	InitialEntityPos    map[string]Position // Initial positions of entities
	InitialCurrentTurn  int          // Initial current turn
}

type StepHistory struct {
	Steps []Step
}

func (sh *StepHistory) AddStep(step Step) {
	sh.Steps = append(sh.Steps, step)
}

func (sh *StepHistory) GetSteps() []Step {
	return sh.Steps
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
	
	// Initialize game state
	gs := &GameState{
		CurrentTurn: 0,
		Grid:        NewGrid(gridWidth, gridHeight),
		Initiative:  entities,
		StepHistory: &StepHistory{},
		InitialEntityPos: make(map[string]Position), // Using ID string as key
	}
	
	// Save initial positions as we place entities
	for _, e := range spawns {
		pos := Position{X: e.Coordinates[0], Y: e.Coordinates[1]}
		gs.Grid.AddEntity(pos, e.Unit)
		gs.InitialEntityPos[e.Unit.Id.String()] = pos
	}
	
	// Roll initiative to determine turn order
	gs.RollInitiative()
	
	// Save initial state AFTER initiative is rolled
	gs.InitialCurrentTurn = gs.CurrentTurn
	
	// Create deep copies of all entities to save initial state
	gs.InitialEntities = make([]*Entity, len(gs.Initiative))
	for i, entity := range gs.Initiative {
		// Deep copy each entity
		copiedEntity := &Entity{
			Id:                 entity.Id,
			Name:               entity.Name,
			HP:                 entity.HP,
			MaxHP:              entity.HP, // Store original HP as MaxHP
			AC:                 entity.AC,
			ActionsRemaining:   entity.ActionsRemaining,
			ReactionsRemaining: entity.ReactionsRemaining,
			Faction:            entity.Faction,
			// Deep copy action cards array
			ActionCards:        make([]*ActionCard, len(entity.ActionCards)),
		}
		
		// Copy action cards
		copy(copiedEntity.ActionCards, entity.ActionCards)
		
		gs.InitialEntities[i] = copiedEntity
	}
	
	// Log initial state
	gs.LogEvent("Game initialized", map[string]interface{}{
		"grid_width": gridWidth,
		"grid_height": gridHeight,
		"entity_count": len(entities),
	})
	
	// Update all entities with their MaxHP
	for _, entity := range gs.Initiative {
		entity.MaxHP = entity.HP
	}
	
	return gs
}

type StartTurnStep struct {
	BaseStep
	Entity *Entity
}

// GetInitialState creates a new GameState instance with the initial state
func (gs *GameState) GetInitialState() *GameState {
	// Create a new grid
	initialGrid := NewGrid(gs.Grid.Width, gs.Grid.Height)
	
	// Create a new state
	initialState := &GameState{
		Grid:         initialGrid,
		Initiative:   make([]*Entity, len(gs.InitialEntities)),
		CurrentTurn:  gs.InitialCurrentTurn,
		Logs:         []LogEntry{},
		StepHistory:  &StepHistory{},
		InitialEntityPos: gs.InitialEntityPos,
	}
	
	// Deep copy all initial entities
	for i, entity := range gs.InitialEntities {
		copiedEntity := &Entity{
			Id:                 entity.Id,
			Name:               entity.Name,
			HP:                 entity.HP,
			MaxHP:              entity.MaxHP,
			AC:                 entity.AC,
			ActionsRemaining:   entity.ActionsRemaining,
			ReactionsRemaining: entity.ReactionsRemaining,
			Faction:            entity.Faction,
			ActionCards:        make([]*ActionCard, len(entity.ActionCards)),
		}
		
		// Copy action cards
		copy(copiedEntity.ActionCards, entity.ActionCards)
		
		initialState.Initiative[i] = copiedEntity
		
		// Place entity on the grid at its initial position
		if pos, ok := gs.InitialEntityPos[entity.Id.String()]; ok {
			initialGrid.AddEntity(pos, copiedEntity)
		}
	}
	
	// Copy initial entities and positions to the new state
	initialState.InitialEntities = make([]*Entity, len(gs.InitialEntities))
	for i, entity := range gs.InitialEntities {
		initialState.InitialEntities[i] = entity
	}
	
	return initialState
}

type EndTurnStep struct {
	BaseStep
	Entity *Entity
}

func (gs *GameState) IsEntityTurn(e *Entity) bool {
	return gs.GetCurrentTurnEntity() == e
}

func (e *EndTurnStep) Type() StepType {
	return EndTurn
}

func (e *EndTurnStep) Metadata() map[string]interface{} {
	return map[string]interface{}{
		"entity_id":   e.Entity.Id.String(),
		"entity_name": e.Entity.Name,
	}
}

func (e *StartTurnStep) Type() StepType {
	return StartTurn
}

func (e *StartTurnStep) Metadata() map[string]interface{} {
	return map[string]interface{}{
		"entity_id":   e.Entity.Id.String(),
		"entity_name": e.Entity.Name,
	}
}

// RollInitiative rolls initiative for each entity and sorts the initiative order
func (gs *GameState) RollInitiative() {
	// Roll initiative for each entity
	for _, entity := range gs.Initiative {
		roll := dice.Roll(20)
		fmt.Printf("%s rolls initiative: %d\n", entity.Name, roll)
		entity.Initiative = roll
	}

	// Sort by initiative score in descending order
	sort.Slice(gs.Initiative, func(i, j int) bool {
		return gs.Initiative[i].Initiative > gs.Initiative[j].Initiative
	})

	fmt.Println("Initiative order determined")
}

func (gs *GameState) LogEvent(message string, metadata map[string]interface{}) {
	jsonData, _ := json.Marshal(metadata)
	logEntry := LogEntry{
		Message:  message,
		Metadata: metadata,
		JSON:     string(jsonData),
	}
	gs.Logs = append(gs.Logs, logEntry)
}

// StartCombat signals that combat is starting and logs it
func (gs *GameState) StartCombat() {
	gs.LogEvent("Combat started", map[string]interface{}{
		"time": time.Now().String(),
	})
}

// NextTurn advances to the next entity in the initiative order and returns the new current entity
func (gs *GameState) NextTurn() *Entity {
	// Get the current entity
	entity := gs.GetCurrentTurnEntity()
	
	// End its turn if it's a valid entity
	if entity != nil {
		// Create end turn step
		endTurnStep := &EndTurnStep{
			Entity: entity,
		}
		
		// Add to step history
		gs.StepHistory.AddStep(endTurnStep)
		
		// Tell controller via callback
		if gs.StepCallback != nil {
			gs.StepCallback(endTurnStep, fmt.Sprintf("%s's turn ends", entity.Name))
		}
		
		// Log the end of the turn
		gs.LogEvent(fmt.Sprintf("%s's turn ends", entity.Name), map[string]interface{}{
			"entity_id":   entity.Id.String(),
			"entity_name": entity.Name,
		})
	}
	
	// Advance to the next entity in initiative order
	gs.CurrentTurn = (gs.CurrentTurn + 1) % len(gs.Initiative)
	
	// If all entities are dead except one, end combat
	aliveCount := 0
	var lastAlive *Entity
	for _, e := range gs.Initiative {
		if e.IsAlive() {
			aliveCount++
			lastAlive = e
		}
	}
	
	// If only one entity is alive, end combat
	if aliveCount <= 1 && lastAlive != nil {
		gs.LogEvent(fmt.Sprintf("%s wins the combat!", lastAlive.Name), map[string]interface{}{
			"winner":      lastAlive.Name,
			"winner_id":   lastAlive.Id.String(),
			"winner_hp":   lastAlive.HP,
			"combat_over": true,
		})
		
		// Output the winner
		fmt.Printf("%s wins the combat!\n", lastAlive.Name)
		
		// Return nil to indicate combat is over
		return nil
	}
	
	// Reset actions for new entity's turn
	entity = gs.GetCurrentTurnEntity()
	entity.ResetTurnResources()
	
	// Skip dead entities
	for !entity.IsAlive() {
		gs.CurrentTurn = (gs.CurrentTurn + 1) % len(gs.Initiative)
		entity = gs.GetCurrentTurnEntity()
		entity.ResetTurnResources()
	}
	
	// Create start turn step
	startTurnStep := &StartTurnStep{
		Entity: entity,
	}
	
	// Add to step history
	gs.StepHistory.AddStep(startTurnStep)
	
	// Tell controller via callback
	if gs.StepCallback != nil {
		gs.StepCallback(startTurnStep, fmt.Sprintf("%s's turn begins", entity.Name))
	}
	
	// Log the start of the turn
	gs.LogEvent(fmt.Sprintf("%s's turn begins", entity.Name), map[string]interface{}{
		"entity_id":   entity.Id.String(),
		"entity_name": entity.Name,
	})
	
	fmt.Printf("%s's turn begins\n", entity.Name)
	
	return entity
}

// GetCurrentTurnEntity returns a pointer to the entity whose turn it currently is
func (gs *GameState) GetCurrentTurnEntity() *Entity {
	if len(gs.Initiative) == 0 {
		return nil
	}
	return gs.Initiative[gs.CurrentTurn]
}

type Spawn struct {
	Unit        *Entity
	Coordinates [2]int
}

// Creates a spawn specification for an entity
func NewSpawn(entity *Entity, x, y int) Spawn {
	return Spawn{
		Unit:        entity,
		Coordinates: [2]int{x, y},
	}
}