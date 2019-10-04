package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

type RelayServer struct {
	Port string
	Open bool
}

func CreateServer(port string) *RelayServer {
	return &RelayServer{
		Port: port,
		Open: false,
	}
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
		return RelayConnection{Open: false}
	}
	return RelayConnection{
		Connection: c,
		Open:       true,
		reader:     bufio.NewReader(c),
		writer:     bufio.NewWriter(c),
	}
}

func (r *RelayServer) RunServer() {
	listener := r.getListener()
	conn := r.getConnection(listener)

	for {
		netData := conn.ReceiveMessage()
		if strings.TrimSpace(string(netData)) == "STOP" {
			fmt.Println("Exiting TCP server!")
			return
		}

		fmt.Print("--> ", netData)
		t := time.Now()
		myTime := t.Format(time.RFC3339) + "\n"

		conn.SendMessage(myTime)
	}
}
