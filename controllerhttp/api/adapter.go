package api

import (
	"pf2eEngine/game"
	"time"
)

// StepToEvent converts internal game.Step objects to frontend-friendly API events
func StepToEvent(step game.Step, message string) GameEvent {
	event := GameEvent{
		EventBase: EventBase{
			Type:      mapStepTypeToEventType(step.Type()),
			Version:   CurrentVersion,
			Timestamp: time.Now(),
			Message:   message,
			Metadata:  step.Metadata(),
		},
	}

	// Convert specific step types to structured events with dedicated data
	switch s := step.(type) {
	case game.BeforeAttackStep:
		event.Data = AttackEventData{
			Attacker: EntityRef{
				ID:   s.Attack.Attacker.Id,
				Name: s.Attack.Attacker.Name,
			},
			Defender: EntityRef{
				ID:   s.Attack.Defender.Id,
				Name: s.Attack.Defender.Name,
			},
			Roll:   s.Attack.Roll,
			Result: s.Attack.Result,
		}

	case game.AfterAttackStep:
		event.Data = AttackEventData{
			Attacker: EntityRef{
				ID:   s.Attack.Attacker.Id,
				Name: s.Attack.Attacker.Name,
			},
			Defender: EntityRef{
				ID:   s.Attack.Defender.Id,
				Name: s.Attack.Defender.Name,
			},
			Degree: s.Attack.Degree.String(),
		}

	case game.BeforeDamageStep:
		event.Data = DamageEventData{
			Source: EntityRef{
				ID:   s.Damage.Source.Id,
				Name: s.Damage.Source.Name,
			},
			Target: EntityRef{
				ID:   s.Damage.Target.Id,
				Name: s.Damage.Target.Name,
			},
			// Sum up all damage amounts for a simplified representation
			Amount: sumDamageAmount(s.Damage.Amount),
		}

	case game.AfterDamageStep:
		event.Data = DamageEventData{
			Source: EntityRef{
				ID:   s.Damage.Source.Id,
				Name: s.Damage.Source.Name,
			},
			Target: EntityRef{
				ID:   s.Damage.Target.Id,
				Name: s.Damage.Target.Name,
			},
			Blocked: s.Damage.Blocked,
			Taken:   s.Damage.Taken,
		}

	case game.StartTurnStep:
		if s.Entity != nil {
			event.Data = TurnEventData{
				Entity: EntityRef{
					ID:   s.Entity.Id,
					Name: s.Entity.Name,
				},
			}
		}

	case game.EndTurnStep:
		if s.Entity != nil {
			event.Data = TurnEventData{
				Entity: EntityRef{
					ID:   s.Entity.Id,
					Name: s.Entity.Name,
				},
			}
		}
	}

	return event
}

// Sum up damage amounts from multiple damage types
func sumDamageAmount(damageMap map[game.DamageType]game.DamageAmount) int {
	total := 0
	for _, amount := range damageMap {
		total += amount.Amount
	}
	return total
}

// Maps internal step types to frontend event types
func mapStepTypeToEventType(stepType game.StepType) string {
	switch stepType {
	case game.BeforeAttack:
		return EventTypeAttack
	case game.AfterAttack:
		return EventTypeAttackResult
	case game.BeforeDamage:
		return EventTypeDamage
	case game.AfterDamage:
		return EventTypeDamageResult
	case game.StartTurn:
		return EventTypeTurnStart
	case game.EndTurn:
		return EventTypeTurnEnd
	default:
		return EventTypeInfo
	}
}

// GameStateToAPIState converts the internal game state to the API representation
func GameStateToAPIState(gs *game.GameState) GameState {
	apiState := GameState{
		Entities:   make([]EntityState, 0),
		GridWidth:  gs.Grid.Width,
		GridHeight: gs.Grid.Height,
		Round:      0, // Round is not explicitly tracked in the current game state
	}

	// Get the current entity ID if there is one
	currentEntity := gs.GetCurrentTurnEntity()
	if currentEntity != nil {
		currentEntityID := currentEntity.Id
		apiState.CurrentTurn = &currentEntityID
	}

	// Convert all entities from the initiative list
	for _, entity := range gs.Initiative {
		apiEntity := EntityToAPIEntity(entity, gs.Grid)
		apiState.Entities = append(apiState.Entities, apiEntity)
	}

	return apiState
}

// EntityToAPIEntity converts an internal entity to its API representation
func EntityToAPIEntity(entity *game.Entity, grid *game.Grid) EntityState {
	// Get the entity's position
	gridPos := grid.GetEntityPosition(entity)
	pos := [2]int{gridPos.X, gridPos.Y}

	// Convert action cards
	actionCards := make([]ActionCardRef, 0, len(entity.ActionCards))
	for _, card := range entity.ActionCards {
		// Convert ActionCardType to string for the action cost
		actionCost := 0
		switch card.Type {
		case game.OneActionCard:
			actionCost = 1
		case game.TwoActionCard:
			actionCost = 2
		case game.ThreeActionCard:
			actionCost = 3
		case game.FreeActionCard:
			actionCost = 0
		}

		actionCards = append(actionCards, ActionCardRef{
			ID:          card.ID,
			Name:        card.Name,
			Description: card.Description,
			ActionCost:  actionCost,
			Type:        string(card.Type),
		})
	}

	// Convert faction to string representation for the frontend
	factionStr := "neutral"
	switch entity.Faction {
	case game.GoodGuys:
		factionStr = "goodGuys"
	case game.BadGuys:
		factionStr = "badGuys"
	}

	return EntityState{
		ID:                 entity.Id,
		Name:               entity.Name,
		HP:                 entity.HP,
		MaxHP:              entity.HP, // For now, use current HP as max HP if not tracked separately
		AC:                 entity.AC,
		ActionsRemaining:   entity.ActionsRemaining,
		ReactionsRemaining: entity.ReactionsRemaining,
		Faction:            factionStr,
		ActionCards:        actionCards,
		Position:           pos,
	}
}