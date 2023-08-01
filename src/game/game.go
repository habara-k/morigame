package game

type Game struct {
	state   *State
	Process chan *Action
}

func NewGame(send chan<- *Event) *Game {
	game := &Game{
		state:   newState(),
		Process: make(chan *Action),
	}
	go game.run(send)
	return game
}

func (g *Game) run(send chan<- *Event) {
	for a := range g.Process {
		a.Step(g.state, send)
	}
}
