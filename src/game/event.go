package game

import "encoding/json"

type Event struct {
	Observer Player
	Body     EventBody
}
type eventType = string
type EventBody interface {
	eventType() eventType
}

type EventBodyDiscard struct {
	player Player
	card   Card
}
type EventBodyDraw struct {
	player Player
	card   Card
}
type EventBodyShuffle struct{}
type EventBodyFlip struct {
	card Card
}
type EventBodyReveal struct {
	player Player
	cards  []Card
}
type EventBodyBurst struct {
	player Player
	card   Card
}
type EventBodyMori struct {
	player Player
	cards  []Card // only unrevealed cards
}
type EventBodyFold struct {
	player Player
}
type EventBodyCounter struct {
	player Player
	cards  []Card // only unrevealed cards
	loser  Player
}
type EventBodyFetch struct {
	top       Card
	hands     [N_PLAYER][]Card
	mode      Mode
	moriQueue []Player
}

func (e *EventBodyDiscard) eventType() eventType {
	return "discard"
}
func (e *EventBodyDraw) eventType() eventType {
	return "draw"
}
func (e *EventBodyShuffle) eventType() eventType {
	return "shuffle"
}
func (e *EventBodyFlip) eventType() eventType {
	return "flip"
}
func (e *EventBodyReveal) eventType() eventType {
	return "reveal"
}
func (e *EventBodyBurst) eventType() eventType {
	return "burst"
}
func (e *EventBodyMori) eventType() eventType {
	return "mori"
}
func (e *EventBodyFold) eventType() eventType {
	return "fold"
}
func (e *EventBodyCounter) eventType() eventType {
	return "counter"
}
func (e *EventBodyFetch) eventType() eventType {
	return "fetch"
}

func (e *EventBodyDiscard) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type   eventType `json:"type"`
		Player Player    `json:"player"`
		Card   Card      `json:"card"`
	}{
		Type:   e.eventType(),
		Player: e.player,
		Card:   e.card,
	})
}
func (e *EventBodyDraw) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type   eventType `json:"type"`
		Player Player    `json:"player"`
		Card   Card      `json:"card"`
	}{
		Type:   e.eventType(),
		Player: e.player,
		Card:   e.card,
	})
}

func (e *EventBodyShuffle) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type eventType `json:"type"`
	}{
		Type: e.eventType(),
	})
}
func (e *EventBodyFlip) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type eventType `json:"type"`
		Card Card      `json:"card"`
	}{
		Type: e.eventType(),
		Card: e.card,
	})
}
func (e *EventBodyReveal) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type   eventType `json:"type"`
		Player Player    `json:"player"`
		Cards  []Card    `json:"cards"`
	}{
		Type:   e.eventType(),
		Player: e.player,
		Cards:  e.cards,
	})
}
func (e *EventBodyBurst) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type   eventType `json:"type"`
		Player Player    `json:"player"`
		Card   Card      `json:"card"`
	}{
		Type:   e.eventType(),
		Player: e.player,
		Card:   e.card,
	})
}
func (e *EventBodyMori) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type   eventType `json:"type"`
		Player Player    `json:"player"`
		Cards  []Card    `json:"cards"`
	}{
		Type:   e.eventType(),
		Player: e.player,
		Cards:  e.cards,
	})
}
func (e *EventBodyFold) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type   eventType `json:"type"`
		Player Player    `json:"player"`
	}{
		Type:   e.eventType(),
		Player: e.player,
	})
}
func (e *EventBodyCounter) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type   eventType `json:"type"`
		Player Player    `json:"player"`
		Cards  []Card    `json:"cards"`
		Loser  Player    `json:"loser"`
	}{
		Type:   e.eventType(),
		Player: e.player,
		Cards:  e.cards,
		Loser:  e.loser,
	})
}
func (e *EventBodyFetch) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type      eventType        `json:"type"`
		Top       Card             `json:"top"`
		Hands     [N_PLAYER][]Card `json:"hands"`
		Mode      Mode             `json:"mode"`
		MoriQueue []Player         `json:"moriQueue"`
	}{
		Type:      e.eventType(),
		Top:       e.top,
		Hands:     e.hands,
		Mode:      e.mode,
		MoriQueue: e.moriQueue,
	})
}
