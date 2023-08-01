package src

import (
	"log"

	"github.com/gorilla/websocket"
)

type Registration struct {
	RoomId RoomId
	Ws     *websocket.Conn
}
type RoomManager struct {
	rooms     map[RoomId]*Room
	openRoom  chan *Room
	closeRoom chan *Room
	Register  chan *Registration
}

func NewRoomManager() *RoomManager {
	roomManager := &RoomManager{
		rooms:     make(map[RoomId]*Room),
		openRoom:  make(chan *Room, BUFFSIZE),
		closeRoom: make(chan *Room, BUFFSIZE),
		Register:  make(chan *Registration, BUFFSIZE),
	}
	go roomManager.run()
	return roomManager
}
func (m *RoomManager) run() {
	for {
		select {
		case r := <-m.openRoom:
			log.Println("RoomManager#run: openRoom")
			m.rooms[r.Id] = r
		case r := <-m.closeRoom:
			log.Println("RoomManager#run: closeRoom")
			delete(m.rooms, r.Id)
		case r := <-m.Register:
			log.Println("RoomManager#run: register")
			room, ok := m.rooms[r.RoomId]
			if !ok {
				log.Println("RoomManager#run: register: createRoom")
				room = NewRoom(r.RoomId, m.openRoom, m.closeRoom)
			}
			NewClient(r.Ws, room.Join, room.Leave, room.Process)
		}
	}
}
