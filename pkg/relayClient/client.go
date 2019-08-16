package relayClient

import "net"

type relayClient struct {
	network string
	address string
	port string
	Connection net.Conn
	hostname string
	messageSize int
}

type RelayClient struct{
	relayClient
}
