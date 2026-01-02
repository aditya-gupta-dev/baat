package main

import "fmt"

type MessageType int

const (
	MsgBroadcast MessageType = iota
	MsgJoin
	MsgLeave
)

type Message struct {
	Type     MessageType
	Username string
	Content  string
}

func (m *Message) Format() string {
	switch m.Type {
	case MsgBroadcast:
		return fmt.Sprintf("[%s] %s\r\n", m.Username, m.Content)
	case MsgJoin:
		return fmt.Sprintf("*** %s joined the chat\r\n", m.Username)
	case MsgLeave:
		return fmt.Sprintf("*** %s left the chat\r\n", m.Username)
	default:
		return ""
	}
}

func (m *Message) Bytes() []byte {
	return []byte(m.Format())
}
