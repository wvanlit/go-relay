package relay


type Command string

const (
	// Closing connections
	CLOSE_CONNECTION = "CLOSE"

	// Creating and Closing pipes
	START_PIPE = "PIPE"
	STOP_PIPE = "EXIT PIPE"
)