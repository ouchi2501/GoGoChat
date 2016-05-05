package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"log"
)

type room struct {
	// forward is message quere
	forward chan []byte
	// join is entry chatroom client canel
	join chan *client
	// leave is chatroom is leave client
	leave chan *client
	// clients is all client
	clients map[*client]bool
}

func newRoom() *room {
	return &room{
		forward:make(chan []byte),
		join:make(chan *client),
		leave:make(chan *client),
		clients:make(map[*client]bool),
	}

}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// join
			r.clients[client] = true
		case client := <-r.leave:
			// leave
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			// all clients send message
			for client := range r.clients {
				select {
				case client.send <- msg:
					// send message
				default:
					delete(r.clients, client)
					close(client.send)
				}

			}
		}

	}
}

const (
	socketBufferSize = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize:socketBufferSize, WriteBufferSize:socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request)  {
	socket, err := upgrader.Upgrade(w,req,nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	client := &client{
		socket:socket,
		send:make(chan []byte, messageBufferSize),
		room:r,
	}
	r.join <- client
	defer func() {r.leave <- client}()
	go client.write()
	client.read()
}