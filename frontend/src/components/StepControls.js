import React from 'react';
import { ChevronRight, PlayCircle, PauseCircle } from 'lucide-react';
import { useGame } from '../context/GameContext';

const StepControls = () => {
    const {
        stepMode,
        toggleStepMode,
        nextStep,
        hasNextStep,
        currentStepIndex,
        pendingSteps,
        connected
    } = useGame();

    return (
        <div className="border-t border-gray-300 p-4 bg-white">
            <div className="flex items-center justify-between">
                <div className="flex items-center">
                    <button
                        className={`mr-4 flex items-center px-4 py-2 rounded ${
                            stepMode ? 'bg-blue-500 text-white' : 'bg-gray-200 text-gray-800'
                        }`}
                        onClick={toggleStepMode}
                    >
                        {stepMode ? (
                            <>
                                <PauseCircle size={18} className="mr-2" />
                                Step Mode
                            </>
                        ) : (
                            <>
                                <PlayCircle size={18} className="mr-2" />
                                Auto Mode
                            </>
                        )}
                    </button>

                    {stepMode && (
                        <button
                            className={`flex items-center px-4 py-2 rounded ${
                                hasNextStep ? 'bg-green-500 text-white' : 'bg-gray-200 text-gray-500'
                            }`}
                            onClick={nextStep}
                            disabled={!hasNextStep}
                        >
                            <ChevronRight size={18} className="mr-2" />
                            Next Step
                        </button>
                    )}
                </div>

                <div className="text-sm text-gray-600">
                    {connected ? (
                        <span className="text-green-600">Connected</span>
                    ) : (
                        <span className="text-red-600">Disconnected</span>
                    )}

                    {stepMode && pendingSteps.length > 0 && (
                        <span className="ml-4">
              Steps: {currentStepIndex} ({pendingSteps.length} pending)
            </span>
                    )}
                </div>
            </div>

            {stepMode && pendingSteps.length > 0 && (
                <div className="mt-2 p-2 bg-gray-100 rounded text-sm text-gray-700">
                    Next: {pendingSteps[0]?.message || "Unknown action"}
                </div>
            )}
        </div>
    );
};

export default StepControls;