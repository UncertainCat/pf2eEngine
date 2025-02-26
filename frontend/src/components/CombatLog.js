import React, { useRef, useEffect } from 'react';
import { useGame } from '../context/GameContext';
import { AlertCircle, Heart, Shield, Sword, SkipForward, CheckCircle, RefreshCw } from 'lucide-react';

const CombatLog = () => {
    const { processedEvents, pendingEvents, historyLoaded } = useGame();
    const logEndRef = useRef(null);

    // Auto-scroll to bottom when new events are added
    useEffect(() => {
        if (logEndRef.current) {
            logEndRef.current.scrollIntoView({ behavior: 'smooth' });
        }
    }, [processedEvents]);

    // Get icon based on event type
    const getIcon = (type) => {
        const typeLower = (type || '').toLowerCase();

        // Check for damage-related events
        if (typeLower.includes('damage')) {
            return <Heart size={16} className="text-red-500 mr-2 flex-shrink-0" />;
        }
        // Check for attack-related events
        else if (typeLower.includes('attack')) {
            return <Sword size={16} className="text-yellow-600 mr-2 flex-shrink-0" />;
        }
        // Check for defense-related events
        else if (typeLower.includes('defense') || typeLower.includes('block')) {
            return <Shield size={16} className="text-blue-500 mr-2 flex-shrink-0" />;
        }
        // Check for turn-related events
        else if (typeLower.includes('turn')) {
            return <SkipForward size={16} className="text-purple-500 mr-2 flex-shrink-0" />;
        }
        // Check for success events
        else if (typeLower.includes('success')) {
            return <CheckCircle size={16} className="text-green-500 mr-2 flex-shrink-0" />;
        }
        // Default for unknown events
        return <AlertCircle size={16} className="text-gray-500 mr-2 flex-shrink-0" />;
    };

    // Get background color based on event type
    const getBackgroundColor = (type) => {
        const typeLower = (type || '').toLowerCase();

        // Check for damage-related events
        if (typeLower.includes('damage')) {
            return 'bg-red-100';
        }
        // Check for attack-related events
        else if (typeLower.includes('attack')) {
            return 'bg-yellow-100';
        }
        // Check for defense-related events
        else if (typeLower.includes('defense') || typeLower.includes('block')) {
            return 'bg-blue-100';
        }
        // Check for turn-related events
        else if (typeLower.includes('turn')) {
            return 'bg-purple-100';
        }
        // Check for success events
        else if (typeLower.includes('success')) {
            return 'bg-green-100';
        }
        // Default for unknown events
        return 'bg-gray-100';
    };

    // Format data object for display
    const formatEventData = (data) => {
        if (!data || typeof data !== 'object') return null;
        
        // Filter out complex nested objects and format display
        return Object.entries(data)
            .filter(([key, value]) => {
                // Filter out null values and complex objects (except entity references)
                if (value === null) return false;
                if (typeof value === 'object') {
                    // Keep entity references which have id and name
                    return value.id && value.name;
                }
                return true;
            })
            .map(([key, value]) => {
                // Format entity references
                if (typeof value === 'object' && value.name) {
                    return (
                        <span key={key} className="mr-2">
                            <span className="font-semibold">{key}:</span> {value.name}
                        </span>
                    );
                }
                
                // Format primitive values
                return (
                    <span key={key} className="mr-2">
                        <span className="font-semibold">{key}:</span> {value.toString()}
                    </span>
                );
            });
    };

    // Render the next pending event preview
    const renderNextEventPreview = () => {
        if (pendingEvents.length === 0) return null;
        
        const nextEvent = pendingEvents[0];
        return (
            <div className="mt-4 p-3 border border-dashed border-gray-300 rounded">
                <h3 className="text-sm font-bold mb-1">Next Event:</h3>
                <div className={`p-2 rounded text-sm flex items-start ${getBackgroundColor(nextEvent.type)}`}>
                    {getIcon(nextEvent.type)}
                    <div>
                        <div className="font-medium">{nextEvent.message || nextEvent.type}</div>
                        {nextEvent.data && (
                            <div className="text-xs text-gray-600 mt-1 flex flex-wrap">
                                {formatEventData(nextEvent.data)}
                            </div>
                        )}
                    </div>
                </div>
            </div>
        );
    };

    return (
        <div className="p-4 flex-1 flex flex-col">
            <h2 className="text-lg font-bold mb-2">Combat Log</h2>
            <div className="space-y-2 flex-1 overflow-y-auto max-h-[400px] border border-gray-200 rounded p-2">
                {!historyLoaded ? (
                    <div className="flex items-center justify-center h-full text-gray-500">
                        <RefreshCw size={18} className="animate-spin mr-2" />
                        <span>Loading combat history...</span>
                    </div>
                ) : processedEvents.length === 0 ? (
                    <div className="text-gray-500 italic">Combat will begin soon...</div>
                ) : (
                    processedEvents.map((event, index) => (
                        <div
                            key={`event-${index}`}
                            className={`p-2 rounded text-sm flex items-start ${getBackgroundColor(event.type)}`}
                        >
                            {getIcon(event.type)}
                            <div>
                                <div className="font-medium">{event.message || event.type}</div>
                                {event.data && (
                                    <div className="text-xs text-gray-600 mt-1 flex flex-wrap">
                                        {formatEventData(event.data)}
                                    </div>
                                )}
                            </div>
                        </div>
                    ))
                )}
                <div ref={logEndRef} />
            </div>
            
            {/* Next event preview */}
            {pendingEvents.length > 0 && renderNextEventPreview()}
        </div>
    );
};

export default CombatLog;