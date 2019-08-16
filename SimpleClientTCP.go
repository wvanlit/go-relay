package main

import (
	"github.com/wvanlit/go-relay/pkg/relayClient"
	"time"
)

func main() {
	// Create Client
	client := relayClient.CreateTCPClient("debug","localhost", "2019", 25)
	client.Connection = client.Start()
	defer client.Stop()
	client.Identify()
	for i := 0; i < 3; i++ {
		time.Sleep(time.Second)
		client.Send([]byte("Hello World!"))
	}
}