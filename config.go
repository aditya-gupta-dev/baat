package main

import "time"

type Config struct {
	Address            string
	MaxConnections     int
	ReadBufferSize     int
	WriteBufferSize    int
	SendChannelSize    int
	BroadcastQueueSize int
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
}

func DefaultConfig() Config {
	return Config{
		Address:            ":8888",
		MaxConnections:     10000,
		ReadBufferSize:     4096,
		WriteBufferSize:    4096,
		SendChannelSize:    256,
		BroadcastQueueSize: 1024,
		ReadTimeout:        5 * time.Minute,
		WriteTimeout:       10 * time.Second,
	}
}
