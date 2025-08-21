package main

import (
	"fmt"
	"log"
	"net/http"
	"wschat-server/internal/socket"
	"wschat-server/internal/util"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var hub = socket.Hub{
	Broadcast:  make(chan socket.Message),
	Register:   make(chan *socket.Client),
	Unregister: make(chan *socket.Client),
	Rooms:      make(map[string]map[*socket.Client]bool),
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
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

	go client.Read(&hub)
	go client.Write()
}

func main() {
	util.LoadEnv()

	go hub.Start()
	http.HandleFunc("/chat", wsHandler)

	log.Printf("Server started on %s:%s", util.Env.Host, util.Env.Port)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", util.Env.Host, util.Env.Port), nil); err != nil {
		log.Fatalf("An Error occured while starting the server: %v", err)
	}
}
