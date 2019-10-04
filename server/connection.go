package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

type RelayConnection struct {
	Connection net.Conn
	Open       bool
	reader     *bufio.Reader
	writer     *bufio.Writer
}

func (c *RelayConnection) SendMessage(message string) {
	_, _ = fmt.Fprintf(c.Connection, message)
}

func (c *RelayConnection) ReceiveMessage() string {
	message, err := c.reader.ReadString('\n')
	if err != nil && err != io.EOF {
		fmt.Printf("Error on Receive Message: %s\n", err)
	}
	return message
}
