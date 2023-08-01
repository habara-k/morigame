package game

import (
	"encoding/json"
	"errors"
)

type Action struct {
	Player Player
	Body   ActionBody
}

func (a *Action) Step(state *State, send chan<- *Event) {
	a.Body.Step(state, a.Player, send)
}

type ActionBody interface {
	Step(*State, Player, chan<- *Event)
}

type ActionBodyDiscard struct {
	Card int `json:"card"`
}

type ActionBodyDraw struct{}
type ActionBodyFlip struct{}
type ActionBodyFold struct{}
type ActionBodyMori struct{}
type ActionBodyCounter struct{}
type ActionBodyFetch struct{}

func (a *ActionBodyDiscard) Step(state *State, p Player, send chan<- *Event) {
	state.Discard(p, a.Card, send)
}
func (a *ActionBodyDraw) Step(state *State, p Player, send chan<- *Event) {
	state.Draw(p, send)
}
func (a *ActionBodyFlip) Step(state *State, p Player, send chan<- *Event) {
	state.Flip(p, send)
}
func (a *ActionBodyFold) Step(state *State, p Player, send chan<- *Event) {
	state.Fold(p, send)
}
func (a *ActionBodyMori) Step(state *State, p Player, send chan<- *Event) {
	state.Mori(p, send)
}
func (a *ActionBodyCounter) Step(state *State, p Player, send chan<- *Event) {
	state.Counter(p, send)
}
func (a *ActionBodyFetch) Step(state *State, p Player, send chan<- *Event) {
	state.Fetch(p, send)
}

type actionType string

const (
	actionTypeDiscard actionType = "discard"
	actionTypeDraw    actionType = "draw"
	actionTypeFlip    actionType = "flip"
	actionTypePass    actionType = "pass"
	actionTypeFold    actionType = "fold"
	actionTypeMori    actionType = "mori"
	actionTypeCounter actionType = "counter"
)

type rawAction struct {
	Type actionType `json:"type"`
	body ActionBody `json:"-"`
}

func (a *Action) ParseBody(data []byte) error {
	raw := &rawAction{}
	if err := json.Unmarshal(data, raw); err != nil {
		return err
	}
	var body ActionBody
	switch raw.Type {
	case actionTypeDiscard:
		body = &ActionBodyDiscard{}
	case actionTypeDraw:
		body = &ActionBodyDraw{}
	case actionTypeFlip:
		body = &ActionBodyFlip{}
	case actionTypeFold:
		body = &ActionBodyFold{}
	case actionTypeMori:
		body = &ActionBodyMori{}
	case actionTypeCounter:
		body = &ActionBodyCounter{}
	default:
		return errors.New("Invalid request type")
	}
	if err := json.Unmarshal(data, body); err != nil {
		return err
	}
	a.Body = body
	return nil
}
