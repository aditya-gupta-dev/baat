package main

import "time"

func main() {
	config := DefaultConfig()

	room := CreateRoom(&config)
	go room.StartListening()
	time.Sleep(10 * time.Second)
}
