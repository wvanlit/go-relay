package relayServer

import (
	"bufio"
	"fmt"
	"github.com/wvanlit/go-relay/pkg/relay"
	"os"
	"strings"
	"time"
)

func (server *relayServer) RunUI() {
	fmt.Println("Starting U.I.")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		server.handleCommand(scanner.Text())
	}

	if scanner.Err() != nil {
		fmt.Println(scanner.Err().Error())
	}
}

func (server *relayServer) handleCommand(command string) {
	command = strings.TrimSpace(command)
	switch command {
	case "help":
		fmt.Println(help())

	case "clear":
		go print("\033[H\033[2J") // ANSI Escape codes for clear screen

	case "list":
		go server.listConnections()

	case "remove broken":
		go server.removeBrokenConnections()

	default:
		fmt.Printf("Unknown Command '%s'\n", command)
	}

}

func help() string {
	return "Relay Commands:\n" + "clear -> Clear the screen\n" + "list -> List all current connections\n" + "remove broken -> Remove all broken connections"

}

func (server *relayServer) listConnections() {
	for _, client := range server.clients {
		fmt.Println(client.id, "- busy:", client.busy, "- close:", client.close, "- address:", client.address)
	}
	fmt.Println(len(server.clients), "total clients.")
}

func (server *relayServer) removeBrokenConnections() {
	for _, client := range server.clients {
		if server.isBroken(client) {
			go server.removeConnection(client.id)
		}
		continue
	}
}

func (server *relayServer) isBroken(client *Client) bool {
	if client.busy {
		return false
	}
	if client.connection != nil {
		_ = client.connection.SetWriteDeadline(time.Now().Add(time.Second))
		_, err := client.connection.Write([]byte(relay.CHECK_CONNECTION))
		if err == nil {
			return false
		}
	}

	return true
}

func (server *relayServer) removeConnection(id string) {
	c, err := server.FindClient(id)
	if err != nil {
		return
	}
	select {
	case c.requests <- relay.CLOSE_CONNECTION:
		fmt.Println("Removing", id)
	case <-time.After(time.Second):
		c.close = true
		fmt.Println("Couldn't close", id)
	}
}

func (server *relayServer) CleanServer(interval time.Duration) {
	for {
		select {
		case <-time.After(interval):
			// Put commands to clean server here
			server.removeBrokenConnections()
		}
	}
}

func (server *relayServer) PurgeServer(interval time.Duration) {
	for {
		select {
		case <-time.After(interval):
			// Put commands to clean server here
			server.purgeConnections()
		}
	}
}

func (server *relayServer) purgeConnections() {
	// Force Delete All Connections
	for _, client := range server.clients {
		server.removeConnection(client.id)
		fmt.Println("PURGE")
	}
}
