package main

import (
	"fmt"
	"net"
)

type Room struct {
	logger   *Logger
	users    []net.Conn
	listener net.Listener
	conn_chn chan net.Conn // connection channel
}

func NewRoom(config *Config, logger *Logger) Room {
	listener, err := net.Listen("tcp", ":"+config.port)
	if err != nil {
		logger.sendErrorf("failed to open tcp server: %s\n", err)
	}

	logger.sendMessageF("started listening on port %s\n", config.port)

	return Room{
		logger:   logger,
		listener: listener,
		users:    []net.Conn{},
		conn_chn: make(chan net.Conn),
	}
}

func connectionHandler(room Room) {
	for conn := range room.conn_chn {

		room.users = append(room.users, conn)

		room.logger.sendMessageF("user-channel-update: new user recieved { %s }\n", conn.RemoteAddr().String())
		room.logger.sendMessageF("user-channel-update: total connections { %d }\n", len(room.users))

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
				room.logger.sendFatalF("failed to write %d %s\n", n, err)
			}
			room.logger.sendMessageF("user-send-message: message send to %s\n", user.RemoteAddr().String())
			counter += 1
		}

		room.logger.sendMessageF("user-send-message: total message recievers %d\n", counter)
	}
	room.logger.sendMessageF("user-channel-update: closed\n")
}

func (room *Room) StartListening() {

	go connectionHandler(*room)

	for {
		new_conn, err := room.listener.Accept()
		if err != nil {
			room.logger.sendFatalF("failed to connect to %s\n", err)
		}
		room.conn_chn <- new_conn

		room.logger.sendMessageF("user-joined: %s\n", new_conn.RemoteAddr().String())
	}
}

func (room *Room) Close() error {
	return room.listener.Close()
}
