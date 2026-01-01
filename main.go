package main

import (
	"log"
	"os"
)

func main() {

	file, err := os.OpenFile("application.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	log.SetOutput(file)

	config := DefaultConfig()

	room := CreateRoom(&config)
	room.StartListening()
}
