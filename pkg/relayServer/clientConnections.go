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
	results     chan net.Conn
	connection  net.Conn
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
		requests:    make(chan string),
		results:     make(chan net.Conn),
	}

	// Add client to list of clients
	server.clients[id] = client
	return client, nil
}

func (server relayServer) FindClient(id string) (Client, error) {
	for key, client := range server.clients {
		if key == id {
			return client, nil
		}
	}

	return Client{}, fmt.Errorf(fmt.Sprintf("ID: %q not found", id))
}

func (server relayServer) HandleConnection(conn net.Conn) {
	defer conn.Close()

	// Identify Connection
	fmt.Println("New Address Connected:",conn.RemoteAddr().String())

	// Create New Client
	data := make([]byte, 100)
	ReadClientMessage(conn, &data)

	// Parse Data
	identity := strings.Split(string(data), "|")
	messageSize, err := strconv.ParseInt(strings.TrimRight(identity[1], "\x00"), 10, 64)
	if err != nil {
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
			fmt.Println("Closing", client.id)
			delete(server.clients, client.id)
		}()
		client.connection = conn
	}
	// Read Data
	for {
		select {
		case request := <-client.requests:
			fmt.Println("Gotten Request:",request)
			client.results <- client.connection
		default:
			data := make([]byte, client.messageSize)
			// Read data from connection
			n, err := conn.Read(data)
			// If data is received, print it
			if err != nil && err != io.EOF {
				log.Println(err.Error())

				return
			}
			if n > 0 {
				server.HandleMessage(client, string(data[:n]))
			}

		}

	}
}

func (server relayServer) HandleMessage(client Client, message string) {
	switch {
	// Connection Request
	case strings.Contains(message, "conn:"):
		// Find Client
		id := strings.Split(message, ":")[1]
		pipeClient, err := server.FindClient(id)
		if err != nil {
			log.Println(err)
			return
		}
		// Send Request
		pipeClient.requests <- "pipe"
		// Get Request Answer
		pipeConnection := <-pipeClient.results
		// Create Pipe
		fmt.Println("Starting Pipe between",pipeClient.id,pipeConnection.RemoteAddr(),"and",client.id,client.connection.RemoteAddr())
		Pipe(client.connection, pipeConnection)

	default:
		fmt.Println("Message:", message)
	}
}
