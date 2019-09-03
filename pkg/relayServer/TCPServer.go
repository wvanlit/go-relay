package relayServer

import "sync"

func CreateTCPServer(port string) RelayServer {
	return RelayServer{
		relayServer{
			network: "tcp",
			port:    ":" + port,
			clients: []*Client{},
			clientLock: &sync.Mutex{},
		},
	}
}