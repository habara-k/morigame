package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/habara-k/morigame/src"
)

func main() {
	upgrader := websocket.Upgrader{}
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	roomManager := src.NewRoomManager()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}
		ws.SetReadDeadline(time.Now().Add(time.Hour * 24))
		ws.SetWriteDeadline(time.Now().Add(time.Hour * 24))
		roomId := src.RoomId(r.FormValue("room"))
		roomManager.Register <- &src.Registration{RoomId: roomId, Ws: ws}
	})

	http.Handle("/", http.FileServer(http.Dir("client/build")))

	port := "80"
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}
