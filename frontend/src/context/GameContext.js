import React, { createContext, useState, useEffect, useContext, useCallback, useReducer } from 'react';
import { apiClient, GameEvent, EntityState } from '../api';

// Create the context
const GameContext = createContext(null);

// Custom hook for using the game context
export function useGame() {
    return useContext(GameContext);
}

// Helper function to get default max HP based on entity type
// This ensures we have reasonable values when the server data is incomplete
function getDefaultMaxHP(entity) {
    // Use existing maxHp if it seems valid (greater than current HP)
    if (entity.maxHp && entity.maxHp >= entity.hp) {
        return entity.maxHp;
    }
    
    // Otherwise, make some guesses based on name/faction
    const name = entity.name.toLowerCase();
    
    if (name.includes("warrior")) {
        return 30; // Warriors have 30 HP in the server code
    } else if (name.includes("goblin")) {
        return 20; // Goblins have 20 HP in the server code
    } else if (entity.faction === "goodGuys") {
        return Math.max(30, entity.hp); // Good guys are usually tougher
    } else {
        return Math.max(20, entity.hp); // Default for others
    }
}

// Game state reducer to handle step-by-step updates
function gameStateReducer(state, action) {
    switch (action.type) {
        case 'SET_INITIAL_STATE':
            // Process entities to ensure proper HP values
            const entities = action.payload.entities.map(e => ({
                ...e,
                // Fix HP if it's incorrect (ensure we start with full health)
                hp: e.hp <= 0 ? getDefaultMaxHP(e) : e.hp,
                // Fix maxHp if it's not set correctly
                maxHp: e.maxHp <= 0 || e.maxHp < e.hp ? getDefaultMaxHP(e) : e.maxHp
            }));

            return {
                ...action.payload,
                entities,
                // Keep track of original positions and HPs for reset
                originalEntities: entities.map(e => ({...e}))
            };
            
        case 'RESET_TO_INITIAL':
            return {
                ...state,
                entities: state.originalEntities.map(e => ({...e})),
                currentTurn: state.originalCurrentTurn
            };
            
        case 'UPDATE_ENTITY_HP':
            return {
                ...state,
                entities: state.entities.map(entity => 
                    entity.id === action.payload.entityId
                        ? { ...entity, hp: action.payload.hp }
                        : entity
                )
            };
            
        case 'UPDATE_ENTITY_POSITION':
            return {
                ...state,
                entities: state.entities.map(entity => 
                    entity.id === action.payload.entityId
                        ? { ...entity, position: action.payload.position }
                        : entity
                )
            };
            
        case 'UPDATE_CURRENT_TURN':
            return {
                ...state,
                currentTurn: action.payload.entityId
            };
            
        case 'UPDATE_ENTITY_ACTIONS':
            return {
                ...state,
                entities: state.entities.map(entity => 
                    entity.id === action.payload.entityId
                        ? { 
                            ...entity, 
                            actionsRemaining: action.payload.actionsRemaining,
                            reactionsRemaining: action.payload.reactionsRemaining
                        }
                        : entity
                )
            };
            
        default:
            return state;
    }
}

export const GameProvider = ({ children }) => {
    // Connection state
    const [connected, setConnected] = useState(false);
    
    // Selection state (UI state)
    const [selectedEntityId, setSelectedEntityId] = useState(null);
    const [selectedActionId, setSelectedActionId] = useState(null);
    const [targetEntityId, setTargetEntityId] = useState(null);
    
    // History/event tracking
    const [processedEvents, setProcessedEvents] = useState([]);
    const [pendingEvents, setPendingEvents] = useState([]);
    const [historyLoaded, setHistoryLoaded] = useState(false);
    const [historyIndex, setHistoryIndex] = useState(0);
    
    // Main game state
    const [gameState, dispatch] = useReducer(gameStateReducer, null);
    
    // Load initial game state 
    useEffect(() => {
        const loadGameState = async () => {
            try {
                const response = await apiClient.getGameState();
                console.log('Initial game state response:', response);
                
                // The response is wrapped in a GameEvent object with the actual state in data
                if (response && response.data) {
                    const state = response.data;
                    console.log('Initial game state:', state);
                    
                    // Store the original current turn for resets
                    state.originalCurrentTurn = state.currentTurn;
                    
                    // Initialize the game state
                    dispatch({ 
                        type: 'SET_INITIAL_STATE', 
                        payload: state 
                    });
                    
                    // Now that we have state, load history
                    loadStepHistory();
                } else {
                    console.error('Invalid game state response format:', response);
                }
            } catch (error) {
                console.error('Failed to load game state:', error);
            }
        };
        
        // Load step history
        const loadStepHistory = async () => {
            try {
                // First load up to 100 steps of history
                const history = await apiClient.getStepHistory(0, 100);
                console.log('Loaded step history:', history);
                
                if (Array.isArray(history) && history.length > 0) {
                    // Add all history events to the pending queue for stepping through
                    setPendingEvents(history);
                    setHistoryLoaded(true);
                    setHistoryIndex(history.length);
                } else {
                    console.log('No step history found or history is empty');
                    setHistoryLoaded(true);
                }
            } catch (error) {
                console.error('Failed to load step history:', error);
                setHistoryLoaded(true); // Still mark as loaded even if there's an error
            }
        };
        
        loadGameState();
    }, []);
    
    // Connect to WebSocket and handle events
    useEffect(() => {
        if (!historyLoaded || !gameState) return; // Don't connect to WebSocket until history is loaded and game state is available
        
        // Handler for connection status
        const handleConnect = () => {
            console.log('Connected to server');
            setConnected(true);
        };
        
        const handleDisconnect = () => {
            console.log('Disconnected from server');
            setConnected(false);
        };
        
        // Handler for game state updates
        const handleGameState = (event) => {
            if (event.data) {
                console.log('Game state update from server:', event.data);
                // We don't automatically update the game state here, as we want
                // all updates to happen through the step queue
                
                // Instead, we compare if there are new entities or major changes
                // and then reset the game state only if needed
                if (gameState && event.data.entities && 
                    (event.data.entities.length !== gameState.entities.length)) {
                    console.log('Major game state change detected, resetting state');
                    dispatch({ 
                        type: 'SET_INITIAL_STATE', 
                        payload: event.data 
                    });
                }
            }
        };
        
        // Handler for all other game events
        const handleGameEvent = (event) => {
            console.log('Game event received via WebSocket:', event);
            // Don't duplicate game state events in the combat log
            if (event.type === 'GAME_STATE') return;
            
            // Only add events that are newer than our history
            // This prevents duplication of events already loaded from history
            const eventTimestamp = new Date(event.timestamp).getTime();
            const shouldAddEvent = processedEvents.length === 0 || 
                !processedEvents.some(e => {
                    const existingTimestamp = new Date(e.timestamp).getTime();
                    return existingTimestamp === eventTimestamp && e.type === event.type;
                });
                
            if (shouldAddEvent) {
                console.log('Adding new event to queue:', event);
                // Add to pending events for manual stepping
                setPendingEvents(prev => [...prev, event]);
            } else {
                console.log('Skipping duplicate event:', event);
            }
        };
        
        // Register event handlers
        apiClient.on('connect', handleConnect);
        apiClient.on('disconnect', handleDisconnect);
        apiClient.on('gamestate', handleGameState); // Note: lowercase to match the event type
        apiClient.on('all', handleGameEvent);
        
        // Connect to server
        apiClient.connect();
        
        // Cleanup when unmounting
        return () => {
            apiClient.off('connect', handleConnect);
            apiClient.off('disconnect', handleDisconnect);
            apiClient.off('gamestate', handleGameState);
            apiClient.off('all', handleGameEvent);
            apiClient.disconnect();
        };
    }, [historyLoaded, gameState, processedEvents]);
    
    // Process an event and update game state
    const processEvent = useCallback((event) => {
        if (!gameState) return;
        
        console.log('Processing event:', event);
        
        switch (event.type) {
            case 'DAMAGE':
            case 'DAMAGE_RESULT':
                // Process damage events by updating entity HP
                if (event.data && event.data.target && event.data.target.id) {
                    if (event.data.taken) {
                        const targetId = event.data.target.id;
                        const entity = gameState.entities.find(e => e.id === targetId);
                        if (entity) {
                            const newHp = Math.max(0, entity.hp - event.data.taken);
                            dispatch({
                                type: 'UPDATE_ENTITY_HP',
                                payload: {
                                    entityId: targetId,
                                    hp: newHp
                                }
                            });
                        }
                    }
                }
                break;
                
            case 'TURN_START':
                // Update current turn
                if (event.data && event.data.entity && event.data.entity.id) {
                    dispatch({
                        type: 'UPDATE_CURRENT_TURN',
                        payload: {
                            entityId: event.data.entity.id
                        }
                    });
                    
                    // Reset actions for the entity
                    const entity = gameState.entities.find(e => e.id === event.data.entity.id);
                    if (entity) {
                        dispatch({
                            type: 'UPDATE_ENTITY_ACTIONS',
                            payload: {
                                entityId: event.data.entity.id,
                                actionsRemaining: 3, // Default for PF2e
                                reactionsRemaining: 1
                            }
                        });
                    }
                }
                break;
                
            case 'ENTITY_MOVE':
                // Update entity position
                if (event.data && event.data.entity && event.data.entity.id && event.data.position) {
                    dispatch({
                        type: 'UPDATE_ENTITY_POSITION',
                        payload: {
                            entityId: event.data.entity.id,
                            position: event.data.position
                        }
                    });
                }
                break;
                
            case 'ACTION_COMPLETE':
                // Update action count
                if (event.data && event.data.entity && event.data.entity.id && 
                    typeof event.data.actionsRemaining !== 'undefined') {
                    dispatch({
                        type: 'UPDATE_ENTITY_ACTIONS',
                        payload: {
                            entityId: event.data.entity.id,
                            actionsRemaining: event.data.actionsRemaining,
                            reactionsRemaining: event.data.reactionsRemaining || 0
                        }
                    });
                }
                break;
                
            default:
                // For other event types, log but don't update state
                console.log('Unhandled event type:', event.type);
                break;
        }
    }, [gameState]);
    
    // Auto-advance through history or on user click
    const nextStep = useCallback(() => {
        if (pendingEvents.length === 0 || !gameState) return;
        
        // Get the next event
        const event = pendingEvents[0];
        
        // Process the event and update game state
        processEvent(event);
        
        // Move the event from pending to processed
        setProcessedEvents(prev => [...prev, event]);
        setPendingEvents(prev => prev.slice(1));
    }, [pendingEvents, gameState, processEvent]);
    
    // Reset game state to initial state
    const resetGameState = useCallback(() => {
        if (!gameState) return;
        
        // Reset to initial state
        dispatch({ type: 'RESET_TO_INITIAL' });
        
        // Clear processed events and reload pending events
        setProcessedEvents([]);
        
        // Reload step history to repopulate pending events
        const loadStepHistory = async () => {
            try {
                const history = await apiClient.getStepHistory(0, 100);
                if (Array.isArray(history) && history.length > 0) {
                    setPendingEvents(history);
                    setHistoryIndex(history.length);
                }
            } catch (error) {
                console.error('Failed to reload step history:', error);
            }
        };
        
        loadStepHistory();
    }, [gameState]);
    
    // Load more history if needed
    const loadMoreHistory = useCallback(async () => {
        if (!gameState) return;
        
        try {
            const moreHistory = await apiClient.getStepHistory(historyIndex, 100);
            if (Array.isArray(moreHistory) && moreHistory.length > 0) {
                console.log(`Loaded ${moreHistory.length} more historical events`);
                setPendingEvents(prev => [...moreHistory, ...prev]);
                setHistoryIndex(historyIndex + moreHistory.length);
            } else {
                console.log('No more history to load');
            }
        } catch (error) {
            console.error('Failed to load more history:', error);
        }
    }, [historyIndex, gameState]);
    
    // Send an action to the server
    const sendAction = useCallback(async (entityId, actionCardId, targetId, params = {}) => {
        if (!entityId || !actionCardId || !targetId) return;
        
        try {
            // Build the command
            const command = {
                entity_id: entityId,
                action_card_id: actionCardId,
                params: {
                    ...params,
                    targetID: targetId
                }
            };
            
            // Send to server
            const response = await apiClient.sendCommand(command);
            console.log('Action response:', response);
            
            // Clear selections
            setSelectedActionId(null);
            setTargetEntityId(null);
            
            return true;
        } catch (error) {
            console.error('Error sending action:', error);
            return false;
        }
    }, []);
    
    // Entity helper functions
    const getEntityById = useCallback((id) => {
        if (!gameState || !id) return null;
        return gameState.entities.find(e => e.id === id) || null;
    }, [gameState]);
    
    // Get current HP directly from game state
    const getEntityCurrentHP = useCallback((id) => {
        if (!gameState || !id) return 0;
        const entity = gameState.entities.find(e => e.id === id);
        return entity ? entity.hp : 0;
    }, [gameState]);
    
    // Context value
    const value = {
        // Connection state
        connected,
        
        // Game state
        gameState,
        entities: gameState?.entities || [],
        currentEntityId: gameState?.currentTurn,
        
        // Selection state
        selectedEntityId,
        setSelectedEntityId,
        selectedActionId, 
        setSelectedActionId,
        targetEntityId,
        setTargetEntityId,
        
        // Entity helpers
        getEntityById,
        getEntityCurrentHP,
        
        // Combat log and steps
        processedEvents,
        pendingEvents,
        hasNextEvent: pendingEvents.length > 0,
        nextStep,
        loadMoreHistory,
        resetGameState,
        
        // Status
        historyLoaded,
        
        // Actions
        sendAction
    };
    
    return (
        <GameContext.Provider value={value}>
            {children}
        </GameContext.Provider>
    );
};