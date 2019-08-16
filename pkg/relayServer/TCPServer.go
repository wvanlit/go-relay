package relayServer

func CreateTCPServer(port string) RelayServer {
	return RelayServer{
		relayServer{
			network: "tcp",
			port:    ":" + port,
		},
	}
}	
