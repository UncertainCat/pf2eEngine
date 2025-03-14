import React, { createContext, useState, useEffect, useContext, useCallback, useReducer } from 'react';
import { gameStateService, processGameEvent } from '../services/gameStateService';
import { websocketService, WsEvent } from '../services/websocketService';

// Create the context
const GameContext = createContext(null);

// Custom hook for using the game context
export function useGame() {
    return useContext(GameContext);
}

// Game state reducer to handle step-by-step updates
function gameStateReducer(state, action) {
    switch (action.type) {
        case 'SET_INITIAL_STATE':
            return action.payload;
            
        case 'PROCESS_EVENT':
            return processGameEvent(state, action.payload);
            
        case 'RESET_STATE':
            // This is now just a signal to perform a reset
            // The actual reset will happen asynchronously
            return state;
            
        case 'SET_RESET_STATE':
            // New action for updating state after reset
            return action.payload;

        case 'RECONCILE_STATE':
            // For handling complete state reconciliation
            return action.payload;
            
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
                const initialState = await gameStateService.getInitialState();
                console.log('Initial game state loaded:', initialState);
                
                // Initialize the game state
                dispatch({ 
                    type: 'SET_INITIAL_STATE', 
                    payload: initialState 
                });
                
                // Now that we have state, load history
                loadStepHistory();
            } catch (error) {
                console.error('Failed to load game state:', error);
            }
        };
        
        // Load step history
        const loadStepHistory = async () => {
            try {
                // First load up to 100 steps of history
                const history = await gameStateService.getStepHistory(0, 100);
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
        
        // Handler for game events
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
        websocketService.on(WsEvent.CONNECT, handleConnect);
        websocketService.on(WsEvent.DISCONNECT, handleDisconnect);
        websocketService.on(WsEvent.GAME_EVENT, handleGameEvent);
        
        // Connect to server
        websocketService.connect();
        
        // Cleanup when unmounting
        return () => {
            websocketService.off(WsEvent.CONNECT, handleConnect);
            websocketService.off(WsEvent.DISCONNECT, handleDisconnect);
            websocketService.off(WsEvent.GAME_EVENT, handleGameEvent);
            websocketService.disconnect();
        };
    }, [historyLoaded, gameState, processedEvents]);
    
    // Auto-advance through history or on user click
    const nextStep = useCallback(() => {
        if (pendingEvents.length === 0 || !gameState) return;
        
        // Get the next event
        const event = pendingEvents[0];
        
        // Process the event and update game state
        dispatch({
            type: 'PROCESS_EVENT',
            payload: event
        });
        
        // Move the event from pending to processed
        setProcessedEvents(prev => [...prev, event]);
        setPendingEvents(prev => prev.slice(1));
    }, [pendingEvents, gameState]);
    
    // Reset game state to initial state
    const resetGameState = useCallback(async () => {
        if (!gameState) return;
        
        // Signal that a reset is happening
        dispatch({ type: 'RESET_STATE' });
        
        try {
            // Get fresh initial state from server
            const initialState = await gameStateService.getInitialState();
            console.log('Loaded fresh initial state for reset:', initialState);
            
            // Set the new initial state
            dispatch({ 
                type: 'SET_RESET_STATE', 
                payload: initialState 
            });
            
            // Clear processed events
            setProcessedEvents([]);
            
            // Reload step history to repopulate pending events
            const history = await gameStateService.getStepHistory(0, 100);
            if (Array.isArray(history) && history.length > 0) {
                setPendingEvents(history);
                setHistoryIndex(history.length);
            }
        } catch (error) {
            console.error('Failed to reset game state:', error);
        }
    }, [gameState]);
    
    // Load more history if needed
    const loadMoreHistory = useCallback(async () => {
        if (!gameState) return;
        
        try {
            const moreHistory = await gameStateService.getStepHistory(historyIndex, 100);
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
    
    // Reconcile with server state if needed
    const reconcileState = useCallback(async () => {
        if (!gameState) return;
        
        try {
            const serverState = await gameStateService.getCurrentState();
            console.log('Reconciling with server state:', serverState);
            
            // Update game state with server state
            dispatch({
                type: 'RECONCILE_STATE',
                payload: serverState.data
            });
        } catch (error) {
            console.error('Failed to reconcile with server state:', error);
        }
    }, [gameState]);
    
    // Send an action to the server
    const sendAction = useCallback(async (entityId, actionCardId, targetId, params = {}) => {
        if (!entityId || !actionCardId) return;
        
        try {
            // Send to server
            const response = await gameStateService.sendAction(entityId, actionCardId, targetId, params);
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
    
    // Get entity current HP
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
        reconcileState,
        
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