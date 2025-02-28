import React from 'react';
import EntityToken from './EntityToken';

/**
 * GridCell represents a single cell on the game board
 * 
 * @param {number} x - X coordinate
 * @param {number} y - Y coordinate
 * @param {Object} entity - Entity at this position (if any)
 * @param {number} currentHP - Current HP of entity
 * @param {string} selectedEntityId - ID of selected entity
 * @param {string} targetEntityId - ID of targeted entity
 * @param {string} currentEntityId - ID of entity whose turn it is
 * @param {Function} onCellClick - Function called when cell is clicked
 */
const GridCell = ({
  x,
  y,
  entity,
  currentHP,
  selectedEntityId,
  targetEntityId,
  currentEntityId,
  onCellClick
}) => {
  // Set the checkerboard pattern
  const cellColor = (x + y) % 2 === 0 ? 'bg-gray-100' : 'bg-gray-200';

  // Handle cell click
  const handleClick = () => {
    if (entity) {
      onCellClick(entity);
    }
  };

  // Only show entities in this cell if:
  // 1. There is an entity
  // 2. The entity's integer position matches this cell's coordinates
  const shouldShowEntity = entity && 
    Math.floor(entity.position[0]) === x && 
    Math.floor(entity.position[1]) === y;

  return (
    <div
      className={`w-12 h-12 border border-gray-300 relative ${cellColor}`}
      data-x={x}
      data-y={y}
    >
      {shouldShowEntity && (
        <EntityToken
          entity={entity}
          currentHP={currentHP}
          isSelected={entity.id === selectedEntityId}
          isTarget={entity.id === targetEntityId}
          isCurrentTurn={entity.id === currentEntityId}
          onClick={handleClick}
        />
      )}
      
      {/* Debugging coordinates - can be removed in production */}
      <div className="absolute bottom-0 right-0 text-xs text-gray-400 opacity-50">
        {x},{y}
      </div>
    </div>
  );
};

export default GridCell;