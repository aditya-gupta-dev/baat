package main

import (
	"log"
	"net"
	"sync"
)

type Server struct {
	config   Config
	hub      *Hub
	listener net.Listener
	wg       sync.WaitGroup
	done     chan struct{}
}

func NewServer(config Config) *Server {
	hub := NewHub(config.BroadcastQueueSize, config.MaxConnections)

	return &Server{
		config: config,
		hub:    hub,
		done:   make(chan struct{}),
	}
}

func (s *Server) Start() error {
	var err error
	s.listener, err = net.Listen("tcp", s.config.Address)
	if err != nil {
		return err
	}

	log.Printf("Server listening on %s", s.config.Address)
	log.Printf("Max connections: %d", s.config.MaxConnections)

	go s.hub.Run()

	s.acceptConnections()
	return nil
}

func (s *Server) acceptConnections() {
	for {
		select {
		case <-s.done:
			return
		default:
		}

		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.done:
				return
			default:
				log.Printf("Accept error: %v", err)
				continue
			}
		}

		s.wg.Add(1)
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer func() {
		s.wg.Done()
		if r := recover(); r != nil {
			log.Printf("Recovered in handleConnection: %v", r)
		}
	}()

	client := NewClient(
		conn,
		s.hub,
		s.config.ReadBufferSize,
		s.config.SendChannelSize,
	)

	client.Start()
}

func (s *Server) Shutdown() {
	log.Println("Initiating server shutdown...")

	close(s.done)

	if s.listener != nil {
		s.listener.Close()
	}

	s.hub.Shutdown()

	s.wg.Wait()

	log.Println("Server shutdown complete")
}

func (s *Server) Stats() map[string]interface{} {
	return map[string]interface{}{
		"active_connections": s.hub.ClientCount(),
		"max_connections":    s.config.MaxConnections,
	}
}
