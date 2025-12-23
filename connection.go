package main

import (
	"log"
	"net"
)

func handleConnection(conn net.Conn) {
	log.Printf("new visitor %s\n", conn.RemoteAddr().String())
	msg := []byte("Welcome bratha \n")

	n, err := conn.Write(msg)

	if n < len(msg) || err != nil {
		log.Printf("error: failed to send message %s %d %d \n", err, n, len(msg))
		conn.Close()
		return
	}

	var buffer = make([]byte, 512)
	conn.Read(buffer)

	log.Printf("message: %s\n", string(buffer))

	// conn.Close()
}
