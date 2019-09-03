package relayServer

import (
	"fmt"
	"github.com/wvanlit/go-relay/pkg/relay"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

// Connection Response codes
type connectionResponse int

const (
	OK    connectionResponse = 0
	CLOSE connectionResponse = 1
	ERROR connectionResponse = 2
)

func (server *relayServer) createClient(id string, messageSize int, address string, port string) (*Client, error) {
	// Check for existing clients with the same ID
	server.clientLock.Lock()
	for _, client := range server.clients {
		if client.id == id {
			return &Client{}, fmt.Errorf("name already used")
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
		busy:        false,
		close:       false,
	}

	// Add client to list of clients
	server.clients = append(server.clients, &client)
	server.clientLock.Unlock()
	return &client, nil
}

func (server *relayServer) FindClient(id string) (*Client, error) {
	server.clientLock.Lock()
	defer server.clientLock.Unlock()
	for _, client := range server.clients {
		if client.id == id {
			return client, nil
		}
	}

	return &Client{}, fmt.Errorf(fmt.Sprintf("ID: %q not found", id))
}

func (server *relayServer) HandleConnection(conn net.Conn) {
	defer func() {
		_ = conn.Close()
	}()

	// Identify Connection
	fmt.Println("New Address Connected:", conn.RemoteAddr().String())

	// Create New Client
	data := make([]byte, 100)
	ReadClientMessage(conn, &data)

	// Parse Data
	identity := strings.Split(string(data), "|")
	if len(identity) <= 1 {
		return
	}
	messageSize, parseError := strconv.ParseInt(strings.TrimRight(identity[1], "\x00"), 10, 64)
	if parseError != nil {
		log.Println(parseError)
		return
	}
	// Get Address Data
	address := strings.Split(conn.RemoteAddr().String(), ":")

	// Create New Client
	client, creationError := server.createClient(identity[0], int(messageSize), address[0], address[1])

	// Report Outcome to Client
	if creationError != nil {
		log.Println(creationError)
		MessageClient(conn, []byte(creationError.Error()))
		_ = conn.Close()
		return

	} else {
		MessageClient(conn, []byte("OK"))
		client.connection = conn

		// Defer shutdown procedure
		defer func() {
			_ = conn.Close()
			server.clientLock.Lock()
			err := server.removeClient(client.id)
			if err != nil {
				fmt.Println(err.Error())
			}
			server.clientLock.Unlock()
		}()

	}

	// Create Channels
	connectionData := make(chan string)
	responses := make(chan connectionResponse)

	// Retrieve Data
	go retrieveData(client, connectionData)
	defer func() {
		select {
		case client.requests <- relay.CLOSE_CONNECTION:
			return
		case <-time.After(time.Second * 5):
			return
		}
	}()

	// Handle Data
	go func() {
		for !client.close {
			select {
			case data := <-connectionData:
				response := server.HandleMessage(client, data)
				responses <- response
				if response == CLOSE {
					client.close = true
				}
			case <-time.After(time.Second):
				if client.close {
					select {
					case responses <- CLOSE:
						continue
					case <-time.After(time.Second):
						continue
					}
				}
			}
		}
	}()

	// Handle Responses
	for !client.close {
		select {
		case response := <-responses:
			switch response {
			case CLOSE:
				client.close = true
				break
			}
		case <-time.After(time.Second):
			continue
		}

	}
}

func retrieveData(client *Client, output chan string) {
	conn := client.connection
	for {
		// Start retrieving data
		select {
		case request := <-client.requests:
			// Handle different requests
			switch request {
			case "pipe":
				client.results <- client.connection
				client.busy = true
			case "pipe close":
				client.busy = false
			case relay.CLOSE_CONNECTION:
				client.close = true
				return
			default:
				fmt.Println("Unknown Request:", request)
			}



		default:
			if client.busy {
				continue
			}

			data := make([]byte, client.messageSize)
			// Read data from connection
			_ = conn.SetReadDeadline(time.Now().Add(time.Second))
			n, readError := conn.Read(data)

			// Handle potential errors
			if readError != nil {
				// Ignore timeout
				if e, ok := readError.(net.Error); ok && e.Timeout() {
					continue
				} else if readError != io.EOF {
					log.Println("Retrieve:", readError.Error())
					return
				}
			}

			// Handle the message if data is received
			if n > 0 {
				select {
				case output <- string(data[:n]):
					continue
				case <-time.After(time.Second):
					continue
				}
			}

		}

	}
}

func (server *relayServer) HandleMessage(client *Client, message string) connectionResponse {
	switch {
	// Pipe Request
	case strings.Contains(message, relay.START_PIPE):
		// Find Client
		id := strings.Split(message, ":")[1]
		pipeClient, err := server.FindClient(id)
		if err != nil {
			log.Println(err)
			return ERROR
		}
		// Send Request
		pipeClient.requests <- "pipe"

		// Get Request Answer
		pipeConnection := <-pipeClient.results

		// Create Pipe
		fmt.Println("Starting Pipe between", pipeClient.id, "(", pipeConnection.RemoteAddr(), ") and", client.id, "(", client.connection.RemoteAddr(), ")")

		// Set Client to Busy
		client.busy = true
		Pipe(client.connection, pipeConnection, client, pipeClient)

		// Close Pipe
		client.busy = false
		pipeClient.requests <- "pipe close"

		fmt.Println("Pipe between", pipeClient.id, "(", pipeConnection.RemoteAddr(), ") and", client.id, "(", client.connection.RemoteAddr(), ") closed")
		return OK

	// Close Connection
	case message == relay.CLOSE_CONNECTION:
		fmt.Println("Closing Connection:", client.id, "(", client.connection.RemoteAddr(), ")")
		return CLOSE

	// Display Message
	default:
		fmt.Println("Message:", message)
		return OK
	}
}
