// Service for handling combat log events

// Format an event for display in the combat log
export function formatEventForDisplay(event) {
    if (!event || !event.type) return null;

    // Default formatting for any event type
    const formatted = {
        id: event.timestamp + '-' + event.type,
        timestamp: new Date(event.timestamp),
        type: event.type,
        message: event.message || '',
        details: null,
        severity: 'info'
    };

    // Type-specific formatting
    switch (event.type) {
        case 'ATTACK':
        case 'ATTACK_RESULT':
            if (event.data && event.data.attacker && event.data.defender) {
                formatted.message = `${event.data.attacker.name} attacks ${event.data.defender.name}`;
                formatted.details = `Roll: ${event.data.roll}, Result: ${event.data.result}`;
                formatted.severity = event.data.degree === 'CRITICAL_SUCCESS' ? 'critical' : 
                                    event.data.degree === 'SUCCESS' ? 'success' : 
                                    event.data.degree === 'FAILURE' ? 'failure' : 'critical-failure';
            }
            break;

        case 'DAMAGE':
        case 'DAMAGE_RESULT':
            if (event.data && event.data.source && event.data.target) {
                formatted.message = `${event.data.source.name} damages ${event.data.target.name}`;
                
                if (event.data.blocked && event.data.blocked > 0) {
                    formatted.details = `${event.data.amount} damage (${event.data.blocked} blocked, ${event.data.taken} taken)`;
                } else {
                    formatted.details = `${event.data.amount} damage`;
                }
                
                formatted.severity = 'damage';
            }
            break;

        case 'TURN_START':
            if (event.data && event.data.entity) {
                formatted.message = `${event.data.entity.name}'s turn begins`;
                formatted.severity = 'turn';
            }
            break;

        case 'TURN_END':
            if (event.data && event.data.entity) {
                formatted.message = `${event.data.entity.name}'s turn ends`;
                formatted.severity = 'turn';
            }
            break;

        case 'ROUND_START':
            formatted.message = `Round ${event.data.round} begins`;
            formatted.severity = 'round';
            break;

        case 'ROUND_END':
            formatted.message = `Round ${event.data.round} ends`;
            formatted.severity = 'round';
            break;

        case 'ENTITY_MOVE':
            if (event.data && event.data.entity) {
                formatted.message = `${event.data.entity.name} moves to position [${event.data.position[0]}, ${event.data.position[1]}]`;
                formatted.severity = 'movement';
            }
            break;

        case 'ENTITY_STATUS':
            if (event.data && event.data.entity) {
                formatted.message = `${event.data.entity.name} is ${event.data.status.toLowerCase()}`;
                formatted.severity = event.data.status === 'DEAD' ? 'critical' : 'status';
            }
            break;

        case 'ACTION_COMPLETE':
            if (event.data && event.data.entity) {
                formatted.message = `${event.data.entity.name} completes action`;
                if (event.data.actionName) {
                    formatted.message = `${event.data.entity.name} completes ${event.data.actionName}`;
                }
                formatted.severity = 'action';
            }
            break;

        default:
            // For any unhandled event types, just use the message directly
            formatted.message = event.message || `Unknown event type: ${event.type}`;
            break;
    }

    return formatted;
}

// Get event color class based on severity
export function getEventColorClass(severity) {
    switch (severity) {
        case 'critical':
            return 'text-red-500 font-bold';
        case 'success':
            return 'text-green-500';
        case 'failure':
            return 'text-yellow-500';
        case 'critical-failure':
            return 'text-red-400';
        case 'damage':
            return 'text-orange-500';
        case 'turn':
            return 'text-blue-500 font-bold';
        case 'round':
            return 'text-purple-500 font-bold';
        case 'movement':
            return 'text-gray-500';
        case 'action':
            return 'text-indigo-500';
        case 'status':
            return 'text-teal-500';
        default:
            return 'text-gray-700';
    }
}