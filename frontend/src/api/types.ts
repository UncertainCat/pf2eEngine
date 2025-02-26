// Types for the API
// These types mirror the backend API models

export interface EntityRef {
  id: string;
  name: string;
}

export interface ActionCardRef {
  id: string;
  name: string;
  description?: string;
  actionCost: number;
  type?: string;
}

export interface EntityState {
  id: string;
  name: string;
  hp: number;
  maxHp: number;
  ac: number;
  actionsRemaining: number;
  reactionsRemaining: number;
  faction: string;
  actionCards?: ActionCardRef[];
  position?: [number, number];
}

export interface GameState {
  entities: EntityState[];
  currentTurn?: string;
  gridWidth: number;
  gridHeight: number;
  round: number;
}

export interface AttackEventData {
  attacker: EntityRef;
  defender: EntityRef;
  roll?: number;
  result?: number;
  degree?: string;
}

export interface DamageEventData {
  source: EntityRef;
  target: EntityRef;
  amount?: number;
  type?: string;
  blocked?: number;
  taken?: number;
}

export interface TurnEventData {
  entity: EntityRef;
}

export interface EventBase {
  type: string;
  version: string;
  timestamp: string;
  message?: string;
  metadata?: Record<string, any>;
}

export interface GameEvent<T = any> extends EventBase {
  data?: T;
}

export interface CommandRequest {
  entity_id: string;
  action_card_id: string;
  params: Record<string, any>;
}

export interface CommandResponse {
  success: boolean;
  message?: string;
  error?: string;
}

// Event type constants (must match backend constants)
export enum EventType {
  INFO = "INFO",
  GAME_SETUP = "GAME_SETUP",
  GAME_STATE = "GAME_STATE",
  ATTACK = "ATTACK",
  ATTACK_RESULT = "ATTACK_RESULT",
  DAMAGE = "DAMAGE",
  DAMAGE_RESULT = "DAMAGE_RESULT",
  TURN_START = "TURN_START",
  TURN_END = "TURN_END",
  ROUND_START = "ROUND_START",
  ROUND_END = "ROUND_END",
  ENTITY_MOVE = "ENTITY_MOVE",
  ENTITY_STATUS = "ENTITY_STATUS",
  ACTION_COMPLETE = "ACTION_COMPLETE"
}