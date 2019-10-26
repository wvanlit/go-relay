package server

import (
	"fmt"
	"github.com/wvanlit/go-relay/global"
)

type RequestPipeMessage struct {
	messageType MessageType
	request     PipeRequest
	feedback    chan PipeResponse
}

type PipeRequest struct {
	target string
	output chan string
}

type PipeResponse struct {
	ok     bool
	output chan string
}

func CreateRequestPipeMessage(req PipeRequest, feedback chan PipeResponse) RequestPipeMessage {
	return RequestPipeMessage{
		messageType: REQUEST_PIPE,
		request:     req,
		feedback:    feedback,
	}
}

func (rpm RequestPipeMessage) GetType() MessageType {
	return rpm.messageType
}

func (rpm RequestPipeMessage) ProcessMessage(s *RelayServer) {
	// Check for user
	user, ok := s.GetUser(rpm.request.target)

	if !ok {
		fmt.Printf("User %s doesn't exist\n", rpm.request.target)
		rpm.feedback <- PipeResponse{ok: false,}
		return
	}

	// Check user availability
	if user.State != global.Open {
		fmt.Printf("User %s not Open\n", user.username)
		rpm.feedback <- PipeResponse{ok: false,}
		return
	}

	// Send users pipe command
	user.RequestPipe <- rpm.request
	rpm.feedback <- PipeResponse{ok: true, output: user.output}
}
