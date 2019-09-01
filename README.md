# go-relay
A Relay Server and Client made in Go that can handle up to 750+ concurrent connections. It's recommended that you keep the connections below 200, especially if you're using a lot of pipes as the overhead slows down the server significantly.

### Protocol
To start a connection the following protocol should be followed:
1. A client is created : `client := relayClient.CreateTCPClient(clientName, hostname, port, messageSize)`
2. The client is started : `client.Start()`
3. The client identifies itself to the server : `client.Identify()`
4. The client is now free to send commands or messages to the server

To stop a connection use `client.Stop()` or manually send `relay.CLOSE_CONNECTION` to the server.

### Commands
When a client is connected to the server, it can send commands to perform actions.

To start a pipe with another connection either use `client.StartPipe(target)` or directly send `relay.START_PIPE+':'+target` to the server.

To stop a pipe use `client.StopPipe()` or directly send `relay.STOP_PIPE` to the server.
