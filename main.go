package main

import (
	"fmt"
	"log"
	"net/http"
	"wschat-server/internal/handler"
	"wschat-server/internal/socket"
	"wschat-server/internal/util"
)

var hub = socket.Hub{
	Broadcast:  make(chan socket.Message),
	Register:   make(chan *socket.Client),
	Unregister: make(chan *socket.Client),
	Rooms:      make(map[string]map[*socket.Client]bool),
}

func main() {
	util.LoadEnv()

	go hub.Start()
	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		handler.WebSocketHandler(w, r, &hub)
	})

	log.Printf("Server started on %s:%s", util.Env.Host, util.Env.Port)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", util.Env.Host, util.Env.Port), nil); err != nil {
		log.Fatalf("An Error occured while starting the server: %v", err)
	}
}
