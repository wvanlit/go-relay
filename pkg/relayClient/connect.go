package relayClient

import (
	"fmt"
	"github.com/wvanlit/go-relay/pkg/relay"
	"log"
	"net"
	"strings"
	"time"
)

func (client *RelayClient) Start() {
	// Setup Connection to Address
	conn, err := net.Dial(client.network, client.address+":"+client.port)
	if err != nil {
		log.Fatal(err)
	}
	client.Connection = conn
}

func (client *RelayClient) Stop() {
	// Setup Connection to Address
	client.SendString(relay.CLOSE_CONNECTION)
	err := client.Connection.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func (client *RelayClient) Send(message []byte) {
	_, err := client.Connection.Write(message)
	if err != nil {
		log.Println(err)
	}
}

func (client *RelayClient) SendString(message string){
	_, err := client.Connection.Write([]byte(message))
	if err != nil {
		log.Println(err)
	}
}

func (client *RelayClient) Read(message []byte) int {
	n, err := client.Connection.Read(message)
	// Check if there is an error that is not a timeout
	if err != nil && !strings.Contains(err.Error(), "timeout") {
		log.Println(err.Error())
	}
	return n
}

func (client *RelayClient) ReadForDuration(message []byte, duration time.Duration) int {
	_ = client.Connection.SetReadDeadline(time.Now().Add(duration))
	n, err := client.Connection.Read(message)
	// Check if there is an error that is not a timeout
	if err != nil {
		if e, ok := err.(net.Error); ok && e.Timeout(){
			return 0
		}
		log.Println(err.Error())
	}
	return n
}

func (client *RelayClient) Identify() bool {
	client.Send([]byte(fmt.Sprint(client.hostname, "|", client.messageSize)))
	data := make([]byte, 100)
	client.Read(data)
	trimmedData := strings.TrimSpace(string(data))
	if !strings.Contains(trimmedData, "OK") {
		fmt.Println(trimmedData)
		return false
	}

	return true
}

func (client *RelayClient) StartPipe(target string){
	client.SendString(relay.START_PIPE+":"+target)
}

func (client *RelayClient) StopPipe(){
	client.SendString(relay.STOP_PIPE)
}