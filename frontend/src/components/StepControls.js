import React, { useState } from 'react';
import { ChevronRight, PlayCircle, PauseCircle, RefreshCw, Clock, RotateCcw } from 'lucide-react';
import { useGame } from '../context/GameContext';

const StepControls = () => {
    const {
        nextStep,
        hasNextEvent,
        pendingEvents,
        processedEvents,
        connected,
        historyLoaded,
        loadMoreHistory,
        resetGameState
    } = useGame();
    
    const [loadingMore, setLoadingMore] = useState(false);
    const [resetting, setResetting] = useState(false);
    
    // Function to handle loading more history
    const handleLoadMore = async () => {
        setLoadingMore(true);
        await loadMoreHistory();
        setLoadingMore(false);
    };
    
    // Function to handle resetting the game state
    const handleReset = async () => {
        setResetting(true);
        await resetGameState();
        setResetting(false);
    };
    
    return (
        <div className="border-t border-gray-300 p-4 bg-white">
            <div className="flex items-center justify-between">
                <div className="flex items-center space-x-2">
                    <button
                        className={`flex items-center px-4 py-2 rounded ${
                            hasNextEvent ? 'bg-green-500 text-white hover:bg-green-600' : 'bg-gray-200 text-gray-500 cursor-not-allowed'
                        }`}
                        onClick={nextStep}
                        disabled={!hasNextEvent}
                    >
                        <ChevronRight size={18} className="mr-2" />
                        Next Step {pendingEvents.length > 0 ? `(${pendingEvents.length})` : ''}
                    </button>
                    
                    {processedEvents.length > 0 && (
                        <button
                            className={`flex items-center px-3 py-2 rounded ${
                                loadingMore ? 'bg-gray-300 text-gray-500' : 'bg-blue-500 text-white hover:bg-blue-600'
                            }`}
                            onClick={handleLoadMore}
                            disabled={loadingMore}
                        >
                            {loadingMore ? (
                                <RefreshCw size={16} className="mr-2 animate-spin" />
                            ) : (
                                <Clock size={16} className="mr-2" />
                            )}
                            Load Earlier
                        </button>
                    )}
                    
                    <button
                        className={`flex items-center px-3 py-2 rounded ${
                            resetting ? 'bg-gray-300 text-gray-500' : 'bg-yellow-500 text-white hover:bg-yellow-600'
                        }`}
                        onClick={handleReset}
                        disabled={resetting || !historyLoaded}
                    >
                        {resetting ? (
                            <RefreshCw size={16} className="mr-2 animate-spin" />
                        ) : (
                            <RotateCcw size={16} className="mr-2" />
                        )}
                        Reset
                    </button>
                </div>

                <div className="flex items-center text-sm text-gray-600">
                    {historyLoaded ? (
                        <>
                            <div className="mr-4">
                                <span className="font-semibold">Events:</span> {processedEvents.length} processed, {pendingEvents.length} pending
                            </div>
                            <div>
                                {connected ? (
                                    <span className="text-green-600">Connected</span>
                                ) : (
                                    <span className="text-red-600">Disconnected</span>
                                )}
                            </div>
                        </>
                    ) : (
                        <div className="flex items-center">
                            <RefreshCw size={16} className="mr-2 animate-spin" />
                            <span>Loading history...</span>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
};

export default StepControls;