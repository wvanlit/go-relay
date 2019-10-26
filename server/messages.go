package server

import (
	"fmt"
)

type MessageType int

const (
	CLOSE_SERVER MessageType = 0
	CREATE_USER  MessageType = 1
	DELETE_USER  MessageType = 2
	CHECK_USER   MessageType = 3
	REQUEST_PIPE MessageType = 4
	CREATE_PIPE  MessageType = 5
)

type Message interface {
	GetType() MessageType
	ProcessMessage(server *RelayServer)
}

type CloseServerMessage struct {
	messageType MessageType
}

func CreateCloseServerMessage() CloseServerMessage {
	return CloseServerMessage{
		messageType: CLOSE_SERVER,
	}
}

func (csm CloseServerMessage) GetType() MessageType {
	return csm.messageType
}

func (csm CloseServerMessage) ProcessMessage(server *RelayServer) {
	server.Open = false
	fmt.Println("Closing Server")
}
