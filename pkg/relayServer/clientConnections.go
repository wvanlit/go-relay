package relayServer

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

type Client struct {
	id          string
	messageSize int
	address     string
	port        string
	requests    chan string
}

func (server relayServer) createClient(id string, messageSize int, address string, port string) (Client, error) {
	// Check for existing clients with the same ID
	for key, _ := range server.clients {
		if key == id {
			return Client{}, fmt.Errorf("name already used")
		}
	}
	// Create new Client
	client := Client{
		id:          id,
		messageSize: messageSize,
		address:     address,
		port:        port,
	}

	// Add client to list of clients
	server.clients[id] = client
	return client, nil
}

func (server relayServer) HandleConnection(conn net.Conn) {
	defer conn.Close()

	// Identify Connection
	println(conn.RemoteAddr().String())

	// Create New Client
	data := make([]byte, 100)
	ReadClientMessage(conn, &data)

	// Parse Data
	identity := strings.Split(string(data), "|")
	fmt.Println("Input",identity)
	messageSize, err := strconv.ParseInt(strings.TrimRight(identity[1],"\x00"), 10, 64)
	if err != nil{
		log.Println(err)
		return
	}
	// Get Address Data
	address := strings.Split(conn.RemoteAddr().String(), ":")

	// Create New Client
	client, err := server.createClient(identity[0], int(messageSize), address[0], address[1])

	// Report Outcome to Client
	if err != nil {
		log.Println(err)
		MessageClient(conn, []byte(err.Error()))
		conn.Close()
		return
	} else {
		MessageClient(conn, []byte("OK"))
		defer func() {
			fmt.Println("Closing",client.id)
			delete(server.clients, client.id)
		}()
	}
	msize := int(client.messageSize)
	println(msize)
	// Read Data
	for {
		select {
		case request := <- client.requests:
			fmt.Println(request)
		default:
			data := make([]byte, msize)
			// Read data from connection
			n, err := conn.Read(data)
			// If data is received, print it
			if err != nil && err != io.EOF{
				log.Println(err.Error())

				return
			}
			if n > 0 {
				fmt.Println("Message :", string(data))
			}

		}

	}
}
