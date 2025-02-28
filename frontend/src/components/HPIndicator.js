import React from 'react';
import { getHpBarColorClass } from '../services/entityService';

/**
 * HPIndicator displays a health bar for an entity
 * @param {number} currentHP - The entity's current HP
 * @param {number} maxHP - The entity's maximum HP
 */
const HPIndicator = ({ currentHP, maxHP }) => {
  // Calculate percentage (with safety checks)
  const percentage = Math.max(0, ((currentHP || 0) / (maxHP || 1)) * 100);
  
  // Get color class based on entity health status
  const barColor = getHpBarColorClass({ maxHp: maxHP }, currentHP);

  return (
    <div className="absolute bottom-0 left-0 right-0 h-1 bg-gray-300">
      <div
        className={`h-1 ${barColor}`}
        style={{ width: `${percentage}%` }}
      />
    </div>
  );
};

export default HPIndicator;