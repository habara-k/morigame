package src

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/habara-k/morigame/src/game"
)

type RoomId string

type Room struct {
	seat    map[*Client]int
	clients [game.N_PLAYER]*Client
	Join    chan *Client
	Leave   chan *Client
	Id      RoomId
	Process chan *RequestFromClient
	Send    chan *game.Event
	game    *game.Game
}

func NewRoom(id RoomId, openRoom chan<- *Room, closeRoom chan<- *Room) *Room {
	send := make(chan *game.Event, BUFFSIZE)
	game := game.NewGame(send)
	room := &Room{
		seat:    make(map[*Client]int),
		Join:    make(chan *Client, BUFFSIZE),
		Leave:   make(chan *Client, BUFFSIZE),
		Id:      id,
		Process: make(chan *RequestFromClient, BUFFSIZE),
		Send:    send,
		game:    game,
	}
	go room.run(closeRoom)

	openRoom <- room
	return room
}
func (r *Room) run(closeRoom chan<- *Room) {
	for {
		select {
		case c := <-r.Join:
			log.Println("Room#run: Join")
			if len(r.seat) >= game.N_PLAYER {
				log.Println("Room#run: Join: full")
				c.Close()
				continue
			}

			id := 0
			for {
				if r.clients[id] == nil {
					break
				}
				id++
			}

			r.seat[c] = id
			r.clients[id] = c
			c.Send <- []byte(fmt.Sprintf(`{"type":"join","id":%v}`, id))
			r.game.Process <- &game.Action{
				Player: id,
				Body:   &game.ActionBodyFetch{},
			}

		case c := <-r.Leave:
			log.Println("Room#run: Leave")
			r.clients[r.seat[c]] = nil
			delete(r.seat, c)
			if len(r.seat) == 0 {
				closeRoom <- r
				break
			}

		case req := <-r.Process:
			a := &game.Action{Player: r.seat[req.client]}
			err := a.ParseBody(req.msg)
			if err != nil {
				log.Println(err)
				continue
			}
			r.game.Process <- a

		case ev := <-r.Send:
			c := r.clients[ev.Observer]
			if c == nil {
				continue
			}
			msg, err := json.Marshal(ev.Body)
			if err != nil {
				log.Fatal(err)
			}
			c.Send <- msg
		}
	}
}
