# go-relay
A simple message relay written in Go using TCP connections.
It can handle multiple concurrent connections that communicate with either the server or to another client.

All connections between clients are one-to-one.

### Starting a Connection

The following protocol is used to start a connection:
1. Start connection with the server
1. Identify yourself via `IDENTITY:USERNAME`
1. You'll receive a response with the availability of the username
1. If your username is available then you've succesfully started a connection

### Starting a Pipe
To start a pipe send the following message `START_PIPE:TARGET` where target is the username of the client you want to connect with.