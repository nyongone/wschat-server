package socket

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	Id     string
	Socket *websocket.Conn
	Send   chan []byte
	RoomId string
}

func (c *Client) Read(h *Hub) {
	defer func() {
		h.Unregister <- c
		_ = c.Socket.Close()
	}()

	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			h.Unregister <- c
			_ = c.Socket.Close()
			break
		}

		h.Broadcast <- Message{Sender: c.Id, Content: string(message), RoomId: c.RoomId}
	}
}

func (c *Client) Write() {
	defer func() {
		_ = c.Socket.Close()
	}()

	for message := range c.Send {
		_ = c.Socket.WriteMessage(websocket.TextMessage, message)
	}

	_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
}
