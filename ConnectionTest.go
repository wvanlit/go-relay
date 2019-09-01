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

	for i := 0;i < 750; i++ {
		go createTCPConnection(host, port, &wg)
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
	if !client.Identify(){
		client.Stop()
	}

	for i := 0; i < 20; i++ {
		client.SendString("StressTest")
		time.Sleep(time.Millisecond*250)
	}

	client.Stop()
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
func generateRandomName(n int) string{
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
