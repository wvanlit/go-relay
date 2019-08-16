package relayClient

func CreateTCPClient(address string, port string) RelayClient {
	return RelayClient{
		relayClient{
			network: "tcp",
			address: address,
			port:    port,
		},
	}
}
