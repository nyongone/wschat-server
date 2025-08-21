package main

type Message struct {
	Sender  string `json:"sender"`
	Content string `json:"content"`
	RoomId  string `json:"room_id"`
}
