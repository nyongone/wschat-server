package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var hub = Hub{
	broadcast:  make(chan Message),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	rooms:      make(map[string]map[*Client]bool),
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

	client := &Client{
		id:     uuid.New().String(),
		socket: conn,
		send:   make(chan []byte),
		roomId: r.URL.Query().Get("roomId"),
	}

	hub.register <- client

	go client.read(&hub)
	go client.write()
}

func main() {
	go hub.start()
	http.HandleFunc("/chat", wsHandler)

	log.Println("Server started on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("An Error occured while starting the server: %v", err)
	}
}
