package main

import (
	"log"
	"net/http"

	"github.com/Densuke-fitness/GoWebSocket/trace"
	"github.com/gorilla/websocket"
)

type room struct {
	// foward is a chanel which sends client all recevied messages .
	forward chan []byte

	//this chan is for client who want to join chat room
	join chan *client

	//this chan is for client who want to leave chat room
	leave chan *client

	//by using map (rather than slice), we can maintain object references while consuming less memory.
	clients map[*client]bool

	//trace will recieve log excuted on chat room
	tracer trace.Tracer
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			//participate
			r.clients[client] = true
			r.tracer.Trace("joined new client")
		case client := <-r.leave:
			//leave
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("client left")
		case msg := <-r.forward:
			r.tracer.Trace("recevied message", string(msg))
			//send message to all client
			for client := range r.clients {
				select {
				case client.send <- msg:
					// send message
				default:
					// Failure to send mail
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBUfferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBUfferSize),
		room:   r,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
	}
}
