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
		print("\033[H\033[2J") // ANSI Escape codes for clear screen

	case "list":
		server.listConnections()

	case "remove broken":
		server.removeBrokenConnections()

	default:
		fmt.Printf("Unknown Command '%s'\n", command)
	}

}

func help() string {
	return "Relay Commands:\n" + "clear -> Clear the screen\n" + "list -> List all current connections\n" + "remove broken -> Remove all broken connections"

}

func (server *relayServer) listConnections() {
	for id, client := range server.clients {
		fmt.Println(id, "- busy:", client.busy, "- address:", client.address)
	}
}

func (server *relayServer) removeBrokenConnections() {
	for id, client := range server.clients {
		fmt.Println(id)
		if client.connection != nil {
			_ = client.connection.SetWriteDeadline(time.Now().Add(time.Second))
			_, err := client.connection.Write([]byte(relay.CHECK_CONNECTION))
			if err == nil {
				// Short circuit to prevent closing
				return
			}
		}

		client.requests <- relay.CLOSE_CONNECTION
		fmt.Println("Closed", id)

	}
}
