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
)

type message interface {
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

type CreateUserMessage struct {
	messageType MessageType
	username    string
	conn        *RelayConnection
}

func CreateCreateUserMessage(name string, conn *RelayConnection) CreateUserMessage {
	return CreateUserMessage{
		messageType: CREATE_USER,
		username:    name,
		conn:        conn,
	}
}

func (cum CreateUserMessage) GetType() MessageType {
	return cum.messageType
}

func (cum CreateUserMessage) ProcessMessage(server *RelayServer) {
	server.SetUser(cum.username, cum.conn)
}

type DeleteUserMessage struct {
	messageType MessageType
	username    string
}

func CreateDeleteUserMessage(name string) DeleteUserMessage {
	return DeleteUserMessage{
		messageType: DELETE_USER,
		username:    name,
	}
}

func (dum DeleteUserMessage) GetType() MessageType {
	return dum.messageType
}

func (dum DeleteUserMessage) ProcessMessage(server *RelayServer) {
	server.DeleteUser(dum.username)
}

type CheckUserMessage struct{
	messageType MessageType
	username string
	feedback chan bool
}

func CreateCheckUserMessage(name string, feedback chan bool) CheckUserMessage{
	return CheckUserMessage{
		messageType: CHECK_USER,
		username:    name,
		feedback:    feedback,
	}
}

func (cum CheckUserMessage) GetType() MessageType {
	return cum.messageType
}

func (cum CheckUserMessage) ProcessMessage(server *RelayServer) {
	cum.feedback <- server.CheckForUser(cum.username)
}