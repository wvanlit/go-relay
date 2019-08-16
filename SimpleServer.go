package main

import(
	"github.com/wvanlit/go-relay/pkg/relayServer"
)

func main(){
	// Create Server
	relay := relayServer.CreateRelayServer()
	relay.Start()
}