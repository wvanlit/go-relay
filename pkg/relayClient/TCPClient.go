package relayClient

func CreateTCPClient(hostname string, address string, port string, messageSize int) RelayClient {
	return RelayClient{
		relayClient{
			network:     "tcp",
			address:     address,
			port:        port,
			hostname:    hostname,
			messageSize: messageSize,
		},
	}
}
