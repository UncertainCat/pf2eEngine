import React from 'react';
import { Heart, Shield, Send } from 'lucide-react';
import { useGame } from '../context/GameContext';

const EntityPanel = () => {
    const {
        selectedEntityId,
        currentEntityId,
        getEntityById,
        selectedActionId,
        setSelectedActionId,
        targetEntityId,
        getEntityCurrentHP,
        sendAction
    } = useGame();

    // Get the selected entity
    const selectedEntity = getEntityById(selectedEntityId);
    const targetEntity = getEntityById(targetEntityId);
    
    if (!selectedEntity) {
        return (
            <div className="p-4 border-b">
                <p className="text-gray-500 italic">Select an entity to view details</p>
            </div>
        );
    }

    const actions = selectedEntity.actionCards || [];
    const isCurrentTurn = selectedEntity.id === currentEntityId;
    const currentHP = getEntityCurrentHP(selectedEntity.id);
    
    // Find the selected action
    const selectedAction = selectedEntity.actionCards?.find(card => card.id === selectedActionId);

    return (
        <div className="p-4 border-b">
            <h2 className="text-lg font-bold flex items-center">
                <span className={`mr-2 w-3 h-3 rounded-full ${selectedEntity.faction === 'goodGuys' ? 'bg-blue-500' : 'bg-red-500'}`}></span>
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
                            style={{ width: `${Math.max(0, (currentHP / (selectedEntity.maxHp || 1)) * 100)}%` }}
                        >
                        </div>
                    </div>
                    <span className="ml-2">{currentHP}/{selectedEntity.maxHp}</span>
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

            {isCurrentTurn && actions.length > 0 && (
                <div className="mt-4">
                    <h3 className="font-bold">Actions:</h3>
                    <div className="mt-2 space-y-2">
                        {actions.map(action => (
                            <button
                                key={action.id}
                                className={`w-full text-left p-2 rounded ${
                                    selectedActionId === action.id ? 'bg-blue-500 text-white' : 'bg-gray-100 hover:bg-gray-200'
                                }`}
                                onClick={() => setSelectedActionId(action.id)}
                            >
                                <div className="font-bold">{action.name}</div>
                                <div className="text-sm">{action.description}</div>
                                <div className="text-xs mt-1">{action.actionCost} action{action.actionCost !== 1 ? 's' : ''}</div>
                            </button>
                        ))}
                    </div>
                </div>
            )}

            {selectedAction && targetEntity && (
                <div className="mt-4 p-2 bg-gray-100 rounded">
                    <h3 className="font-bold">Target: {targetEntity.name}</h3>
                    <button
                        className="mt-2 w-full bg-green-500 text-white p-2 rounded flex items-center justify-center"
                        onClick={() => sendAction(selectedEntity.id, selectedAction.id, targetEntity.id)}
                    >
                        <Send size={16} className="mr-2" />
                        Execute {selectedAction.name}
                    </button>
                </div>
            )}
        </div>
    );
};

export default EntityPanel;