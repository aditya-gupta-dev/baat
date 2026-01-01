package main

import (
	"log"
	"net"
)

type Room struct {
	clients  []bool
	listener net.Listener
}

func CreateRoom(config *Config) Room {
	listener, err := net.Listen("tcp", ":"+config.port)
	if err != nil {
		log.Fatalf("failed to open tcp server: %s\n", err)
	}

	log.Printf("started listening on port %s\n", config.port)

	return Room{
		clients:  []bool{},
		listener: listener,
	}
}

func (room *Room) StartListening() {
	defer func() {
		err := room.listener.Close()
		if err != nil {
			log.Fatalf("failed to close server %s\n", err)
		}
	}()

	for {
		_, err := room.listener.Accept()
		if err != nil {
			log.Fatalf("failed to connect to %s\n", err)
		}
	}
}
