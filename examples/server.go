package examples

import (
	"fmt"
	relayServer "github.com/wvanlit/go-relay/src/server"
	"os"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide port number")
		return
	}

	server := relayServer.CreateServer(arguments[1])
	server.RunServer()
}
