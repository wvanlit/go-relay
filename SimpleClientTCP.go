package main

import (
	"fmt"
	"github.com/wvanlit/go-relay/pkg/relayClient"
	"os"
	"time"
)

func main() {
	// Get Client Name and Target Client Name
	name := os.Args[1]
	var target string
	if len(os.Args) > 2{
		target = os.Args[2]
	}

	// Create Client
	client := relayClient.CreateTCPClient(name,"localhost", "2019", 25)

	// Establish Connection to Server
	client.Connection = client.Start()
	//_ = client.Connection.SetReadDeadline(time.Unix(0,1000))
	defer client.Stop()

	// Identify Connection to Server
	client.Identify()

	// Start Pipe
	time.Sleep(time.Second)
	if target != ""{
		client.Send([]byte("conn:"+target))
	}else{
		// Wait before sending through created pipe
		time.Sleep(time.Second)
	}

	// Send Through Pipe
	for {
		time.Sleep(time.Second)
		client.Send([]byte("Hello World! - from "+name))
		time.Sleep(time.Second)
		data := make([]byte, 100)
		n := client.Read(data)
		if n >= 0{
			fmt.Println("Received:",string(data),"on",name,client.Connection.LocalAddr()	)
		}

	}
}