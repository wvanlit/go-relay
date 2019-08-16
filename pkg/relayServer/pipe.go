package relayServer

import (
	"fmt"
	"io"
	"log"
	"net"
)

func createChannelForConnection(conn net.Conn) chan []byte {
	channel := make(chan []byte)

	// Read Connection and Dump in channel
	go readAndDumpIntoChannel(conn, channel)

	return channel
}

func readAndDumpIntoChannel(conn net.Conn, channel chan []byte) {
	data := make([]byte, 1024)
	for {
		n, err := conn.Read(data)
		// Handle Errors
		if err != nil && err != io.EOF {
			log.Printf(err.Error())
			return
		}

		// Dump message
		if n > 0 {
			message := make([]byte, n)
			copy(message, data)
			channel <- message
		}

	}
}

func Pipe(connection1 net.Conn, connection2 net.Conn) {
	channel1 := createChannelForConnection(connection1)
	channel2 := createChannelForConnection(connection2)

	for {
		select {
		case messageTo1 := <-channel2:
			if messageTo1 != nil {
				_, err := connection1.Write(messageTo1)
				if err != nil {
					log.Println(err.Error())
				}
			}
		case messageTo2 := <-channel1:
			if messageTo2 != nil {
				_, err := connection2.Write(messageTo2)
				if err != nil {
					log.Println(err.Error())
				}
			}
		}
	}
}
