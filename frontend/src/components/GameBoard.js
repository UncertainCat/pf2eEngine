import React from 'react';
import { X } from 'lucide-react';
import { useGame } from '../context/GameContext';

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

    const width = gameState.gridWidth || 10;
    const height = gameState.gridHeight || 10;
    const grid = [];

    // Create the grid cells
    for (let y = 0; y < height; y++) {
        const row = [];
        for (let x = 0; x < width; x++) {
            // Find entity at this position
            const entity = entities.find(e => e.position && e.position[0] === x && e.position[1] === y);
            
            // Get current HP for this entity if one exists
            const currentHP = entity ? getEntityCurrentHP(entity.id) : 0;
            
            row.push(
                <div
                    key={`${x}-${y}`}
                    className={`w-12 h-12 border border-gray-300 relative ${
                        (x + y) % 2 === 0 ? 'bg-gray-100' : 'bg-gray-200'
                    }`}
                >
                    {entity && (
                        <div
                            className={`absolute inset-0 flex items-center justify-center cursor-pointer ${
                                entity.id === selectedEntityId ? 'ring-2 ring-blue-500' : ''
                            } ${
                                entity.id === targetEntityId ? 'ring-2 ring-red-500' : ''
                            } ${
                                entity.id === currentEntityId ? 'ring-2 ring-yellow-500' : ''
                            }`}
                            onClick={() => {
                                if (selectedEntityId && entity.id !== selectedEntityId) {
                                    setTargetEntityId(entity.id);
                                } else {
                                    setSelectedEntityId(entity.id);
                                    setTargetEntityId(null);
                                }
                            }}
                        >
                            <div
                                className={`w-10 h-10 rounded-full flex items-center justify-center text-white font-bold ${
                                    entity.faction === 'goodGuys' ? 'bg-blue-500' : 'bg-red-500'
                                }`}
                            >
                                {entity.name.charAt(0)}
                                {currentHP <= 0 && (
                                    <div className="absolute inset-0 flex items-center justify-center">
                                        <X size={24} className="text-black" />
                                    </div>
                                )}
                            </div>

                            {/* HP indicator */}
                            <div className="absolute bottom-0 left-0 right-0 h-1 bg-gray-300">
                                <div
                                    className="h-1 bg-green-500"
                                    style={{ width: `${Math.max(0, (currentHP || 0) / (entity.maxHp || 1) * 100)}%` }}
                                />
                            </div>
                        </div>
                    )}
                </div>
            );
        }
        grid.push(
            <div key={y} className="flex">
                {row}
            </div>
        );
    }

    return (
        <div className="inline-block border border-gray-400">
            {grid}
        </div>
    );
};

export default GameBoard;