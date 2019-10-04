package main

import (
	"fmt"
	relayServer "github.com/wvanlit/go-relay/server"
	"os"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide port number")
		return
	}

	server := relayServer.CreateServer(arguments[1])
	server.RunServer()
}
