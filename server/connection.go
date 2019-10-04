package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

type RelayConnection struct {
	Connection net.Conn
	Open       bool
	reader     *bufio.Reader
	writer     *bufio.Writer
	serverPipe chan message
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

func (c *RelayConnection) HandleConnection() {
	for {
		netData := c.ReceiveMessage()
		if strings.TrimSpace(string(netData)) == "STOP" {
			c.serverPipe <- CreateCloseServerMessage()
			break
		}

		fmt.Print("--> ", netData)
		t := time.Now()
		myTime := t.Format(time.RFC3339) + "\n"

		c.SendMessage(myTime)
	}

}
