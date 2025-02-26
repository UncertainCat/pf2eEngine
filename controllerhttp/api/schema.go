package api

import (
	"github.com/google/uuid"
	"time"
)

// Version identifier for the API
const (
	CurrentVersion = "v1"
)

// ApiSchema contains all the model definitions for the public API
// This provides a clear contract between frontend and backend
// and allows versioning of the API in the future

// Common models shared across different events

// EntityRef is a lightweight reference to an entity
type EntityRef struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// EventBase is the base structure for all events
type EventBase struct {
	Type      string                 `json:"type"`
	Version   string                 `json:"version"`
	Timestamp time.Time              `json:"timestamp"`
	Message   string                 `json:"message,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// GameEvent represents a game event sent to the frontend
type GameEvent struct {
	EventBase
	Data interface{} `json:"data,omitempty"`
}

// ActionCardRef represents a reference to an action card
type ActionCardRef struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	ActionCost  int       `json:"actionCost"`
	Type        string    `json:"type,omitempty"`
}

// EntityState represents the complete state of an entity
type EntityState struct {
	ID                 uuid.UUID       `json:"id"`
	Name               string          `json:"name"`
	HP                 int             `json:"hp"`
	MaxHP              int             `json:"maxHp"` // Added to ensure frontend knows the max HP
	AC                 int             `json:"ac"`
	ActionsRemaining   int             `json:"actionsRemaining"`
	ReactionsRemaining int             `json:"reactionsRemaining"`
	Faction            string          `json:"faction"`
	ActionCards        []ActionCardRef `json:"actionCards,omitempty"`
	Position           [2]int          `json:"position,omitempty"`
}

// GameState represents the entire game state
type GameState struct {
	Entities    []EntityState `json:"entities"`
	CurrentTurn *uuid.UUID    `json:"currentTurn,omitempty"`
	GridWidth   int           `json:"gridWidth"`
	GridHeight  int           `json:"gridHeight"`
	Round       int           `json:"round"`
}

// CommandRequest represents a command sent from the frontend to the backend
type CommandRequest struct {
	EntityID     uuid.UUID              `json:"entity_id"`
	ActionCardID uuid.UUID              `json:"action_card_id"`
	Params       map[string]interface{} `json:"params"`
}

// CommandResponse represents a response to a command
type CommandResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// Specialized event data structures

// AttackEventData represents an attack event
type AttackEventData struct {
	Attacker EntityRef `json:"attacker"`
	Defender EntityRef `json:"defender"`
	Roll     int       `json:"roll"`
	Result   int       `json:"result"`
	Degree   string    `json:"degree,omitempty"`
}

// DamageEventData represents a damage event
type DamageEventData struct {
	Source  EntityRef `json:"source"`
	Target  EntityRef `json:"target"`
	Amount  int       `json:"amount,omitempty"`
	Type    string    `json:"type,omitempty"`
	Blocked int       `json:"blocked,omitempty"`
	Taken   int       `json:"taken,omitempty"`
}

// TurnEventData represents a turn event
type TurnEventData struct {
	Entity EntityRef `json:"entity"`
}

// GameSetupData represents initial game setup data
type GameSetupData struct {
	GridWidth  int           `json:"gridWidth"`
	GridHeight int           `json:"gridHeight"`
	Entities   []EntityState `json:"entities"`
}

// EventType constants define all possible event types
// This allows for explicit typing of events in the frontend
const (
	// Core event types
	EventTypeInfo           = "INFO"
	EventTypeGameSetup      = "GAME_SETUP"
	EventTypeGameState      = "GAME_STATE"
	EventTypeAttack         = "ATTACK"
	EventTypeAttackResult   = "ATTACK_RESULT"
	EventTypeDamage         = "DAMAGE"
	EventTypeDamageResult   = "DAMAGE_RESULT"
	EventTypeTurnStart      = "TURN_START"
	EventTypeTurnEnd        = "TURN_END"
	EventTypeRoundStart     = "ROUND_START"
	EventTypeRoundEnd       = "ROUND_END"
	EventTypeEntityMove     = "ENTITY_MOVE"
	EventTypeEntityStatus   = "ENTITY_STATUS"
	EventTypeActionComplete = "ACTION_COMPLETE"
)