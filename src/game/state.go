package game

import (
	"log"
	"math/rand"
	"time"
)

type Player = int
type Card = int

const (
	N_DECK   = 10000
	N_PLAYER = 4
)

func num(c Card) int {
	return c & 0b1111
}
func suit(c Card) int {
	return c >> 4 & 0b11
}
func revealed(c Card) bool {
	return (c >> 6 & 1) == 1
}
func reveal(c Card) Card {
	return c | 0b1000000
}
func unreveal(c Card) Card {
	return c & ^0b1000000
}
func newCard(num, suit, revealed, deckID int) Card {
	return num | suit<<4 | revealed<<6 | deckID<<7
}

func next(p Player) Player {
	return (p + 1) % N_PLAYER
}
func prev(p Player) Player {
	return (p - 1 + N_PLAYER) % N_PLAYER
}

func newDeck() []Card {
	deck := make([]Card, 0, N_DECK*52)
	for n := 1; n <= 13; n++ {
		for s := 0; s < 4; s++ {
			for id := 0; id < N_DECK; id++ {
				deck = append(deck, newCard(n, s, 0, id))
			}
		}
	}
	shuffle(deck)
	return deck
}

func shuffle(a []Card) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
}

func find(a []Card, x Card) Card {
	for i, y := range a {
		if x == y {
			return i
		}
	}
	return -1
}

type Mode int

const (
	ModeInit Mode = iota
	ModeFlip
	ModeDiscard
	ModeDraw
	ModeLock
	ModeMori
)

type State struct {
	deck      []Card
	trash     []Card
	hand      [N_PLAYER][]Card
	mode      Mode
	turn      Player
	loser     Player
	moriQueue []Player
}

func newState() *State {
	deck := newDeck()
	return &State{
		deck:  deck,
		trash: []Card{},
		hand:  [N_PLAYER][]Card{},
		mode:  ModeInit,
		turn:  -1,
		loser: -1,
	}
}

func (s *State) draw(p Player, send chan<- *Event) {
	card := s.deck[len(s.deck)-1]
	s.deck = s.deck[:len(s.deck)-1]
	s.hand[p] = append(s.hand[p], card)

	for i := 0; i < N_PLAYER; i++ {
		c := -1
		if i == p {
			c = card
		}
		send <- &Event{
			Observer: i,
			Body: &EventBodyDraw{
				player: p,
				card:   c,
			},
		}
	}

	if len(s.deck) == 0 {
		s.shuffle(send)
	}
}

func (s *State) Draw(p Player, send chan<- *Event) {
	if len(s.deck) == 0 {
		log.Fatal("Deck is empty")
	}
	if !(s.mode == ModeInit && len(s.hand[p]) < 3 || (s.mode == ModeDiscard || s.mode == ModeDraw) && s.turn == p) {
		return
	}

	s.draw(p, send)

	if s.mode == ModeInit {
		return
	}

	top := s.trash[len(s.trash)-1]
	if len(s.hand[p]) == 5 {
		canDiscard := false
		for _, c := range s.hand[p] {
			if num(c) == num(top) || suit(c) == suit(top) {
				canDiscard = true
				break
			}
		}
		if canDiscard {
			s.mode = ModeLock
		} else {
			s.reveal(p, send)
		}
	} else if len(s.hand[p]) == 6 {
		canDiscard := false
		for _, c := range s.hand[p] {
			if num(c) == num(top) || suit(c) == suit(top) {
				canDiscard = true
				break
			}
		}
		if canDiscard {
			s.mode = ModeLock
		} else {
			s.burst(p, send)
		}
	} else {
		s.mode = ModeDraw
		s.turn = next(p)
	}
}

func (s *State) shuffle(send chan<- *Event) {
	s.deck = append(s.deck, s.trash[:len(s.trash)-1]...)
	shuffle(s.deck)
	s.trash = s.trash[len(s.trash)-1:]

	for i := 0; i < N_PLAYER; i++ {
		send <- &Event{
			Observer: i,
			Body:     &EventBodyShuffle{},
		}
	}
}

func (s *State) Flip(p Player, send chan<- *Event) {
	if len(s.deck) == 0 {
		log.Fatal("Deck is empty")
	}
	if s.loser != -1 && s.loser != p {
		log.Println("Invalid player")
		return
	}

	switch s.mode {
	case ModeInit:
		for _, h := range s.hand {
			if len(h) < 3 {
				log.Println("Hands less than 3")
				return
			}
		}
	case ModeFlip:
		if len(s.trash) > 0 {
			top := s.trash[len(s.trash)-1]
			for _, h := range s.hand {
				for _, c := range h {
					if num(c) == num(top) {
						log.Println("Someone can disard")
						return
					}
				}
			}
		}
	default:
		return
	}

	card := s.deck[len(s.deck)-1]
	s.deck = s.deck[:len(s.deck)-1]
	s.trash = append(s.trash, unreveal(card))
	if s.mode == ModeInit {
		s.mode = ModeFlip
	}

	for i := 0; i < N_PLAYER; i++ {
		send <- &Event{
			Observer: i,
			Body: &EventBodyFlip{
				card: card,
			},
		}
	}

	if len(s.deck) == 0 {
		s.shuffle(send)
	}
}

func (s *State) Discard(p Player, card Card, send chan<- *Event) {
	if len(s.trash) == 0 {
		return
	}

	top := s.trash[len(s.trash)-1]

	switch s.mode {
	case ModeFlip:
		if num(top) != num(card) {
			return
		}
	case ModeDiscard:
		if !(num(top) == num(card) || suit(top) == suit(card) && s.turn == p) {
			return
		}
	case ModeDraw:
		if !(num(top) == num(card) || suit(top) == suit(card) && (s.turn == p || s.turn == next(p))) {
			return
		}
	case ModeLock:
		if s.turn != p {
			return
		}
		if !(num(top) == num(card) || suit(top) == suit(card)) {
			return
		}
	default:
		return
	}

	idx := find(s.hand[p], card)

	if idx == -1 {
		return
	}

	s.trash = append(s.trash, unreveal(card))
	s.hand[p][idx] = s.hand[p][len(s.hand[p])-1]
	s.hand[p] = s.hand[p][:len(s.hand[p])-1]

	for i := 0; i < N_PLAYER; i++ {
		send <- &Event{
			Observer: i,
			Body: &EventBodyDiscard{
				player: p,
				card:   card,
			},
		}
	}

	if len(s.hand[p]) == 0 {
		s.draw(p, send)
		s.draw(p, send)
	}

	s.turn = next(p)
	s.mode = ModeDiscard
}

func (s *State) reveal(p Player, send chan<- *Event) {
	top := s.trash[len(s.trash)-1]
	for _, h := range s.hand[p] {
		if num(h) == num(top) || suit(h) == suit(top) {
			log.Println("You can discard some cards")
			return
		}
	}

	cards := []Card{} // only unreveald cards
	for i, h := range s.hand[p] {
		if !revealed(h) {
			s.hand[p][i] = reveal(h)
			cards = append(cards, reveal(h))
		}
	}

	for i := 0; i < N_PLAYER; i++ {
		_cards := []Card{}
		for _, c := range cards {
			_cards = append(_cards, c)
		}
		send <- &Event{
			Observer: i,
			Body: &EventBodyReveal{
				player: p,
				cards:  _cards,
			},
		}
	}

	s.turn = next(s.turn)
	s.mode = ModeDraw
}

func (s *State) burst(p Player, send chan<- *Event) {
	top := s.trash[len(s.trash)-1]
	for _, h := range s.hand[p] {
		if num(h) == num(top) || suit(h) == suit(top) {
			log.Println("You can discard some cards")
			return
		}
	}

	for i := 0; i < N_PLAYER; i++ {
		send <- &Event{
			Observer: i,
			Body: &EventBodyBurst{
				player: p,
				card:   s.hand[p][5],
			},
		}
	}

	s.trash = s.trash[:len(s.trash)-1]
	for _, h := range s.hand[p] {
		s.trash = append(s.trash, unreveal(h))
	}
	s.trash = append(s.trash, top)
	s.hand[p] = nil

	s.draw(p, send)
	s.draw(p, send)
	s.draw(p, send)
	s.mode = ModeDraw
	s.turn = next(p)
}

func (s *State) Mori(p Player, send chan<- *Event) {
	if find(s.moriQueue, p) != -1 {
		return
	}
	if !(s.mode == ModeDiscard || s.mode == ModeMori || s.mode == ModeFlip && s.loser != -1) {
		return
	}
	if !can_mori(s.hand[p], s.trash[len(s.trash)-1]) {
		return
	}
	s.mode = ModeMori
	s.moriQueue = append(s.moriQueue, p)

	cards := []Card{} // only unreveald cards
	for i, c := range s.hand[p] {
		if !revealed(c) {
			s.hand[p][i] = reveal(c)
			cards = append(cards, reveal(c))
		}
	}

	for i := 0; i < N_PLAYER; i++ {
		_cards := []Card{}
		for _, c := range cards {
			_cards = append(_cards, c)
		}
		send <- &Event{
			Observer: i,
			Body: &EventBodyMori{
				player: p,
				cards:  _cards,
			},
		}
	}
}

func can_mori(hand []Card, discard Card) bool {
	d_num := num(discard)
	s := 0
	for _, h := range hand {
		s += num(h)
	}
	if s == d_num {
		return true
	}

	if len(hand) != 2 {
		return false
	}

	a, b := num(hand[0]), num(hand[1])
	if a < b {
		a, b = b, a
	}
	return a-b == d_num || a*b == d_num || a%b == 0 && a/b == d_num
}

func (s *State) Fold(p Player, send chan<- *Event) {
	if s.mode != ModeMori {
		return
	}

	loser := prev(s.turn)
	if loser != p {
		return
	}

	for i := 0; i < N_PLAYER; i++ {
		send <- &Event{
			Observer: i,
			Body: &EventBodyFold{
				player: p,
			},
		}
	}

	s.loser = loser
	for _, p := range s.moriQueue {
		for _, h := range s.hand[p] {
			s.trash = append(s.trash, unreveal(h))
		}
		s.hand[p] = nil
	}
	s.moriQueue = nil
	s.mode = ModeInit
}

func (s *State) Counter(p Player, send chan<- *Event) {
	if s.mode != ModeMori {
		return
	}
	winner := prev(s.turn)
	if winner != p {
		return
	}
	if !can_mori(s.hand[p], s.trash[len(s.trash)-1]) {
		return
	}

	s.loser = s.moriQueue[len(s.moriQueue)-1]
	for i := 0; i < N_PLAYER; i++ {
		hand := []Card{}
		for _, c := range s.hand[p] {
			hand = append(hand, c)
		}
		send <- &Event{
			Observer: i,
			Body: &EventBodyCounter{
				player: p,
				cards:  hand,
				loser:  s.loser,
			},
		}
	}

	for _, h := range s.hand[p] {
		s.trash = append(s.trash, unreveal(h))
	}
	s.hand[p] = nil

	for _, i := range s.moriQueue {
		for _, h := range s.hand[i] {
			s.trash = append(s.trash, unreveal(h))
		}
		s.hand[i] = nil
	}
	s.moriQueue = nil
	s.mode = ModeInit
}

func (s *State) Fetch(p Player, send chan<- *Event) {

	top := -1
	if len(s.trash) > 0 {
		top = s.trash[len(s.trash)-1]
	}

	hands := [N_PLAYER][]Card{}
	for i, hand := range s.hand {
		hands[i] = []Card{}
		for _, c := range hand {
			if !revealed(c) && i != p {
				c = -1
			}
			hands[i] = append(hands[i], c)
		}
	}

	moriQueue := []Player{}
	for _, p := range s.moriQueue {
		moriQueue = append(moriQueue, p)
	}

	send <- &Event{
		Observer: p,
		Body: &EventBodyFetch{
			top:       top,
			hands:     hands,
			mode:      s.mode,
			moriQueue: moriQueue,
		},
	}
}
