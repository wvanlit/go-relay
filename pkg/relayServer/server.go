package relayServer

import (
	"fmt"
	"net"
	"sync"
)

type relayServer struct {
	network string
	port string
	clients []*Client
	clientLock *sync.Mutex
}

type RelayServer struct{
	relayServer
}

type Client struct {
	id          string
	messageSize int
	address     string
	port        string
	requests    chan string
	results     chan net.Conn
	connection  net.Conn
	busy        bool
	close       bool
}

func (server *relayServer) removeClient(id string) error{
	index := -1
	// Find index
	for i, client := range server.clients{
		if client.id == id{
			index = i
		}
	}

	// Check if not found
	if index == -1 {
		return fmt.Errorf("%s not found in client list", id)
	}

	// Move last item to removed index (this removes the item at index)
	server.clients[index] = server.clients[len(server.clients)-1]
	// Set last item to an empty struct
	server.clients[len(server.clients)-1] = &Client{}
	// Shorten slice
	server.clients = server.clients[:len(server.clients)-1]

	return nil
}