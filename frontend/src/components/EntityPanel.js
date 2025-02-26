import React from 'react';
import { Heart, Shield, Send } from 'lucide-react';
import { useGame } from '../context/GameContext';

const EntityPanel = () => {
    const {
        selectedEntity,
        currentEntity,
        selectedAction,
        setSelectedAction,
        selectedTarget,
        getEntityActions,
        sendAction
    } = useGame();

    if (!selectedEntity) {
        return (
            <div className="p-4 border-b">
                <p className="text-gray-500 italic">Select an entity to view details</p>
            </div>
        );
    }

    const actions = getEntityActions(selectedEntity);
    const isCurrentTurn = selectedEntity.id === currentEntity?.id;

    return (
        <div className="p-4 border-b">
            <h2 className="text-lg font-bold flex items-center">
                <span className={`mr-2 w-3 h-3 rounded-full ${selectedEntity.faction === 0 ? 'bg-blue-500' : 'bg-red-500'}`}></span>
                {selectedEntity.name}
                {isCurrentTurn && (
                    <span className="ml-2 text-xs bg-yellow-300 text-yellow-800 px-2 py-1 rounded">Current Turn</span>
                )}
            </h2>

            <div className="mt-2 space-y-2">
                <div className="flex items-center">
                    <Heart size={16} className="text-red-500 mr-2" />
                    <div className="w-full bg-gray-200 rounded-full h-4">
                        <div
                            className="bg-red-500 h-4 rounded-full"
                            style={{ width: `${((selectedEntity.hp || 0) / (selectedEntity.maxHp || selectedEntity.hp || 1)) * 100}%` }}
                        >
                        </div>
                    </div>
                    <span className="ml-2">{selectedEntity.hp || 0}/{selectedEntity.maxHp || selectedEntity.hp || 30}</span>
                </div>
                <div className="flex items-center">
                    <Shield size={16} className="text-blue-500 mr-2" />
                    <span>AC: {selectedEntity.ac}</span>
                </div>
                {isCurrentTurn && (
                    <div className="flex items-center mt-1">
            <span className="text-sm text-gray-700">
              Actions: {selectedEntity.actionsRemaining || 3} remaining
            </span>
                    </div>
                )}
            </div>

            {isCurrentTurn && (
                <div className="mt-4">
                    <h3 className="font-bold">Actions:</h3>
                    <div className="mt-2 space-y-2">
                        {actions.map(action => (
                            <button
                                key={action.id}
                                className={`w-full text-left p-2 rounded ${selectedAction?.id === action.id ? 'bg-blue-500 text-white' : 'bg-gray-100 hover:bg-gray-200'}`}
                                onClick={() => setSelectedAction(action)}
                            >
                                <div className="font-bold">{action.name}</div>
                                <div className="text-sm">{action.description}</div>
                                <div className="text-xs mt-1">{action.actionCost} action{action.actionCost !== 1 ? 's' : ''}</div>
                            </button>
                        ))}
                    </div>
                </div>
            )}

            {selectedAction && selectedTarget && (
                <button
                    className="mt-4 w-full bg-green-500 text-white p-2 rounded flex items-center justify-center"
                    onClick={sendAction}
                >
                    <Send size={16} className="mr-2" />
                    Execute Action
                </button>
            )}
        </div>
    );
};

export default EntityPanel;