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

	room := Room{
		clients:  nil,
		listener: listener,
	}

}
