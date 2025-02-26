import React from 'react';
import { useEntity, useActions } from '../hooks';
import { ActionCardRef } from '../types';

interface EntityCardProps {
  entityId: string;
  onTargetSelect?: () => void;
}

const EntityCard: React.FC<EntityCardProps> = ({ entityId, onTargetSelect }) => {
  const { entity, isCurrentTurn } = useEntity(entityId);
  const { sendAction, isSubmitting, error } = useActions();

  // Handle action selection and execution
  const handleActionSelect = async (action: ActionCardRef) => {
    if (!entity) return;
    
    // In a real implementation, you would show a UI to select a target
    // For this example, we're just using a dummy target
    const targetId = "00000000-0000-0000-0000-000000000000";
    
    await sendAction(entity.id, action.id, { targetID: targetId });
  };

  if (!entity) {
    return <div className="p-4 bg-gray-100 rounded shadow">Loading entity...</div>;
  }

  // Calculate health percentage for the progress bar
  const healthPercentage = Math.max(0, Math.min(100, (entity.hp / entity.maxHp) * 100));
  
  // Determine health bar color based on percentage
  let healthColor = 'bg-green-500';
  if (healthPercentage < 30) {
    healthColor = 'bg-red-500';
  } else if (healthPercentage < 70) {
    healthColor = 'bg-yellow-500';
  }

  return (
    <div className={`p-4 bg-white rounded shadow ${isCurrentTurn ? 'ring-2 ring-blue-500' : ''}`}>
      <div className="flex justify-between items-start mb-2">
        <h3 className="text-lg font-bold">{entity.name}</h3>
        <span className="px-2 py-1 text-xs bg-gray-200 rounded">
          {entity.faction}
        </span>
      </div>
      
      {/* Health bar */}
      <div className="mb-4">
        <div className="flex justify-between text-sm mb-1">
          <span>HP</span>
          <span>{entity.hp} / {entity.maxHp}</span>
        </div>
        <div className="w-full bg-gray-200 rounded-full h-2.5">
          <div 
            className={`h-2.5 rounded-full ${healthColor}`} 
            style={{ width: `${healthPercentage}%` }}
          ></div>
        </div>
      </div>
      
      {/* Stats */}
      <div className="grid grid-cols-3 gap-2 mb-4 text-sm">
        <div>
          <span className="font-semibold">AC:</span> {entity.ac}
        </div>
        <div>
          <span className="font-semibold">Actions:</span> {entity.actionsRemaining}
        </div>
        <div>
          <span className="font-semibold">Reactions:</span> {entity.reactionsRemaining}
        </div>
      </div>
      
      {/* Actions */}
      {entity.actionCards && entity.actionCards.length > 0 && (
        <div>
          <h4 className="font-semibold mb-2">Actions</h4>
          <div className="space-y-2">
            {entity.actionCards.map(action => (
              <button
                key={action.id}
                onClick={() => handleActionSelect(action)}
                disabled={isSubmitting || !isCurrentTurn || entity.actionsRemaining < action.actionCost}
                className={`w-full p-2 text-left text-sm rounded ${
                  isCurrentTurn && entity.actionsRemaining >= action.actionCost
                    ? 'bg-blue-100 hover:bg-blue-200'
                    : 'bg-gray-100 cursor-not-allowed opacity-50'
                }`}
              >
                <div className="flex justify-between">
                  <span>{action.name}</span>
                  <span>{action.actionCost} {action.actionCost === 1 ? 'action' : 'actions'}</span>
                </div>
                {action.description && <p className="text-xs text-gray-600 mt-1">{action.description}</p>}
              </button>
            ))}
          </div>
        </div>
      )}
      
      {/* Target button */}
      {onTargetSelect && (
        <button
          onClick={onTargetSelect}
          className="mt-4 w-full p-2 bg-green-100 hover:bg-green-200 rounded"
        >
          Select as Target
        </button>
      )}
      
      {/* Error message */}
      {error && (
        <div className="mt-2 p-2 bg-red-100 text-red-700 text-sm rounded">
          {error.message}
        </div>
      )}
    </div>
  );
};

export default EntityCard;