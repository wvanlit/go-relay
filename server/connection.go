package server

import (
	"bufio"
	"fmt"
	"github.com/wvanlit/go-relay/global"
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
	username   string
}

func (c *RelayConnection) closeConnection() {
	c.serverPipe <- CreateDeleteUserMessage(c.username)
	c.Open = false
	fmt.Printf("Closing connection to user %s.\n", c.username)
}

func (c *RelayConnection) SendMessage(message string) {
	_, err := fmt.Fprintln(c.Connection, message)
	if err != nil {
		// Close on network error
		_, ok := err.(*net.OpError)
		if ok {
			c.closeConnection()
		} else {
			fmt.Println("Error on sending message:", err)
		}
	}
}

func (c *RelayConnection) ReceiveMessage() string {
	message, err := c.reader.ReadString('\n')
	if err != nil && err != io.EOF {
		// Close on network error
		_, ok := err.(*net.OpError)
		if ok {
			c.closeConnection()
		} else {
			fmt.Printf("Error on Receive Message: %s\n", err)
		}
	}
	// Make sure the message is not only a '\n'
	if message == "\n" {
		return c.ReceiveMessage()
	}

	return message
}

func (c *RelayConnection) HandleConnection() {
	c.Identify()
	c.ReadMessages()
}

func (c *RelayConnection) Identify() {
	feedback := make(chan bool, 0)

	c.SendMessage("Please Identify Yourself.")
	for c.Open {
		input := c.ReceiveMessage()

		// Restart if not an identification
		if !strings.Contains(input, string(global.IDENTIFY)) {
			c.SendMessage(fmt.Sprintf("Please use %s:USERNAME to identify yourself.", string(global.IDENTIFY)))
			continue
		}
		// Get Username
		identity := global.GetIdentification(input)
		fmt.Printf("Server Received Identification %s\n", identity)
		// Check for existing username
		c.serverPipe <- CreateCheckUserMessage(identity, feedback)
		positiveResponse := <-feedback

		if positiveResponse {
			c.SendMessage(string(global.IDENTITY_EXISTS))
			continue
		}

		// Register User
		c.serverPipe <- CreateCreateUserMessage(identity, c)
		c.username = identity
		c.SendMessage(string(global.IDENTITY_OK))
		break
	}
}

func (c *RelayConnection) ReadMessages() {
	for c.Open {
		netData := c.ReceiveMessage()
		if strings.TrimSpace(netData) == "STOP" {
			c.serverPipe <- CreateCloseServerMessage()
			return
		}
		fmt.Printf("%s> %s", c.username, netData)
		t := time.Now()
		myTime := t.Format(time.RFC3339)

		c.SendMessage(myTime)
	}
}

func (c *RelayConnection) PipeMessages() {

}
