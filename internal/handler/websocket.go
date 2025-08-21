package handler

import (
	"log"
	"net/http"
	"wschat-server/internal/socket"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request, hub *socket.Hub) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &socket.Client{
		Id:     uuid.New().String(),
		Socket: conn,
		Send:   make(chan []byte),
		RoomId: r.URL.Query().Get("roomId"),
	}

	hub.Register <- client

	go client.Read(hub)
	go client.Write()
}
