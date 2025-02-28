import React from 'react';
import { useGame } from '../context/GameContext';
import GameGrid from './GameGrid';

/**
 * GameBoard is the main component that displays the game grid
 * It handles entity selection and targeting
 */
const GameBoard = () => {
    const {
        gameState,
        entities,
        selectedEntityId,
        setSelectedEntityId,
        targetEntityId,
        setTargetEntityId,
        currentEntityId,
        getEntityCurrentHP
    } = useGame();

    // If no game state, show loading
    if (!gameState) {
        return (
            <div className="flex items-center justify-center h-64 bg-gray-100 rounded border border-gray-300">
                <div className="text-lg text-gray-500">Loading game state...</div>
            </div>
        );
    }
    
    // Grid dimensions
    const width = gameState.gridWidth || 10;
    const height = gameState.gridHeight || 10;
    
    // Handle entity click based on selection state
    const handleEntityClick = (entity) => {
        if (selectedEntityId && entity.id !== selectedEntityId) {
            // If we already have a selection and clicked a different entity, target it
            setTargetEntityId(entity.id);
        } else {
            // Otherwise select this entity and clear any target
            setSelectedEntityId(entity.id);
            setTargetEntityId(null);
        }
    };

    return (
        <GameGrid
            width={width}
            height={height}
            entities={entities}
            getEntityCurrentHP={getEntityCurrentHP}
            selectedEntityId={selectedEntityId}
            targetEntityId={targetEntityId}
            currentEntityId={currentEntityId}
            onEntityClick={handleEntityClick}
        />
    );
};

export default GameBoard;