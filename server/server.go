package server

import (
	"bufio"
	"fmt"
	"github.com/wvanlit/go-relay/global"
	"net"
	"sync"
)

type RelayServer struct {
	Port               string
	Open               bool
	ConnectionMessages chan Message
	users              map[string]*RelayConnection
	lock               sync.Mutex
}

func CreateServer(port string) *RelayServer {
	return &RelayServer{
		Port:               port,
		Open:               true,
		ConnectionMessages: make(chan Message, 5),
		lock:               sync.Mutex{},
		users:              map[string]*RelayConnection{},
	}
}

func (r *RelayServer) SetUser(name string, conn *RelayConnection) {
	r.lock.Lock()
	r.users[name] = conn
	r.lock.Unlock()
}

func (r *RelayServer) DeleteUser(name string) {
	r.lock.Lock()
	delete(r.users, name)
	r.lock.Unlock()
}

func (r *RelayServer) CheckForUser(name string) bool {
	_, ok := r.users[name]
	return ok
}

func (r *RelayServer) GetUser(name string) (*RelayConnection, bool) {
	user, ok := r.users[name]
	return user, ok
}

func (r *RelayServer) getListener() net.Listener {
	l, err := net.Listen("tcp", ":"+r.Port)
	if err != nil {
		fmt.Println("Error on listening to connection:", err)
		return nil
	}
	return l
}

func (r *RelayServer) getConnection(listener net.Listener) RelayConnection {
	c, err := listener.Accept()
	if err != nil {
		fmt.Println("Error on accepting connection:", err)
		return RelayConnection{State: global.Offline}
	}
	return RelayConnection{
		Connection:    c,
		State:         global.Open,
		reader:        bufio.NewReader(c),
		writer:        bufio.NewWriter(c),
		input:         make(chan string, 1),
		output:        make(chan string, 1),
		serverChannel: r.ConnectionMessages,
		RequestPipe:   make(chan PipeRequest, 0),
	}
}

func (r *RelayServer) handleConnectionMessages() {
	for {
		select {
		case command := <-r.ConnectionMessages:
			command.ProcessMessage(r)
		}
	}
}

func (r *RelayServer) RunServer() {
	go r.handleConnectionMessages()
	listener := r.getListener()

	// Accept Open Connections
	go func() {
		for r.Open {
			conn := r.getConnection(listener)
			go conn.HandleConnection()
		}
	}()

	for r.Open {
	}
}
