package main

import(
	"github.com/wvanlit/go-relay/pkg/relayServer"
	"time"
)

func main(){
	// Create Server
	relay := relayServer.CreateTCPServer("2019")
	go relay.RunUI()
	go relay.CleanServer(time.Second*15)
	go relay.PurgeServer(time.Minute)
	relay.Run()
}