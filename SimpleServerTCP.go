package main

import(
	"github.com/wvanlit/go-relay/pkg/relayServer"
)

func main(){
	// Create Server
	relay := relayServer.CreateTCPServer("2019")
	go relay.RunUI()
	relay.Run()
}