package relayServer

import (
	"fmt"
	"log"
	"net"
)

func (server relayServer) Start() {
	fmt.Println("Starting Server")
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

		// Print Address
		go func(conn net.Conn) {
			println(conn.LocalAddr().String())
		}(conn)
	}

}
