package relayServer

import "sync"

type relayServer struct {
	network string
	port string
	clients map[string]Client
	clientLock *sync.Mutex
}

type RelayServer struct{
	relayServer
}
