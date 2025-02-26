import React from 'react';
import { X } from 'lucide-react';
import { useGame } from '../context/GameContext';

const GameBoard = () => {
    const {
        entities,
        selectedEntity,
        setSelectedEntity,
        selectedTarget,
        setSelectedTarget,
        currentEntity,
        selectedAction
    } = useGame();

    const width = 10;
    const height = 10;
    const grid = [];

    // Create the grid cells
    for (let y = 0; y < height; y++) {
        const row = [];
        for (let x = 0; x < width; x++) {
            // Find entity at this position
            const entity = entities.find(e => e.x === x && e.y === y);

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
                                entity === selectedEntity ? 'ring-2 ring-blue-500' : ''
                            } ${
                                entity === selectedTarget ? 'ring-2 ring-red-500' : ''
                            } ${
                                entity.id === currentEntity?.id ? 'ring-2 ring-yellow-500' : ''
                            }`}
                            onClick={() => {
                                if (selectedAction && selectedEntity && entity.id !== selectedEntity.id) {
                                    setSelectedTarget(entity);
                                } else {
                                    setSelectedEntity(entity);
                                    setSelectedTarget(null);
                                }
                            }}
                        >
                            <div
                                className={`w-10 h-10 rounded-full flex items-center justify-center text-white font-bold ${
                                    entity.faction === 0 ? 'bg-blue-500' : 'bg-red-500'
                                }`}
                            >
                                {entity.name.charAt(0)}
                                {entity.hp <= 0 && (
                                    <div className="absolute inset-0 flex items-center justify-center">
                                        <X size={24} className="text-black" />
                                    </div>
                                )}
                            </div>

                            {/* HP indicator */}
                            <div className="absolute bottom-0 left-0 right-0 h-1 bg-gray-300">
                                <div
                                    className="h-1 bg-green-500"
                                    style={{ width: `${Math.max(0, (entity.hp || 0) / (entity.maxHp || entity.hp || 1) * 100)}%` }}
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