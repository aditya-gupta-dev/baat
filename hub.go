package main

import (
	"log"
	"sync"
	"sync/atomic"
)

type Hub struct {
	clients    sync.Map // map[*Client]bool
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
	count      int64
	maxClients int
}

func NewHub(broadcastSize int, maxClients int) *Hub {
	return &Hub{
		broadcast:  make(chan *Message, broadcastSize),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		maxClients: maxClients,
	}
}

func (h *Hub) Run() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered in hub.Run: %v", r)
		}
	}()

	for {
		select {
		case client := <-h.register:
			h.handleRegister(client)
		case client := <-h.unregister:
			h.handleUnregister(client)
		case message := <-h.broadcast:
			h.handleBroadcast(message)
		}
	}
}

func (h *Hub) handleRegister(client *Client) {
	if atomic.LoadInt64(&h.count) >= int64(h.maxClients) {
		client.conn.Write([]byte("Server is full. Try again later.\r\n"))
		client.conn.Close()
		return
	}

	h.clients.Store(client, true)
	atomic.AddInt64(&h.count, 1)

	log.Printf("Client registered: %s (total: %d)", client.Username(), atomic.LoadInt64(&h.count))

	client.Welcome()

	joinMsg := &Message{
		Type:     MsgJoin,
		Username: client.Username(),
	}
	h.broadcastToAll(joinMsg)
}

func (h *Hub) handleUnregister(client *Client) {
	if _, ok := h.clients.LoadAndDelete(client); ok {
		atomic.AddInt64(&h.count, -1)
		log.Printf("Client unregistered: %s (total: %d)", client.Username(), atomic.LoadInt64(&h.count))

		leaveMsg := &Message{
			Type:     MsgLeave,
			Username: client.Username(),
		}
		h.broadcastToAll(leaveMsg)
	}
}

func (h *Hub) handleBroadcast(message *Message) {
	h.broadcastToAll(message)
}

func (h *Hub) broadcastToAll(message *Message) {
	data := message.Bytes()

	h.clients.Range(func(key, value interface{}) bool {
		client := key.(*Client)
		client.Send(data)
		return true
	})
}

func (h *Hub) ClientCount() int64 {
	return atomic.LoadInt64(&h.count)
}

func (h *Hub) Shutdown() {
	log.Println("Shutting down hub...")

	h.clients.Range(func(key, value interface{}) bool {
		client := key.(*Client)
		client.Close()
		return true
	})

	close(h.broadcast)
	close(h.register)
	close(h.unregister)
}
