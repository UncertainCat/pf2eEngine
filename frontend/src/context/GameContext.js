import React, { createContext, useState, useEffect, useContext } from 'react';

// Create the context
const GameContext = createContext();

// Custom hook for using the game context
export function useGame() {
    return useContext(GameContext);
}

export const GameProvider = ({ children }) => {
    // WebSocket connection
    const [socket, setSocket] = useState(null);
    const [connected, setConnected] = useState(false);

    // Game state
    const [entities, setEntities] = useState([]);
    const [currentEntity, setCurrentEntity] = useState(null);
    const [selectedEntity, setSelectedEntity] = useState(null);
    const [selectedAction, setSelectedAction] = useState(null);
    const [selectedTarget, setSelectedTarget] = useState(null);
    const [logs, setLogs] = useState([]);

    // Step-by-step playback controls
    const [stepMode, setStepMode] = useState(true);
    const [steps, setSteps] = useState([]);
    const [currentStepIndex, setCurrentStepIndex] = useState(0);
    const [lastFetchedStepIndex, setLastFetchedStepIndex] = useState(0);
    const [pendingSteps, setPendingSteps] = useState([]);

    // Connect to WebSocket
    useEffect(() => {
        const ws = new WebSocket(`ws://${window.location.hostname}:8080/ws`);

        ws.onopen = () => {
            console.log('Connected to server');
            setConnected(true);
        };

        ws.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                console.log('Message from server:', data);

                if (data.type === 'gameUpdate') {
                    if (data.entities) {
                        // Process entities to ensure they have maxHp
                        const processedEntities = data.entities.map(entity => ({
                            ...entity,
                            maxHp: entity.maxHp || entity.hp || 30 // Use hp as fallback or default to 30
                        }));
                        setEntities(processedEntities);
                    }
                    if (data.currentTurn !== undefined) {
                        const current = data.entities?.find(e => e.id === data.currentTurn);
                        setCurrentEntity(current || null);
                    }
                }

                // Add to logs if it's a message
                if (data.message || data.type) {
                    // Extract metadata to ensure consistent format with HTTP steps
                    const meta = data.metadata || data.data || {};
                    const eventType = (data.type || '').toUpperCase();
                    let message = data.message;

                    // If no message provided, try to generate one
                    if (!message) {
                        // Handle common event types
                        if (eventType.includes('ATTACK')) {
                            if (meta.attacker && meta.defender) {
                                message = `${meta.attacker} attacks ${meta.defender}`;
                            }
                        }
                        else if (eventType.includes('DAMAGE')) {
                            if (meta.source && meta.target) {
                                const amount = meta.damage ? ` ${meta.damage}` : '';
                                message = `${meta.source} dealt${amount} damage to ${meta.target}`;
                            }
                        }
                        else if (eventType.includes('TURN')) {
                            if (meta.entity) {
                                message = eventType.includes('START') ?
                                    `${meta.entity}'s turn begins` :
                                    `${meta.entity}'s turn ends`;
                            }
                        }

                        // If we still don't have a message, create one from the event type
                        if (!message && eventType) {
                            // Convert from SNAKE_CASE to Title Case with spaces
                            const readableType = eventType
                                .split('_')
                                .map(word => word.charAt(0) + word.slice(1).toLowerCase())
                                .join(' ');
                            message = readableType;
                        }
                    }

                    const newLog = {
                        message: message || "Game event occurred",
                        type: data.type || "INFO",
                        timestamp: new Date(),
                        data: meta,
                        processed: false
                    };

                    if (stepMode) {
                        // In step mode, queue steps for manual advance
                        setPendingSteps(prev => [...prev, newLog]);
                    } else {
                        // In auto mode, add directly to logs
                        setLogs(prev => [...prev, newLog]);
                    }
                }
            } catch (error) {
                console.error('Error parsing WebSocket message:', error);
                // If it's just a plain text message
                const newLog = {
                    message: event.data,
                    type: "INFO",
                    timestamp: new Date(),
                    processed: false
                };

                if (stepMode) {
                    setPendingSteps(prev => [...prev, newLog]);
                } else {
                    setLogs(prev => [...prev, newLog]);
                }
            }
        };

        ws.onclose = () => {
            console.log('Disconnected from server');
            setConnected(false);
        };

        setSocket(ws);

        return () => {
            ws.close();
        };
    }, []);

    // Fetch initial game state and updates
    useEffect(() => {
        const fetchSteps = async () => {
            try {
                const response = await fetch(`/steps?index=${lastFetchedStepIndex}`);
                if (response.ok) {
                    const newSteps = await response.json();
                    if (newSteps && newSteps.length > 0) {
                        setLastFetchedStepIndex(lastFetchedStepIndex + newSteps.length);

                        // Process new steps
                        const formattedSteps = newSteps.map(step => {
                            // Create descriptive messages based on step data
                            let message = step.message;
                            const stepType = (step.type || '').toUpperCase();
                            const meta = step.metadata || {};
                            let logData = { ...meta };

                            // If we don't have a message, try to generate one based on available data
                            if (!message) {
                                // Handle common step types
                                if (stepType.includes('ATTACK')) {
                                    if (meta.Attacker && meta.Defender) {
                                        if (stepType.includes('BEFORE') || stepType === 'ATTACK') {
                                            message = `${meta.Attacker} attacks ${meta.Defender}`;
                                        } else if (stepType.includes('AFTER') || stepType.includes('RESULT')) {
                                            const result = meta.Degree ? ` (${meta.Degree})` : '';
                                            message = `${meta.Attacker}'s attack on ${meta.Defender} resolved${result}`;
                                        }
                                    }
                                }
                                else if (stepType.includes('DAMAGE')) {
                                    if (meta.Source && meta.Target) {
                                        if (stepType.includes('BEFORE')) {
                                            message = `${meta.Source} is about to deal damage to ${meta.Target}`;
                                        } else if (stepType.includes('AFTER') || stepType.includes('RESULT')) {
                                            const amount = meta.Taken ? ` ${meta.Taken}` : '';
                                            message = `${meta.Source} dealt${amount} damage to ${meta.Target}`;
                                        }
                                    }
                                }
                                else if (stepType.includes('TURN')) {
                                    const entity = meta.Entity || meta.entity;
                                    if (entity) {
                                        if (stepType.includes('START')) {
                                            message = `${entity}'s turn begins`;
                                        } else if (stepType.includes('END')) {
                                            message = `${entity}'s turn ends`;
                                        }
                                    }
                                }

                                // If we still don't have a message, try some more generic approaches
                                if (!message) {
                                    // Look for any entity references
                                    if (meta.Entity || meta.entity) {
                                        const entity = meta.Entity || meta.entity;
                                        message = `${entity} performs an action`;
                                    }
                                    else if (meta.Attacker || meta.Source) {
                                        const actor = meta.Attacker || meta.Source;
                                        message = `${actor} performs an action`;
                                    }
                                    else if (meta.Defender || meta.Target) {
                                        const target = meta.Defender || meta.Target;
                                        message = `${target} is affected`;
                                    }
                                    // Still no message? Create one from the step type
                                    else if (stepType) {
                                        // Convert from SNAKE_CASE to Title Case with spaces
                                        const readableType = stepType
                                            .split('_')
                                            .map(word => word.charAt(0) + word.slice(1).toLowerCase())
                                            .join(' ');
                                        message = readableType;
                                    }
                                }
                            }

                            return {
                                message: message || "Game event occurred",
                                type: step.type || "INFO",
                                timestamp: new Date(),
                                data: logData,
                                processed: false
                            };
                        });

                        if (stepMode) {
                            // In step mode, queue steps for manual advance
                            setPendingSteps(prev => [...prev, ...formattedSteps]);
                        } else {
                            // In auto mode, add directly to logs and process
                            setLogs(prev => [...prev, ...formattedSteps]);
                            processSteps(formattedSteps);
                        }
                    }
                }
            } catch (error) {
                console.error('Error fetching steps:', error);
            }
        };

        // Initial fetch and set up interval
        fetchSteps();
        const interval = setInterval(fetchSteps, 1000);

        return () => clearInterval(interval);
    }, [lastFetchedStepIndex, stepMode]);

    // Process a collection of steps
    const processSteps = (stepsToProcess) => {
        stepsToProcess.forEach(step => {
            if (step.processed) return;

            // Mark as processed
            step.processed = true;

            // Update game state based on step type
            if (step.type === 'DAMAGE_RESULT' && step.data) {
                // Update HP for the entity that took damage
                const { target, taken } = step.data;
                if (target && target.id) {
                    setEntities(prev => prev.map(e =>
                        e.id === target.id
                            ? { ...e, hp: Math.max(0, e.hp - (taken || 0)) }
                            : e
                    ));
                }
            }
            else if (step.type === 'TURN_START' && step.data && step.data.entity) {
                // Update current entity
                const entityId = step.data.entity.id;
                setCurrentEntity(prev => entities.find(e => e.id === entityId) || prev);
            }
        });
    };

    // Advance to the next step
    const nextStep = () => {
        if (pendingSteps.length === 0) return;

        // Get the next step
        const step = pendingSteps[0];

        // Process the step
        processSteps([step]);

        // Move the step from pending to logs
        setLogs(prev => [...prev, step]);
        setPendingSteps(prev => prev.slice(1));

        // Increment the step index
        setCurrentStepIndex(prev => prev + 1);
    };

    // Toggle step mode
    const toggleStepMode = () => {
        setStepMode(prev => !prev);

        if (stepMode) {
            // Switching from step mode to auto mode
            // Process all pending steps
            processSteps(pendingSteps);
            setLogs(prev => [...prev, ...pendingSteps]);
            setPendingSteps([]);
        }
    };

    // Send action to server
    const sendAction = () => {
        if (!selectedEntity || !selectedAction || !selectedTarget) return;

        const command = {
            entity_id: selectedEntity.id,
            action_card_id: selectedAction.id,
            params: {
                targetID: selectedTarget.id
            }
        };

        // Send via HTTP
        fetch(`/action`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(command)
        }).then(response => {
            if (!response.ok) {
                console.error('Error sending action');
            }
        }).catch(error => {
            console.error('Error:', error);
        });

        // Clear selections
        setSelectedAction(null);
        setSelectedTarget(null);
    };

    // Get available actions for the selected entity
    const getEntityActions = (entity) => {
        if (!entity || !entity.actionCards) {
            return [
                { id: '101', name: 'Strike', description: 'Make a melee attack', actionCost: 1 },
                { id: '102', name: 'Shield Block', description: 'Use your shield to block damage', actionCost: 0 }
            ];
        }
        return entity.actionCards;
    };

    const value = {
        // Connection state
        connected,

        // Game state
        entities,
        currentEntity,
        selectedEntity,
        setSelectedEntity,
        selectedAction,
        setSelectedAction,
        selectedTarget,
        setSelectedTarget,
        logs,

        // Actions
        sendAction,
        getEntityActions,

        // Step-by-step playback
        stepMode,
        toggleStepMode,
        nextStep,
        pendingSteps,
        currentStepIndex,
        hasNextStep: pendingSteps.length > 0
    };

    return (
        <GameContext.Provider value={value}>
            {children}
        </GameContext.Provider>
    );
};