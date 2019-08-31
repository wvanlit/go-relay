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

	// Identify Connection to Server, close if failed
	if !client.Identify(){
		return
	}

	// Start Pipe
	time.Sleep(time.Second)
	if target != ""{
		client.Send([]byte("PIPE:"+target))
	}else{
		// Wait before sending through created pipe
		time.Sleep(time.Second)
	}

	// Send Through Pipe
	for i := 0; i < 3; i++{
		fmt.Println("Rotation",i)
		time.Sleep(time.Millisecond*500)
		client.Send([]byte("Hello World! - from "+name))
		time.Sleep(time.Millisecond*500)
		data := make([]byte, 100)
		n := client.Read(data)
		if n >= 0{
			fmt.Println("Received:",string(data),"on",name,client.Connection.LocalAddr()	)
		}
	}
	// Stop Pipe
	client.Send([]byte("CLOSE"))
	time.Sleep(time.Second)
}