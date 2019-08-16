package relayClient

import (
	"log"
	"net"
)

func (client RelayClient) Start(){
	connection, err := net.Dial(client.network, client.address+":"+client.port)
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()
}
