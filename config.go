package main

type Config struct {
	port             string
	connection_limit int
	max_message_size int
}

func DefaultConfig() Config {
	return Config{
		port:             "8080",
		connection_limit: 12,
		max_message_size: 1024,
	}
}

// TODO: implement getting config from env variables
func GetConfig() Config {
	return DefaultConfig()
}
