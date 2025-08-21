package main

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	id     string
	socket *websocket.Conn
	send   chan []byte
	roomId string
}

func (c *Client) read(h *Hub) {
	defer func() {
		h.unregister <- c
		_ = c.socket.Close()
	}()

	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			h.unregister <- c
			_ = c.socket.Close()
			break
		}

		h.broadcast <- Message{Sender: c.id, Content: string(message), RoomId: c.roomId}
	}
}

func (c *Client) write() {
	defer func() {
		_ = c.socket.Close()
	}()

	for message := range c.send {
		_ = c.socket.WriteMessage(websocket.TextMessage, message)
	}

	_ = c.socket.WriteMessage(websocket.CloseMessage, []byte{})
}
