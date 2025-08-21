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

func (c *Client) read(m *ClientManager) {
	defer func() {
		m.unregister <- c
		_ = c.socket.Close()
	}()

	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			m.unregister <- c
			_ = c.socket.Close()
			break
		}

		m.broadcast <- &Message{Sender: c.id, Content: string(message), RoomId: c.roomId}
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
