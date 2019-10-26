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
	Connection    net.Conn
	State         global.State
	reader        *bufio.Reader
	writer        *bufio.Writer
	input         chan string
	output        chan string
	pipeOutput    chan string
	serverChannel chan Message
	RequestPipe   chan PipeRequest
	username      string
}

func (c *RelayConnection) closeConnection() {
	c.serverChannel <- CreateDeleteUserMessage(c.username)
	c.State = global.Offline
	fmt.Printf("Closing connection to user %s.\n", c.username)
}

func (c *RelayConnection) SendMessage(message string) {
	c.output <- message
}

func (c *RelayConnection) sendingWorker() {
	for c.State != global.Offline {
		select {
		// Check if there is output
		case message := <-c.output:
			if strings.Contains(message, string(global.STOP_PIPE)) {
				c.State = global.Open
				c.pipeOutput = nil
			}
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
		// Interrupt
		case <-time.After(time.Minute):
			continue
		}

	}

}

func (c *RelayConnection) ReceiveMessage() string {
	return <-c.input
}

func (c *RelayConnection) receivingWorker() {
	for c.State != global.Offline {
		message, err := c.reader.ReadString('\n')
		if err != nil && err != io.EOF {
			// Close on network error
			_, ok := err.(*net.OpError)
			if ok {
				fmt.Println(err.(*net.OpError))
				c.closeConnection()
			} else {
				fmt.Printf("Error on Receive Message: %s\n", err)
			}
		}
		// Make sure the message is not only a '\n'
		if message == "\n" || len(message) == 0 {
			continue
		}

		if c.State == global.InPipe {
			// Remove '\n'
			c.pipeOutput <- message[:len(message)-1]
			if strings.Contains(message, string(global.STOP_PIPE)) {
				c.State = global.Open
				c.pipeOutput = nil
			}
		} else {
			c.input <- message
		}
	}
}

func (c *RelayConnection) startWorkers() {
	go c.receivingWorker()
	go c.sendingWorker()
}

func (c *RelayConnection) HandleConnection() {
	c.startWorkers()
	c.identify()
	c.readMessages()
}

func (c *RelayConnection) identify() {
	feedback := make(chan bool, 0)

	c.SendMessage("Please identify Yourself.")
	for c.State == global.Open {
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
		c.serverChannel <- CreateCheckUserMessage(identity, feedback)
		positiveResponse := <-feedback

		if positiveResponse {
			c.SendMessage(string(global.IDENTITY_EXISTS))
			continue
		}

		// Register User
		c.serverChannel <- CreateCreateUserMessage(identity, c)
		c.username = identity
		c.SendMessage(string(global.IDENTITY_OK))
		break
	}
}

func (c *RelayConnection) readMessages() {
	for c.State != global.Offline {
		select {
		case message := <-c.input:
			c.processMessage(message)
		case req := <-c.RequestPipe:
			c.State = global.InPipe
			c.pipeOutput = req.output
			//for c.State == global.InPipe {
			//	time.Sleep(time.Second)
			//}
		case <-time.After(time.Minute):
			continue
		}
	}
}

func (c *RelayConnection) processMessage(message string) {
	fmt.Printf("%s> %s", c.username, message)

	// Handle Pipe Request
	if strings.Contains(message, string(global.START_PIPE)) {
		target := PipeRequest{
			target: global.GetNameFromPipeCommand(message),
			output: c.output,
		}
		feedback := make(chan PipeResponse, 0)
		c.serverChannel <- CreateRequestPipeMessage(target, feedback)
		resp := <-feedback
		if resp.ok {
			c.pipeOutput = resp.output
			c.State = global.InPipe
			c.SendMessage("PIPE OK")
			return
		} else {
			c.SendMessage("Not able to create channel.")
			return
		}
	}

	// Return server time
	t := time.Now()
	myTime := t.Format(time.RFC3339)

	c.SendMessage(myTime)
}

func (c *RelayConnection) startPipe(target string) {

}
