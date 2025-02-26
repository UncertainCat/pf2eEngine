import React, { useRef, useEffect } from 'react';
import { useGame } from '../context/GameContext';
import { AlertCircle, Heart, Shield, Sword, SkipForward, CheckCircle } from 'lucide-react';

const CombatLog = () => {
    const { logs } = useGame();
    const logEndRef = useRef(null);

    // Auto-scroll to bottom when new logs are added
    useEffect(() => {
        if (logEndRef.current) {
            logEndRef.current.scrollIntoView({ behavior: 'smooth' });
        }
    }, [logs]);

    // Get icon based on log type
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

    // Get background color based on log type
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

    return (
        <div className="p-4 flex-1 flex flex-col">
            <h2 className="text-lg font-bold mb-2">Combat Log</h2>
            <div className="space-y-2 flex-1 overflow-y-auto">
                {logs.length === 0 ? (
                    <div className="text-gray-500 italic">Combat will begin soon...</div>
                ) : (
                    logs.map((log, index) => (
                        <div
                            key={index}
                            className={`p-2 rounded text-sm flex items-start ${getBackgroundColor(log.type)}`}
                        >
                            {getIcon(log.type)}
                            <div>
                                {log.message || `Event: ${log.type || 'Unknown'}`}
                                {log.data && Object.keys(log.data).length > 0 && (
                                    <div className="text-xs text-gray-600 mt-1">
                                        {Object.entries(log.data).map(([key, value]) => {
                                            // Skip rendering complex nested objects
                                            if (typeof value === 'object' && value !== null) return null;
                                            return (
                                                <span key={key} className="mr-2">
                          {key}: {value}
                        </span>
                                            );
                                        })}
                                    </div>
                                )}
                            </div>
                        </div>
                    ))
                )}
                <div ref={logEndRef} />
            </div>
        </div>
    );
};

export default CombatLog;