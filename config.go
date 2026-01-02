package main

type Config struct {
	port             string
	connection_limit int
	max_message_size int
	debug_mode       bool
	log_bufsize      int
}

func DefaultConfig() Config {
	return Config{
		port:             "8080",
		connection_limit: 12,
		max_message_size: 1024,
		debug_mode:       true,
		log_bufsize:      12,
	}
}

// TODO: implement getting config from env variables
func GetConfig() Config {
	return DefaultConfig()
}
