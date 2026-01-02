package main

import "log"

func main() {

	config := GetConfig()

	logger := NewLogger(config.debug_mode, config.log_bufsize)

	room := NewRoom(&config, &logger)
	room.StartListening()
	if err := room.Close(); err != nil {
		log.Fatalf("failed to close room %s", err)
	}
}
