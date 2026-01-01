package main

import (
	"fmt"
	"log"
	"net"
)

type Room struct {
	clients  []net.Conn
	listener net.Listener
}

func CreateRoom(config *Config) Room {
	listener, err := net.Listen("tcp", ":"+config.port)
	if err != nil {
		log.Fatalf("failed to open tcp server: %s\n", err)
	}

	log.Printf("started listening on port %s\n", config.port)

	return Room{
		clients:  []net.Conn{},
		listener: listener,
	}
}

var connections []net.Conn = []net.Conn{}

func connectionHandler(conn_chn chan net.Conn) {
	for conn := range conn_chn {

		log.Printf("user-channel-update: new user recieved { %s }\n", conn.RemoteAddr().String())
		log.Printf("user-channel-update: total connections { %d }\n", len(connections))

		if len(connections) == 1 {
			continue
		}

		var counter int = 0
		for _, connection := range connections {
			if connection.RemoteAddr().String() == conn.LocalAddr().String() {
				continue
			}
			n, err := connection.Write([]byte(fmt.Sprintf("user-joined: %s\n", connection.RemoteAddr().String())))
			if err != nil {
				log.Fatalf("failed to write %d %s\n", n, err)
			}
			log.Printf("user-send-message: message send to %s\n", connection.RemoteAddr().String())
			counter += 1
		}

		log.Printf("user-send-message: total message recievers %d\n", counter)
	}
}

func (room *Room) StartListening() {
	conn_chn := make(chan net.Conn)

	go connectionHandler(conn_chn)

	for {
		new_conn, err := room.listener.Accept()
		if err != nil {
			log.Fatalf("failed to connect to %s\n", err)
		}
		conn_chn <- new_conn

		log.Printf("user-joined: %s\n", new_conn.RemoteAddr().String())
	}
}

func (room *Room) Close() error {
	return room.listener.Close()
}
