import React from 'react';
import GridCell from './GridCell';
import { findEntityAtPosition } from '../services/entityService';

/**
 * GameGrid renders the grid of cells for the game board
 * 
 * @param {number} width - Grid width
 * @param {number} height - Grid height
 * @param {Array} entities - List of all entities
 * @param {Function} getEntityCurrentHP - Function to get entity HP
 * @param {string} selectedEntityId - ID of selected entity
 * @param {string} targetEntityId - ID of targeted entity
 * @param {string} currentEntityId - ID of entity whose turn it is
 * @param {Function} onEntityClick - Function called when entity is clicked
 */
const GameGrid = ({
  width,
  height,
  entities,
  getEntityCurrentHP,
  selectedEntityId,
  targetEntityId,
  currentEntityId,
  onEntityClick
}) => {
  // Generate grid rows and cells
  const renderGrid = () => {
    const grid = [];
    
    for (let y = 0; y < height; y++) {
      const row = [];
      
      for (let x = 0; x < width; x++) {
        // Find entity at this position using the entityService
        const entity = findEntityAtPosition(entities, x, y);
        
        // Get current HP for this entity if one exists
        const currentHP = entity ? getEntityCurrentHP(entity.id) : 0;
        
        row.push(
          <GridCell
            key={`${x}-${y}`}
            x={x}
            y={y}
            entity={entity}
            currentHP={currentHP}
            selectedEntityId={selectedEntityId}
            targetEntityId={targetEntityId}
            currentEntityId={currentEntityId}
            onCellClick={onEntityClick}
          />
        );
      }
      
      grid.push(
        <div key={`row-${y}`} className="flex">
          {row}
        </div>
      );
    }
    
    return grid;
  };

  return (
    <div className="inline-block border border-gray-400">
      {renderGrid()}
    </div>
  );
};

export default GameGrid;