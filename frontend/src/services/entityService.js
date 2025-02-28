/**
 * Entity Service - Provides utility functions for working with entities
 */

/**
 * Calculate entity status based on HP
 * @param {Object} entity - The entity object
 * @param {number} currentHP - Current HP value (if different from entity.hp)
 * @returns {string} - Status string: 'dead', 'critical', 'wounded', 'injured', or 'healthy'
 */
export function getEntityStatus(entity, currentHP) {
    if (!entity) return 'unknown';
    
    // Use provided currentHP or fall back to entity.hp
    const hp = currentHP !== undefined ? currentHP : entity.hp;
    
    if (hp <= 0) return 'dead';
    
    const hpPercent = (hp / entity.maxHp) * 100;
    
    if (hpPercent <= 25) return 'critical';
    if (hpPercent <= 50) return 'wounded';
    if (hpPercent <= 75) return 'injured';
    return 'healthy';
}

/**
 * Check if entity can take actions
 * @param {Object} entity - The entity object
 * @param {string} currentTurn - ID of the entity whose turn it is
 * @returns {boolean} - True if entity can act
 */
export function canEntityAct(entity, currentTurn) {
    if (!entity) return false;
    
    // Must be this entity's turn
    if (entity.id !== currentTurn) return false;
    
    // Must be alive
    if (entity.hp <= 0) return false;
    
    // Must have actions remaining
    if (entity.actionsRemaining <= 0) return false;
    
    return true;
}

/**
 * Check if an entity can use a specific action card
 * @param {Object} entity - The entity object
 * @param {Object} actionCard - The action card object
 * @param {string} currentTurn - ID of the entity whose turn it is
 * @returns {boolean} - True if entity can use the action card
 */
export function canUseActionCard(entity, actionCard, currentTurn) {
    if (!canEntityAct(entity, currentTurn)) return false;
    if (!actionCard) return false;
    
    // Check if entity has enough actions for this card
    return entity.actionsRemaining >= actionCard.actionCost;
}

/**
 * Get color class for entity HP bar
 * @param {Object} entity - The entity object
 * @param {number} currentHP - Current HP value (if different from entity.hp)
 * @returns {string} - Tailwind CSS class for HP bar color
 */
export function getHpBarColorClass(entity, currentHP) {
    if (!entity) return 'bg-gray-500';
    
    const status = getEntityStatus(entity, currentHP);
    
    switch (status) {
        case 'dead':
            return 'bg-red-800';
        case 'critical':
            return 'bg-red-600';
        case 'wounded':
            return 'bg-orange-500';
        case 'injured':
            return 'bg-yellow-500';
        case 'healthy':
        default:
            return 'bg-green-500';
    }
}

/**
 * Get text color class for entity HP display
 * @param {Object} entity - The entity object
 * @param {number} currentHP - Current HP value (if different from entity.hp)
 * @returns {string} - Tailwind CSS class for text color
 */
export function getHpTextColorClass(entity, currentHP) {
    if (!entity) return 'text-gray-500';
    
    const status = getEntityStatus(entity, currentHP);
    
    switch (status) {
        case 'dead':
            return 'text-red-800';
        case 'critical':
            return 'text-red-600';
        case 'wounded':
            return 'text-orange-500';
        case 'injured':
            return 'text-yellow-500';
        case 'healthy':
        default:
            return 'text-green-500';
    }
}

/**
 * Get background color class for entity token based on faction
 * @param {Object} entity - The entity object
 * @returns {string} - Tailwind CSS class for token background
 */
export function getEntityFactionClass(entity) {
    if (!entity) return 'bg-gray-500';
    
    return entity.faction === 'goodGuys' ? 'bg-blue-500' : 'bg-red-500';
}

/**
 * Get border class for entity based on selection state and faction
 * @param {Object} entity - The entity object
 * @param {boolean} isSelected - Whether entity is selected
 * @param {boolean} isTarget - Whether entity is targeted
 * @param {boolean} isTurn - Whether it's this entity's turn
 * @returns {string} - Tailwind CSS class for border styling
 */
export function getEntityBorderClass(entity, isSelected, isTarget, isTurn) {
    if (!entity) return '';
    
    let classes = [];
    
    // Base faction color
    if (entity.faction === 'goodGuys') {
        classes.push('border-blue-500');
    } else {
        classes.push('border-red-500');
    }
    
    // Selection state
    if (isSelected) {
        classes.push('border-4');
    } else if (isTarget) {
        classes.push('border-4 border-yellow-300');
    } else if (isTurn) {
        classes.push('border-4 border-green-300');
    } else {
        classes.push('border-2');
    }
    
    return classes.join(' ');
}

/**
 * Find entity at specific grid position
 * @param {Array} entities - Array of entities
 * @param {number} x - X coordinate
 * @param {number} y - Y coordinate
 * @returns {Object|null} - Entity at position or null
 */
export function findEntityAtPosition(entities, x, y) {
    if (!entities || !Array.isArray(entities)) return null;
    
    return entities.find(e => 
        e.position && 
        e.position[0] === x && 
        e.position[1] === y
    ) || null;
}

/**
 * Calculate current animated position for an entity
 * @param {Object} entity - The entity object
 * @param {number} duration - Animation duration in milliseconds
 * @returns {Array} - Calculated [x, y] position for rendering
 */
export function calculateAnimatedPosition(entity, duration = 500) {
    if (!entity || !entity.isMoving) {
        return entity?.position || [0, 0];
    }
    
    const startPos = entity.prevPosition || [0, 0];
    const endPos = entity.position || [0, 0];
    const startTime = entity.moveStartTime || 0;
    const currentTime = Date.now();
    const elapsedTime = currentTime - startTime;
    
    // If animation is complete, return final position
    if (elapsedTime >= duration) {
        return endPos;
    }
    
    // Calculate progress (0 to 1)
    const progress = Math.min(elapsedTime / duration, 1);
    
    // Easing function for smoother motion (ease-out)
    const easedProgress = 1 - Math.pow(1 - progress, 2);
    
    // Interpolate position
    const x = startPos[0] + (endPos[0] - startPos[0]) * easedProgress;
    const y = startPos[1] + (endPos[1] - startPos[1]) * easedProgress;
    
    return [x, y];
}

/**
 * Check if an entity's movement animation is complete
 * @param {Object} entity - The entity object
 * @param {number} duration - Animation duration in milliseconds 
 * @returns {boolean} - True if entity has finished moving
 */
export function isAnimationComplete(entity, duration = 500) {
    if (!entity || !entity.isMoving || !entity.moveStartTime) {
        return true;
    }
    
    const currentTime = Date.now();
    const elapsedTime = currentTime - entity.moveStartTime;
    
    return elapsedTime >= duration;
}