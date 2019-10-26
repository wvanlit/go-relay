package examples

import (
	"fmt"
	relayClient "github.com/wvanlit/go-relay/src/client"
	"github.com/wvanlit/go-relay/src/global"
	"os"
	"time"
)

func main() {
	arguments := os.Args
	if len(arguments) <= 3 {
		fmt.Println("Please provide a host, port number and username")
		return
	}
	client := relayClient.CreateRelayClient(arguments[1], arguments[2])
	err := client.StartClient(arguments[3])
	if err != nil{
		fmt.Println(err.Error())
	}
	client.SendMessage("Hello Server!")
	client.SendMessage("Bye Server!")
	client.SendMessage(global.CreatePipeCommand("B"))
	client.SendMessage("Hi B!")
	time.Sleep(time.Millisecond*500)
	client.SendMessage("Bye B!")
	client.SendMessage(string(global.STOP_PIPE))
	time.Sleep(time.Millisecond*500)
	client.StopClient()
}
