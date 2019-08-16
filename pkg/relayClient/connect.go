package relayClient

import (
	"fmt"
	"log"
	"net"
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

func (client RelayClient) Read(message []byte) {
	_, err := client.Connection.Read(message)
	if err != nil {
		log.Println(err)
	}
}

func (client RelayClient) Identify() bool {
	client.Send([]byte(fmt.Sprint(client.hostname, "|", client.messageSize)))
	data := make([]byte, 10)
	client.Read(data)
	if string(data) != "OK"{
		fmt.Println(string(data))
		return false
	}

	return true
}
