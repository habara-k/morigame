package src

import (
	"log"

	"github.com/gorilla/websocket"
)

type ClientId string
type Client struct {
	ws   *websocket.Conn
	Send chan []byte
}

type RequestFromClient struct {
	client *Client
	msg    []byte
}

func NewClient(
	ws *websocket.Conn,
	join chan<- *Client,
	leave chan<- *Client,
	process chan<- *RequestFromClient,
) *Client {
	log.Println("NewClient")
	client := &Client{
		ws:   ws,
		Send: make(chan []byte, BUFFSIZE),
	}
	go client.readPump(leave, process)
	go client.writePump()

	join <- client
	return client
}
func (c *Client) Close() {
	c.ws.Close()
}

func (c *Client) readPump(leave chan<- *Client, process chan<- *RequestFromClient) {
	defer func() {
		c.Close()
		close(c.Send)
		leave <- c
	}()
	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			log.Println("Client#read:", err)
			break
		}

		process <- &RequestFromClient{client: c, msg: msg}
	}
	log.Println("Client#read: end")
}
func (c *Client) writePump() {
	for msg := range c.Send {
		w, err := c.ws.NextWriter(websocket.TextMessage)
		if err != nil {
			log.Fatal(err)
		}

		w.Write(msg)

		if err := w.Close(); err != nil {
			log.Fatal(err)
		}
	}
	log.Println("Client#write: end")
}
