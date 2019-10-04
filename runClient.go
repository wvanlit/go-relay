package main

import (
	"fmt"
	relayClient "github.com/wvanlit/go-relay/client"
	"os"
)

func main() {
	arguments := os.Args
	if len(arguments) <= 2 {
		fmt.Println("Please provide a host and port number")
		return
	}
	client := relayClient.CreateRelayClient(arguments[1], arguments[2])
	client.RunClient()
}
