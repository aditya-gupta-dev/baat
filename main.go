package main

import (
	"log"
	"net"
)

const PORT = "8080"
const CONNECTION_LIMIT = 2

func main() {
	listener, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}
	defer listener.Close()

	log.Printf("tcp server listening on port: %s\n", PORT)

	connections := []net.Conn{}

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Printf("error: failed to connect %s\n", err)
			continue
		}

		if len(connections) >= CONNECTION_LIMIT {
			connection.Write([]byte("max visitor count reached.\n"))
			if err := connection.Close(); err != nil {
				log.Fatalf("error: %s", err)
			}
			continue
		}

		connections = append(connections, connection)

		go handleConnection(connection)
	}
}
