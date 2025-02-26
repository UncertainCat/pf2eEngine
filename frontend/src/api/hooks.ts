import { useState, useEffect, useCallback } from 'react';
import { apiClient } from './index';
import { GameState, GameEvent, EntityState, EventType } from './types';

/**
 * Custom hook for accessing the game state
 * Automatically updates when game state events are received
 */
export function useGameState(): {
  gameState: GameState | null;
  isLoading: boolean;
  error: Error | null;
  refresh: () => Promise<void>;
} {
  const [gameState, setGameState] = useState<GameState | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [error, setError] = useState<Error | null>(null);

  // Function to fetch the game state
  const fetchGameState = useCallback(async () => {
    try {
      setIsLoading(true);
      setError(null);
      const state = await apiClient.getGameState();
      setGameState(state);
    } catch (err) {
      setError(err instanceof Error ? err : new Error(String(err)));
    } finally {
      setIsLoading(false);
    }
  }, []);

  // Initial fetch and WebSocket connection
  useEffect(() => {
    fetchGameState();

    // Listen for game state updates
    const handleGameStateUpdate = (event: GameEvent<GameState>) => {
      if (event.data) {
        setGameState(event.data);
      }
    };

    apiClient.on<GameState>('gameState', handleGameStateUpdate);

    return () => {
      apiClient.off('gameState', handleGameStateUpdate);
    };
  }, [fetchGameState]);

  return { gameState, isLoading, error, refresh: fetchGameState };
}

/**
 * Custom hook for accessing a specific entity
 * Returns null if the entity is not found or the game state is not loaded
 */
export function useEntity(entityId: string | undefined): {
  entity: EntityState | null;
  isCurrentTurn: boolean;
} {
  const { gameState } = useGameState();
  const [entity, setEntity] = useState<EntityState | null>(null);
  const [isCurrentTurn, setIsCurrentTurn] = useState<boolean>(false);

  useEffect(() => {
    if (!gameState || !entityId) {
      setEntity(null);
      setIsCurrentTurn(false);
      return;
    }

    const foundEntity = gameState.entities.find(e => e.id === entityId) || null;
    setEntity(foundEntity);
    
    // Check if it's this entity's turn
    setIsCurrentTurn(gameState.currentTurn === entityId);
  }, [gameState, entityId]);

  return { entity, isCurrentTurn };
}

/**
 * Custom hook for tracking and managing steps/events
 */
export function useGameEvents(): {
  events: GameEvent[];
  pendingEvents: GameEvent[];
  nextEvent: () => void;
  stepMode: boolean;
  toggleStepMode: () => void;
} {
  const [events, setEvents] = useState<GameEvent[]>([]);
  const [pendingEvents, setPendingEvents] = useState<GameEvent[]>([]);
  const [stepMode, setStepMode] = useState<boolean>(true);

  // Process the next pending event
  const nextEvent = useCallback(() => {
    if (pendingEvents.length === 0) return;
    
    const event = pendingEvents[0];
    setEvents(prev => [...prev, event]);
    setPendingEvents(prev => prev.slice(1));
  }, [pendingEvents]);

  // Toggle between step mode and auto mode
  const toggleStepMode = useCallback(() => {
    setStepMode(prev => {
      const newMode = !prev;
      
      // If switching to auto mode, process all pending events
      if (!newMode && pendingEvents.length > 0) {
        setEvents(prev => [...prev, ...pendingEvents]);
        setPendingEvents([]);
      }
      
      return newMode;
    });
  }, [pendingEvents]);

  // Listen for game events via WebSocket
  useEffect(() => {
    const handleGameEvent = (event: GameEvent) => {
      // Ignore game state update events as they're handled separately
      if (event.type === EventType.GAME_STATE) {
        return;
      }
      
      if (stepMode) {
        // In step mode, add to pending events queue
        setPendingEvents(prev => [...prev, event]);
      } else {
        // In auto mode, add directly to processed events
        setEvents(prev => [...prev, event]);
      }
    };

    apiClient.on('all', handleGameEvent);

    return () => {
      apiClient.off('all', handleGameEvent);
    };
  }, [stepMode]);

  return { events, pendingEvents, nextEvent, stepMode, toggleStepMode };
}

/**
 * Custom hook for actions (sending commands to the server)
 */
export function useActions() {
  const [isSubmitting, setIsSubmitting] = useState<boolean>(false);
  const [error, setError] = useState<Error | null>(null);

  const sendAction = useCallback(async (
    entityId: string, 
    actionCardId: string, 
    params: Record<string, any> = {}
  ) => {
    try {
      setIsSubmitting(true);
      setError(null);
      
      await apiClient.sendCommand({
        entity_id: entityId,
        action_card_id: actionCardId,
        params
      });
      
      return true;
    } catch (err) {
      setError(err instanceof Error ? err : new Error(String(err)));
      return false;
    } finally {
      setIsSubmitting(false);
    }
  }, []);

  return { sendAction, isSubmitting, error };
}