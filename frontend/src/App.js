import React, { useState, useEffect } from 'react';
import './App.css';
import GameBoard from './components/GameBoard';
import EntityPanel from './components/EntityPanel';
import CombatLog from './components/CombatLog';
import StepControls from './components/StepControls';
import { GameProvider } from './context/GameContext';

function App() {
    return (
        <GameProvider>
            <div className="flex flex-col h-screen bg-gray-100">
                <header className="bg-gray-800 text-white p-4">
                    <h1 className="text-xl font-bold">PF2E Combat Simulator</h1>
                </header>

                <div className="flex flex-1 overflow-hidden">
                    {/* Game Board and Controls */}
                    <div className="flex-1 p-4 overflow-auto flex flex-col">
                        <div className="mb-4 flex-1">
                            <h2 className="text-lg font-bold mb-2">Combat Grid</h2>
                            <GameBoard />
                        </div>
                        <StepControls />
                    </div>

                    {/* Side Panel */}
                    <div className="w-80 bg-white shadow-lg overflow-auto flex flex-col">
                        <EntityPanel />
                        <CombatLog />
                    </div>
                </div>
            </div>
        </GameProvider>
    );
}

export default App;