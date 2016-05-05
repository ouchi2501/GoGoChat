package main

import "github.com/gorilla/websocket"

type client struct {
	// socket is this client webSockets
	socket *websocket.Conn
	// send is send message chanel
	send chan []byte
	// room is this client join chat room
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
