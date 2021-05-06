package main

import (
	"github.com/gorilla/websocket"
)

//this struct expresses one user chatting .
type client struct {
	socket *websocket.Conn

	send chan []byte
	// chat room which the clinent partcipates in
	room *room
}

func (c *client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
