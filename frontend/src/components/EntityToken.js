import React, { useEffect, useState, useRef } from 'react';
import { X } from 'lucide-react';
import HPIndicator from './HPIndicator';
import { getEntityFactionClass, calculateAnimatedPosition, isAnimationComplete } from '../services/entityService';

/**
 * EntityToken represents an entity on the game board
 * 
 * @param {Object} entity - The entity data
 * @param {number} currentHP - Current HP value
 * @param {boolean} isSelected - Whether this entity is selected
 * @param {boolean} isTarget - Whether this entity is targeted
 * @param {boolean} isCurrentTurn - Whether it's this entity's turn
 * @param {Function} onClick - Function to call when clicked
 */
const EntityToken = ({ 
  entity, 
  currentHP, 
  isSelected, 
  isTarget, 
  isCurrentTurn, 
  onClick 
}) => {
  // Animation state
  const [position, setPosition] = useState(entity.position || [0, 0]);
  const animationFrameRef = useRef();
  const ANIMATION_DURATION = 500; // ms
  
  // Update position with animation
  useEffect(() => {
    // Skip if no entity or not moving
    if (!entity || !entity.isMoving) {
      setPosition(entity.position || [0, 0]);
      return;
    }
    
    // Animation loop
    const updatePosition = () => {
      // Calculate current animated position
      const animatedPos = calculateAnimatedPosition(entity, ANIMATION_DURATION);
      setPosition(animatedPos);
      
      // Continue animation if not complete
      if (!isAnimationComplete(entity, ANIMATION_DURATION)) {
        animationFrameRef.current = requestAnimationFrame(updatePosition);
      }
    };
    
    // Start animation loop
    animationFrameRef.current = requestAnimationFrame(updatePosition);
    
    // Cleanup animation on unmount or when entity changes
    return () => {
      if (animationFrameRef.current) {
        cancelAnimationFrame(animationFrameRef.current);
      }
    };
  }, [entity, entity?.isMoving, entity?.moveStartTime]);
  
  // Status classes for selection, targeting, and turn indication
  const statusRingClasses = [
    isSelected ? 'ring-2 ring-blue-500' : '',
    isTarget ? 'ring-2 ring-red-500' : '',
    isCurrentTurn ? 'ring-2 ring-yellow-500' : ''
  ].filter(Boolean).join(' ');

  // Get faction color from entity service
  const factionClass = getEntityFactionClass(entity);
  
  // Is entity defeated
  const isDefeated = currentHP <= 0;
  
  // Calculate rendering transform based on animated position
  // We need to transform within the cell, which is positioned absolutely
  const renderX = position[0] - Math.floor(position[0]);
  const renderY = position[1] - Math.floor(position[1]);
  
  // Convert to percentage of cell width/height
  const transformX = renderX * 100;
  const transformY = renderY * 100;
  
  // Apply transform style for smooth animation
  const transformStyle = {
    transform: `translate(${transformX}%, ${transformY}%)`,
    transition: entity?.isMoving ? 'none' : 'transform 0.1s ease-out', // Only apply transition when not actively animating
  };

  return (
    <div
      className={`absolute inset-0 flex items-center justify-center cursor-pointer ${statusRingClasses}`}
      onClick={onClick}
      style={transformStyle}
    >
      <div className={`w-10 h-10 rounded-full flex items-center justify-center text-white font-bold ${factionClass}`}>
        {/* First letter of entity name */}
        {entity.name.charAt(0)}
        
        {/* X mark if defeated */}
        {isDefeated && (
          <div className="absolute inset-0 flex items-center justify-center">
            <X size={24} className="text-black" />
          </div>
        )}
      </div>

      {/* HP bar */}
      <HPIndicator currentHP={currentHP} maxHP={entity.maxHp} />
    </div>
  );
};

export default EntityToken;