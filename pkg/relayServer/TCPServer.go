package relayServer

import "sync"

func CreateTCPServer(port string) RelayServer {
	return RelayServer{
		relayServer{
			network: "tcp",
			port:    ":" + port,
			clients: map[string]Client{},
			clientLock: &sync.Mutex{},
		},
	}
}