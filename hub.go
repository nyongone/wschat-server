package main

type Hub struct {
	rooms      map[string]map[*Client]bool
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
}

func (h *Hub) start() {
	for {
		select {
		case conn := <-h.register:
			if connections := h.rooms[conn.roomId]; connections == nil {
				connections = make(map[*Client]bool)
				h.rooms[conn.roomId] = connections
			}
			h.rooms[conn.roomId][conn] = true
		case conn := <-h.unregister:
			connections := h.rooms[conn.roomId]
			if connections != nil {
				if _, ok := connections[conn]; ok {
					delete(connections, conn)
					close(conn.send)
					if len(connections) == 0 {
						delete(h.rooms, conn.roomId)
					}
				}
			}
		case message := <-h.broadcast:
			connections := h.rooms[message.RoomId]
			for c := range connections {
				select {
				case c.send <- []byte(message.Content):
				default:
					close(c.send)
					delete(connections, c)
					if len(connections) == 0 {
						delete(h.rooms, message.RoomId)
					}
				}
			}
		}
	}
}

func (h *Hub) send(roomId string, message []byte, ignore *Client) {
	for conn := range h.rooms[roomId] {
		if conn != ignore {
			conn.send <- message
		}
	}
}
