package socket

type Hub struct {
	Rooms      map[string]map[*Client]bool
	Broadcast  chan Message
	Register   chan *Client
	Unregister chan *Client
}

func (h *Hub) Start() {
	for {
		select {
		case conn := <-h.Register:
			if connections := h.Rooms[conn.RoomId]; connections == nil {
				connections = make(map[*Client]bool)
				h.Rooms[conn.RoomId] = connections
			}
			h.Rooms[conn.RoomId][conn] = true
		case conn := <-h.Unregister:
			connections := h.Rooms[conn.RoomId]
			if connections != nil {
				if _, ok := connections[conn]; ok {
					delete(connections, conn)
					close(conn.Send)
					if len(connections) == 0 {
						delete(h.Rooms, conn.RoomId)
					}
				}
			}
		case message := <-h.Broadcast:
			connections := h.Rooms[message.RoomId]
			for c := range connections {
				select {
				case c.Send <- []byte(message.Content):
				default:
					close(c.Send)
					delete(connections, c)
					if len(connections) == 0 {
						delete(h.Rooms, message.RoomId)
					}
				}
			}
		}
	}
}

func (h *Hub) send(roomId string, message []byte, ignore *Client) {
	for conn := range h.Rooms[roomId] {
		if conn != ignore {
			conn.Send <- message
		}
	}
}
