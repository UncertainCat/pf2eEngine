import React, { useState } from 'react';
import { useGameState, useEntity, useGameEvents } from '../hooks';
import { EntityState } from '../types';
import EntityCard from './EntityCard';

const GameBoard: React.FC = () => {
  const { gameState, isLoading, error, refresh } = useGameState();
  const { events, pendingEvents, nextEvent, stepMode, toggleStepMode } = useGameEvents();
  const [selectedEntityId, setSelectedEntityId] = useState<string | null>(null);
  const [targetEntityId, setTargetEntityId] = useState<string | null>(null);

  // When selecting an entity as a target
  const handleTargetSelect = (entity: EntityState) => {
    setTargetEntityId(entity.id);
  };

  if (isLoading) {
    return <div className="p-6 text-center">Loading game state...</div>;
  }

  if (error) {
    return (
      <div className="p-6 text-center">
        <div className="p-4 bg-red-100 text-red-700 rounded mb-4">
          {error.message}
        </div>
        <button
          onClick={refresh}
          className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
        >
          Retry
        </button>
      </div>
    );
  }

  if (!gameState) {
    return <div className="p-6 text-center">No game state available</div>;
  }

  // Group entities by faction
  const goodGuys = gameState.entities.filter(e => e.faction === 'goodGuys');
  const badGuys = gameState.entities.filter(e => e.faction === 'badGuys');
  const others = gameState.entities.filter(e => e.faction !== 'goodGuys' && e.faction !== 'badGuys');

  return (
    <div className="container mx-auto p-4">
      {/* Game grid */}
      <div className="mb-8">
        <h2 className="text-xl font-bold mb-4">Combat Grid</h2>
        <div 
          className="grid gap-1 bg-green-100 p-2 rounded" 
          style={{ 
            gridTemplateColumns: `repeat(${gameState.gridWidth}, minmax(0, 1fr))`,
            gridTemplateRows: `repeat(${gameState.gridHeight}, minmax(0, 60px))`
          }}
        >
          {Array.from({ length: gameState.gridWidth * gameState.gridHeight }).map((_, index) => {
            const x = index % gameState.gridWidth;
            const y = Math.floor(index / gameState.gridWidth);
            
            // Find entity at this position
            const entity = gameState.entities.find(e => 
              e.position && e.position[0] === x && e.position[1] === y
            );
            
            return (
              <div 
                key={`cell-${x}-${y}`} 
                className={`border border-green-300 relative ${
                  entity ? 'bg-white' : 'bg-green-50'
                }`}
                onClick={() => entity && setSelectedEntityId(entity.id)}
              >
                {entity && (
                  <div 
                    className={`absolute inset-0 flex items-center justify-center cursor-pointer ${
                      entity.faction === 'goodGuys' ? 'text-blue-600' :
                      entity.faction === 'badGuys' ? 'text-red-600' : 'text-gray-600'
                    } ${entity.id === gameState.currentTurn ? 'font-bold' : ''}`}
                  >
                    {entity.name}
                  </div>
                )}
              </div>
            );
          })}
        </div>
      </div>
      
      {/* Game state info */}
      <div className="mb-8 grid grid-cols-2 gap-4">
        <div className="bg-white p-4 rounded shadow">
          <h3 className="font-bold mb-2">Game Info</h3>
          <div className="text-sm">
            <p>Round: {gameState.round}</p>
            <p>Current Turn: {
              gameState.currentTurn 
                ? gameState.entities.find(e => e.id === gameState.currentTurn)?.name || 'Unknown'
                : 'None'
            }</p>
          </div>
        </div>
        
        <div className="bg-white p-4 rounded shadow">
          <h3 className="font-bold mb-2">Step Controls</h3>
          <div className="flex space-x-2">
            <button
              onClick={toggleStepMode}
              className={`px-3 py-1 rounded ${
                stepMode ? 'bg-blue-500 text-white' : 'bg-gray-200'
              }`}
            >
              {stepMode ? 'Step Mode' : 'Auto Mode'}
            </button>
            
            <button
              onClick={nextEvent}
              disabled={pendingEvents.length === 0}
              className={`px-3 py-1 rounded ${
                pendingEvents.length > 0
                  ? 'bg-green-500 text-white hover:bg-green-600'
                  : 'bg-gray-200 cursor-not-allowed'
              }`}
            >
              Next Step {pendingEvents.length > 0 && `(${pendingEvents.length})`}
            </button>
          </div>
        </div>
      </div>
      
      {/* Selected entity details */}
      {selectedEntityId && (
        <div className="mb-8">
          <h2 className="text-xl font-bold mb-4">Selected Entity</h2>
          <EntityCard entityId={selectedEntityId} />
        </div>
      )}
      
      {/* Combat log */}
      <div className="mb-8">
        <h2 className="text-xl font-bold mb-4">Combat Log</h2>
        <div className="bg-white p-4 rounded shadow max-h-60 overflow-y-auto">
          {events.length === 0 ? (
            <p className="text-gray-500 italic">No events yet</p>
          ) : (
            <div className="space-y-2">
              {events.map((event, index) => (
                <div key={`event-${index}`} className="text-sm border-b pb-2">
                  <span className="font-semibold">{event.type}:</span>{' '}
                  <span>{event.message || 'No message'}</span>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
      
      {/* Entity lists by faction */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Good guys */}
        <div>
          <h2 className="text-xl font-bold mb-4 text-blue-600">Allies</h2>
          <div className="space-y-4">
            {goodGuys.map(entity => (
              <div key={entity.id} onClick={() => setSelectedEntityId(entity.id)}>
                <EntityCard 
                  entityId={entity.id} 
                  onTargetSelect={() => handleTargetSelect(entity)} 
                />
              </div>
            ))}
          </div>
        </div>
        
        {/* Bad guys */}
        <div>
          <h2 className="text-xl font-bold mb-4 text-red-600">Enemies</h2>
          <div className="space-y-4">
            {badGuys.map(entity => (
              <div key={entity.id} onClick={() => setSelectedEntityId(entity.id)}>
                <EntityCard 
                  entityId={entity.id}
                  onTargetSelect={() => handleTargetSelect(entity)}
                />
              </div>
            ))}
          </div>
        </div>
        
        {/* Others */}
        {others.length > 0 && (
          <div className="md:col-span-2">
            <h2 className="text-xl font-bold mb-4 text-gray-600">Neutrals</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {others.map(entity => (
                <div key={entity.id} onClick={() => setSelectedEntityId(entity.id)}>
                  <EntityCard 
                    entityId={entity.id}
                    onTargetSelect={() => handleTargetSelect(entity)}
                  />
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default GameBoard;