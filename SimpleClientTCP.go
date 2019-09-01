package main

import (
	"fmt"
	"github.com/wvanlit/go-relay/pkg/relay"
	"github.com/wvanlit/go-relay/pkg/relayClient"
	"os"
	"time"
)

func main() {
	// Get Client Name and Target Client Name
	name := os.Args[1]
	var target string
	if len(os.Args) > 2 {
		target = os.Args[2]
	}

	// Create Client
	client := relayClient.CreateTCPClient(name, "localhost", "2019", 25)

	// Establish Connection to Server
	client.Start()

	defer client.Stop()

	// Identify Connection to Server, close if failed
	if !client.Identify() {
		return
	}

	// Start Pipe
	time.Sleep(time.Second)
	if target != "" {
		client.StartPipe(target)
	} else {
		// Wait before sending through created pipe
		time.Sleep(time.Second)
	}

	// Send Through Pipe
	for i := 0; i < 3; i++ {

		fmt.Println("Rotation", i)
		time.Sleep(time.Millisecond * 1000)
		client.SendString("Hello World! - from " + name)
		time.Sleep(time.Millisecond * 1000)

		//data := make([]byte, 100)
		//n := client.Read(data)
		//if n >= 0{
		//	fmt.Println("Received:",string(data),"on",name,client.Connection.LocalAddr()	)
		//}
	}
	time.Sleep(time.Second*5)
	// Stop Pipe
	client.SendString(relay.STOP_PIPE)
	time.Sleep(time.Second*5)

	// Stop Connection
	client.Stop()
	time.Sleep(time.Second*5)
}
