package main

import "github.com/wvanlit/go-relay/pkg/relayClient"

func main() {
	// Create Client
	client := relayClient.CreateTCPClient("localhost", "2019")
	client.Start()
}