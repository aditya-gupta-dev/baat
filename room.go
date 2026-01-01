package main

import (
	"fmt"
	"log"
	"net"
)

type Room struct {
	users    []net.Conn
	listener net.Listener
	conn_chn chan net.Conn // connection channel
}

func CreateRoom(config *Config) Room {
	listener, err := net.Listen("tcp", ":"+config.port)
	if err != nil {
		log.Fatalf("failed to open tcp server: %s\n", err)
	}

	log.Printf("started listening on port %s\n", config.port)

	return Room{
		users:    []net.Conn{},
		listener: listener,
		conn_chn: make(chan net.Conn),
	}
}

func connectionHandler(room Room) {
	for conn := range room.conn_chn {

		room.users = append(room.users, conn)

		log.Printf("user-channel-update: new user recieved { %s }\n", conn.RemoteAddr().String())
		log.Printf("user-channel-update: total connections { %d }\n", len(room.users))

		if len(room.users) == 1 {
			continue
		}

		var counter int = 0
		for _, user := range room.users {
			if user.RemoteAddr().String() == conn.LocalAddr().String() {
				continue
			}
			n, err := user.Write([]byte(fmt.Sprintf("user-joined: %s\n", user.RemoteAddr().String())))
			if err != nil {
				log.Fatalf("failed to write %d %s\n", n, err)
			}
			log.Printf("user-send-message: message send to %s\n", user.RemoteAddr().String())
			counter += 1
		}

		log.Printf("user-send-message: total message recievers %d\n", counter)
	}
	log.Printf("user-channel-update: closed\n")
}

func (room *Room) StartListening() {

	go connectionHandler(*room)

	for {
		new_conn, err := room.listener.Accept()
		if err != nil {
			log.Fatalf("failed to connect to %s\n", err)
		}
		room.conn_chn <- new_conn

		log.Printf("user-joined: %s\n", new_conn.RemoteAddr().String())
	}
}

func (room *Room) Close() error {
	return room.listener.Close()
}
