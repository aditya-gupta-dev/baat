package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type Client struct {
	conn     net.Conn
	username string
	send     chan []byte
	hub      *Hub
	reader   *bufio.Reader
	writer   *bufio.Writer
	once     sync.Once
}

func NewClient(conn net.Conn, hub *Hub, bufSize int, sendSize int) *Client {
	return &Client{
		conn:   conn,
		send:   make(chan []byte, sendSize),
		hub:    hub,
		reader: bufio.NewReaderSize(conn, bufSize),
		writer: bufio.NewWriterSize(conn, bufSize),
	}
}

func (c *Client) Start() {
	if err := c.getUsername(); err != nil {
		c.conn.Close()
		return
	}

	c.hub.register <- c

	go c.readPump()
	go c.writePump()
}

func (c *Client) getUsername() error {
	c.conn.SetReadDeadline(time.Now().Add(30 * time.Second))

	c.writer.WriteString("Enter your username: ")
	c.writer.Flush()

	username, err := c.reader.ReadString('\n')
	if err != nil {
		return err
	}

	username = strings.TrimSpace(username)
	username = strings.ReplaceAll(username, "\r", "")
	username = strings.ReplaceAll(username, "\n", "")

	if username == "" {
		username = "anonymous"
	}

	c.username = username
	c.conn.SetReadDeadline(time.Time{}) // Remove deadline
	return nil
}

func (c *Client) readPump() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered in readPump: %v", r)
		}
		c.Close()
	}()

	for {
		c.conn.SetReadDeadline(time.Now().Add(5 * time.Minute))

		line, err := c.reader.ReadString('\n')
		if err != nil {
			return
		}

		line = strings.TrimSpace(line)
		line = strings.ReplaceAll(line, "\r", "")
		line = strings.ReplaceAll(line, "\n", "")

		if line == "" {
			continue
		}

		if strings.ToLower(line) == "/quit" {
			return
		}

		msg := &Message{
			Type:     MsgBroadcast,
			Username: c.username,
			Content:  line,
		}
		c.hub.broadcast <- msg
	}
}

func (c *Client) writePump() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered in writePump: %v", r)
		}
		c.Close()
	}()

	for data := range c.send {
		c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

		if _, err := c.writer.Write(data); err != nil {
			return
		}

		if err := c.writer.Flush(); err != nil {
			return
		}
	}
}

func (c *Client) Send(data []byte) {
	select {
	case c.send <- data:
	default:
		log.Printf("Dropped message for slow client: %s", c.username)
	}
}

func (c *Client) Close() {
	c.once.Do(func() {
		c.hub.unregister <- c
		close(c.send)
		c.conn.Close()
	})
}

func (c *Client) Username() string {
	return c.username
}

func (c *Client) Welcome() {
	welcome := fmt.Sprintf("\r\nWelcome to the chat, %s!\r\n", c.username)
	welcome += "Type your messages and press Enter to send.\r\n"
	welcome += "Type /quit to exit.\r\n\r\n"
	c.Send([]byte(welcome))
}
