package main

import (
	"log"
	"net"
)

const PORT = "8080"

func main() {
	listener, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}
	defer listener.Close()

	log.Printf("tcp server listening on port: %s\n", PORT)

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Printf("error: failed to connect %s\n", err)
			continue
		}

		go handleConnection(connection)
	}
}
