package relayServer

import (
	"fmt"
	"log"
	"net"
)

func MessageClient(conn net.Conn, message []byte){
	_, err := conn.Write(message)
	if err != nil{
		log.Println(err)
	}
}

func ReadClientMessage(conn net.Conn, data *[]byte) bool{
	_, err := conn.Read(*data)
	if err != nil{
		log.Println(err)
		return false
	}
	if len(*data) <= 0{
		return false
	}
	return true
}

func (server relayServer) Run() {
	listen, err := net.Listen(server.network, server.port)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server Running on", listen.Addr(), "over", listen.Addr().Network())

	defer listen.Close()
	for {
		// Accept connections
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// Handle Connection
		go server.HandleConnection(conn)
	}
}

