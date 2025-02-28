import { apiClient } from '../api';

// Helper function to get default max HP based on entity type
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

// Process an incoming state from the server to ensure it has complete and valid data
export function processStateData(stateData) {
    if (!stateData || !stateData.entities) {
        return stateData;
    }

    // Process entities to ensure proper HP values
    const entities = stateData.entities.map(e => ({
        ...e,
        // Fix HP if it's incorrect (ensure we start with full health)
        hp: e.hp <= 0 ? getDefaultMaxHP(e) : e.hp,
        // Fix maxHp if it's not set correctly
        maxHp: e.maxHp <= 0 || e.maxHp < e.hp ? getDefaultMaxHP(e) : e.maxHp
    }));

    return {
        ...stateData,
        entities,
        // Keep track of original positions and HPs for reset
        originalEntities: entities.map(e => ({...e}))
    };
}

// Process a single game event and update the current state accordingly
export function processGameEvent(state, event) {
    if (!state || !event) return state;

    switch (event.type) {
        case 'DAMAGE':
        case 'DAMAGE_RESULT':
            // Process damage events by updating entity HP
            if (event.data && event.data.target && event.data.target.id && event.data.taken) {
                const targetId = event.data.target.id;
                const entity = state.entities.find(e => e.id === targetId);
                if (entity) {
                    const newHp = Math.max(0, entity.hp - event.data.taken);
                    return {
                        ...state,
                        entities: state.entities.map(e => 
                            e.id === targetId ? { ...e, hp: newHp } : e
                        )
                    };
                }
            }
            break;
            
        case 'TURN_START':
            // Update current turn and reset actions
            if (event.data && event.data.entity && event.data.entity.id) {
                const entityId = event.data.entity.id;
                return {
                    ...state,
                    currentTurn: entityId,
                    entities: state.entities.map(e => 
                        e.id === entityId ? {
                            ...e,
                            actionsRemaining: 3, // Default for PF2e
                            reactionsRemaining: 1
                        } : e
                    )
                };
            }
            break;
            
        case 'ENTITY_MOVE':
            // Update entity position with animation
            if (event.data && event.data.entity && event.data.entity.id && event.data.position) {
                const entityId = event.data.entity.id;
                const entity = state.entities.find(e => e.id === entityId);
                
                if (entity) {
                    const prevPosition = entity.position || [0, 0];
                    const newPosition = event.data.position;
                    
                    return {
                        ...state,
                        entities: state.entities.map(e => 
                            e.id === entityId ? {
                                ...e,
                                position: newPosition,
                                prevPosition: prevPosition,  // Store previous position
                                isMoving: true,              // Flag that entity is moving
                                moveStartTime: Date.now(),   // Track when movement started
                            } : e
                        )
                    };
                }
            }
            break;
            
        case 'ACTION_COMPLETE':
            // Update action count
            if (event.data && event.data.entity && event.data.entity.id && 
                typeof event.data.actionsRemaining !== 'undefined') {
                return {
                    ...state,
                    entities: state.entities.map(e => 
                        e.id === event.data.entity.id ? {
                            ...e,
                            actionsRemaining: event.data.actionsRemaining,
                            reactionsRemaining: event.data.reactionsRemaining || 0
                        } : e
                    )
                };
            }
            break;
            
        case 'GAME_STATE':
            // Only use for reconciliation (complete state replacement) 
            // when specifically requested, not automatic updates
            if (event.metadata && event.metadata.isReconciliation) {
                return processStateData(event.data);
            }
            break;
    }
    
    return state;
}

// Reset game state to initial state - now deprecated
// This is now handled directly in the GameContext component
export function resetGameState(state) {
    console.warn("resetGameState in gameStateService is deprecated. Use the resetGameState function from GameContext instead.");
    
    // For backward compatibility, just return the state as-is
    return state;
}

// Game state API functions
export const gameStateService = {
    // Load initial game state
    async getInitialState() {
        try {
            const response = await apiClient.getInitialState();
            if (response && response.data) {
                // Process the state data to ensure it's valid
                const processedState = processStateData(response.data);
                processedState.originalCurrentTurn = processedState.currentTurn;
                return processedState;
            }
            throw new Error('Invalid state data format');
        } catch (error) {
            console.error('Failed to load initial game state:', error);
            throw error;
        }
    },

    // Load step history
    async getStepHistory(index = 0, limit = 100) {
        try {
            return await apiClient.getStepHistory(index, limit);
        } catch (error) {
            console.error('Failed to load step history:', error);
            throw error;
        }
    },

    // Get current server state (without any processing)
    async getCurrentState() {
        try {
            return await apiClient.getCurrentState();
        } catch (error) {
            console.error('Failed to get current state:', error);
            throw error;
        }
    },

    // Get current server state for reconciliation
    async reconcileState() {
        try {
            const response = await apiClient.getCurrentState();
            if (response && response.data) {
                // Mark this as a reconciliation update
                response.metadata = { ...response.metadata, isReconciliation: true };
                return response;
            }
            throw new Error('Invalid state data for reconciliation');
        } catch (error) {
            console.error('Failed to reconcile game state:', error);
            throw error;
        }
    },

    // Send an action to the server
    async sendAction(entityId, actionCardId, targetId, params = {}) {
        if (!entityId || !actionCardId) {
            throw new Error('Missing required parameters for action');
        }

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
            return await apiClient.sendCommand(command);
        } catch (error) {
            console.error('Error sending action:', error);
            throw error;
        }
    }
};