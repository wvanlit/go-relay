package relayClient

import (
	"fmt"
	"log"
	"net"
	"strings"
)

func (client RelayClient) Start() net.Conn {
	// Setup Connection to Address
	conn, err := net.Dial(client.network, client.address+":"+client.port)
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

func (client RelayClient) Stop() {
	// Setup Connection to Address
	err := client.Connection.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func (client RelayClient) Send(message []byte) {
	_, err := client.Connection.Write(message)
	if err != nil {
		log.Println(err)
	}
}

func (client RelayClient) Read(message []byte) int{
	n, err := client.Connection.Read(message)
	// Check if there is an error that is not a timeout
	if err != nil && !strings.Contains(err.Error(), "timeout") {
		log.Println(err.Error())
	}
	return n
}

func (client RelayClient) Identify() bool {
	client.Send([]byte(fmt.Sprint(client.hostname, "|", client.messageSize)))
	data := make([]byte, 100)
	client.Read(data)
	if string(data) != "OK"{
		fmt.Println(string(data))
		return false
	}

	return true
}
