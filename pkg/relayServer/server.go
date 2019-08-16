package relayServer

type relayServer struct {
	network string
	port string
	clients map[string]Client
}

type RelayServer struct{
	relayServer
}
