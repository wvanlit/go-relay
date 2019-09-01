package main

import (
	"fmt"
	"github.com/wvanlit/go-relay/pkg/relayClient"
	"math/rand"
	"sync"
	"time"
)

func main() {
	// Set Seed
	rand.Seed(42)

	// Host
	host := "localhost"
	port := "2019"

	// Create WaitGroup for syncing
	var wg sync.WaitGroup

	// Create Message Sending Connections
	for i := 0; i < 100; i++ {
		go createTCPConnection(host, port, &wg)
		wg.Add(1)
	}

	// Create Pipes
	for i := 0; i < 100; i++ {
		go createPipe(host, port, &wg)
		wg.Add(1)
	}

	wg.Wait()
	fmt.Println("Done")

}

func createTCPConnection(hostname string, port string, group *sync.WaitGroup) {
	defer group.Done()
	name := generateRandomName(5)
	client := relayClient.CreateTCPClient(name, hostname, port, 100)
	client.Start()
	if !client.Identify() {
		client.Stop()
	}

	for i := 0; i < 20; i++ {
		client.SendString("StressTest")
		time.Sleep(time.Millisecond * 250)
	}

	client.Stop()
}

func createPipe(hostname string, port string, group *sync.WaitGroup) {
	defer group.Done()
	// Start Client 1
	c1Name := generateRandomName(6)
	c1 := relayClient.CreateTCPClient(c1Name, hostname, port, 100)

	c1.Start()
	if !c1.Identify() {
		c1.Stop()
		return
	} else {
		defer c1.Stop()
	}

	// Start Client 2
	c2Name := generateRandomName(6)
	c2 := relayClient.CreateTCPClient(c2Name, hostname, port, 100)

	c2.Start()
	if !c2.Identify() {
		c2.Stop()
		return
	} else {
		defer c2.Stop()
	}

	// Start Pipe
	c1.StartPipe(c2Name)
	time.Sleep(time.Second)
	// Send & Receive Messages
	n := 10

	go func() {
		for i := 0; i < n; i++ {
			c1.SendString("PipeTest - 1")
			data := make([]byte, 100)
			c1.ReadForDuration(data, time.Second)
			time.Sleep(time.Second)
		}
	}()

	for i := 0; i < n; i++ {
		c2.SendString("PipeTest - 2")
		data := make([]byte, 100)
		c2.ReadForDuration(data, time.Second)
		time.Sleep(time.Second)
	}

	c1.StopPipe()
	time.Sleep(time.Second*5)
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateRandomName(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
