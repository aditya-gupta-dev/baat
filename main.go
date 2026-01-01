package main

func main() {

	config := GetConfig()

	logger := NewLogger(config.debug_mode, config.log_bufsize)

	room := NewRoom(&config, &logger)
	room.StartListening()
}
