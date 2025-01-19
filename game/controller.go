package game

// Controller interface for entity controllers
type Controller interface {
	NextAction(gs *GameState, entity *Entity) Action
}
