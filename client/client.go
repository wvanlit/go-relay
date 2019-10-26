package client

import (
	"bufio"
	"fmt"
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
		Open:       true,
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

func (r *RelayClient) RunWorkers() {
	go r.SendingWorker()
	go r.ReceivingWorker()
}

func (r *RelayClient) SendMessage(message string) {
	r.output <- message
}

func (r *RelayClient) SendingWorker() {
	for r.Open {
		message := <-r.output
		_, err := fmt.Fprint(r.Connection, message)
		if err != nil {
			fmt.Println("Error on sending message:", err)
		}
	}
}

func (r *RelayClient) ReceiveMessage() string {
	return <-r.input
}

func (r *RelayClient) ReceivingWorker() {
	for r.Open {
		select {

		default:
			message, err := r.reader.ReadString('\n')
			if err != nil && err != io.EOF {
				fmt.Printf("Error on Receive Message: %s\n", err)
			}
			r.input <- message
		}

	}
}

func (r *RelayClient) RunClient() {
	defer r.Connection.Close()
	if !r.Open {
		fmt.Println("Connection not open!")
		return
	}

	go r.RunWorkers()
	time.Sleep(time.Millisecond * 100)
	fmt.Print(r.ReceiveMessage())
	go func() {
		for {
			r.handleUserInput()
			time.Sleep(time.Millisecond*500)
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
