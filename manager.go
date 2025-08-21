package main

type ClientManager struct {
	rooms      map[string]map[*Client]bool
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
}

func (m *ClientManager) start() {
	for {
		select {
		case conn := <-m.register:
			if connections := m.rooms[conn.roomId]; connections == nil {
				connections = make(map[*Client]bool)
				m.rooms[conn.roomId] = connections
			}
			m.rooms[conn.roomId][conn] = true
		case conn := <-m.unregister:
			connections := m.rooms[conn.roomId]
			if connections != nil {
				if _, ok := connections[conn]; ok {
					delete(connections, conn)
					close(conn.send)
					if len(connections) == 0 {
						delete(m.rooms, conn.roomId)
					}
				}
			}
		case message := <-m.broadcast:
			connections := m.rooms[message.RoomId]
			for c := range connections {
				select {
				case c.send <- []byte(message.Content):
				default:
					close(c.send)
					delete(connections, c)
					if len(connections) == 0 {
						delete(m.rooms, message.RoomId)
					}
				}
			}
		}
	}
}

func (m *ClientManager) send(roomId string, message []byte, ignore *Client) {
	for conn := range m.rooms[roomId] {
		if conn != ignore {
			conn.send <- message
		}
	}
}
