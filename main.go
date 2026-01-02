package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	addr := flag.String("addr", ":8888", "TCP address to listen on")
	maxConn := flag.Int("max", 10000, "Maximum concurrent connections")
	flag.Parse()

	config := Config{
		Address:            *addr,
		MaxConnections:     *maxConn,
		ReadBufferSize:     4096,
		WriteBufferSize:    4096,
		SendChannelSize:    256,
		BroadcastQueueSize: 1024,
	}

	server := NewServer(config)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Starting chat server on %s", *addr)
		if err := server.Start(); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-sigChan
	log.Println("Shutting down server...")
	server.Shutdown()
	log.Println("Server stopped")
}
