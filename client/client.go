package client

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

type RelayClient struct {
	Address    string
	Connection net.Conn
	Open       bool
	reader     *bufio.Reader
	writer     *bufio.Writer
}

func CreateRelayClient(address string, port string) *RelayClient {
	addr := address + ":" + port
	conn, err := net.Dial("tcp", addr)

	relay := &RelayClient{
		Address:    addr,
		Connection: conn,
		Open:       true,
		reader:     bufio.NewReader(conn),
		writer:     bufio.NewWriter(conn),
	}

	if err != nil {
		fmt.Printf("Error on dial to %s: %s\n", relay.Address, err)
		relay.Open = false
	}
	return relay
}

func (r *RelayClient) SendMessage(message string) {
	_, _ = fmt.Fprintf(r.Connection, message)
}

func (r *RelayClient) ReceiveMessage() string {
	message, err := r.reader.ReadString('\n')
	if err != nil && err != io.EOF {
		fmt.Printf("Error on Receive Message: %s\n", err)
	}
	return message
}

func (r *RelayClient) RunClient() {
	defer r.Connection.Close()
	if !r.Open {
		fmt.Println("Connection not open!")
		return
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')

		r.SendMessage(text)

		message := r.ReceiveMessage()
		fmt.Printf("-> %s", message)
		if strings.TrimSpace(string(text)) == "STOP" {
			fmt.Println("TCP client exiting...")
			return
		}
	}
}
