package client

import (
	"bufio"
	"fmt"
	"github.com/wvanlit/go-relay/src/global"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

type RelayClient struct {
	Address    string
	Connection net.Conn
	Open       bool
	reader     *bufio.Reader
	writer     *bufio.Writer
	output     chan string
	input      chan string
}

func CreateRelayClient(address string, port string) *RelayClient {
	addr := address + ":" + port
	conn, err := net.Dial("tcp", addr)

	relay := &RelayClient{
		Address:    addr,
		Connection: conn,
		Open:       false,
		reader:     bufio.NewReader(conn),
		writer:     bufio.NewWriter(conn),
		output:     make(chan string, 0),
		input:      make(chan string, 2),
	}

	if err != nil {
		fmt.Printf("Error on dial to %s: %s\n", relay.Address, err)
		relay.Open = false
	}
	return relay
}

func (r *RelayClient) runWorkers() {
	go r.sendingWorker()
	go r.receivingWorker()
}

func (r *RelayClient) SendMessage(message string) {
	r.output <- message
}

func (r *RelayClient) sendingWorker() {
	for r.Open {
		message := <-r.output
		_, err := fmt.Fprintln(r.Connection, message)
		if err != nil && err != io.EOF {
			// Close on network error
			_, ok := err.(*net.OpError)
			if ok {
				fmt.Println("Stopping Client -> ", err.(*net.OpError))
				r.StopClient()
			} else {
				fmt.Printf("Error on Receive Message: %s\n", err)
			}
		}
	}
}

func (r *RelayClient) ReceiveMessage() string {
	return <-r.input
}

func (r *RelayClient) receivingWorker() {
	for r.Open {
		select {
		default:
			message, err := r.reader.ReadString('\n')
			if err != nil && err != io.EOF {
				// Close on network error
				_, ok := err.(*net.OpError)
				if ok {
					fmt.Println("Stopping Client -> ", err.(*net.OpError))
					r.StopClient()
				} else {
					fmt.Printf("Error on Receive Message: %s\n", err)
				}
			}
			r.input <- message
		}

	}
}

func (r *RelayClient) RunManualClient() {
	defer r.Connection.Close()
	if r.Open {
		fmt.Println("Connection already open!")
		return
	}
	r.Open = true
	go r.runWorkers()
	time.Sleep(time.Millisecond * 100)
	fmt.Print(r.ReceiveMessage())
	go func() {
		for {
			r.handleUserInput()
			time.Sleep(time.Millisecond * 500)
		}
	}()

	for {
		message := r.ReceiveMessage()
		fmt.Printf("-> %s", message)
	}
}

func (r *RelayClient) handleUserInput() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(">> ")
	text, _ := reader.ReadString('\n')
	r.SendMessage(text)
	if strings.TrimSpace(string(text)) == "STOP" {
		fmt.Println("TCP client exiting...")
		r.Open = false
		os.Exit(0)
	}
}

func (r *RelayClient) StartClient(username string) error {
	if r.Open {
		return fmt.Errorf("connection already open")
	}
	r.Open = true

	go r.runWorkers()
	// Identify
	_ = r.ReceiveMessage()
	r.SendMessage(global.CreateIdentification(username) + "\n")
	message := r.ReceiveMessage()
	if !strings.Contains(message, string(global.IDENTITY_OK)) {
		return fmt.Errorf("username '%s' not allowed -> %s", username, message)
	}

	return nil
}

func (r *RelayClient) StopClient() {
	r.SendMessage(string(global.STOP_CONNECTION) + "\n")
	r.Open = false
	time.Sleep(time.Second)
	r.Connection.Close()
}
