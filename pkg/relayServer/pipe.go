package relayServer

import (
	"fmt"
	"github.com/wvanlit/go-relay/pkg/relay"
	"io"
	"log"
	"net"
	"reflect"
	"time"
)

func createChannelForConnection(conn net.Conn, isOn *bool) chan []byte {
	channel := make(chan []byte)
	// Read Connection and Dump in channel
	go readAndDumpIntoChannel(conn, channel, isOn)
	return channel
}

func readAndDumpIntoChannel(conn net.Conn, channel chan []byte, isOn *bool) {
	data := make([]byte, 1024)

	for *isOn {
		_ = conn.SetReadDeadline(time.Now().Add(time.Second))
		n, err := conn.Read(data)
		// Handle Errors
		if err != nil {
			if err == io.EOF {
				continue
			} else if reflect.TypeOf(err) == reflect.TypeOf(&net.OpError{}) {
				// is Timeout
				if e, ok := err.(net.Error); ok && e.Timeout() {
					continue
				}

				// Is not a reset of the connection or timeout -> This is an actual error
				if *isOn {
					log.Printf(err.Error())
				}

			} else {
				log.Printf(err.Error())
				log.Printf(reflect.TypeOf(err).String())
				return
			}
		}

		// Dump message
		if n > 0 {
			message := make([]byte, n)
			copy(message, data)
			select {
			case channel <- message:
				continue
			case <-time.After(time.Second):
				continue
			}

		}

	}

}

// Create a Pipe between 2 Connections, sending data from one directly to the other.
func Pipe(connection1 net.Conn, connection2 net.Conn, client1 *Client, client2 *Client) {
	isOn := true

	channel1 := createChannelForConnection(connection1, &isOn)
	channel2 := createChannelForConnection(connection2, &isOn)

	defer func() {
		fmt.Println("Stopping pipe on", connection1.RemoteAddr(), "and", connection2.RemoteAddr())
		isOn = false
		time.Sleep(time.Second)
	}()

	for client1.close || client2.close {
		select {
		case messageTo1 := <-channel2:
			//fmt.Println("1:", string(messageTo1))
			if messageTo1 != nil {
				if string(messageTo1) == relay.STOP_PIPE {
					return
				}
				_, err := connection1.Write(messageTo1)
				if err != nil {
					log.Println("PIPE:", err.Error())
				}

			}
		case messageTo2 := <-channel1:
			//fmt.Println("2:", string(messageTo2))
			if messageTo2 != nil {
				if string(messageTo2) == relay.STOP_PIPE {
					return
				}
				_, err := connection2.Write(messageTo2)
				if err != nil {
					log.Println("PIPE:", err.Error())
				}
			}
		case <-time.After(time.Second * 5):
			continue
		}
	}
}
