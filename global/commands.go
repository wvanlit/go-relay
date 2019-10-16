package global

import "strings"

type Command string

const (
	STOP_CONNECTION Command = "STOP_CONNECTION"
	STOP_PIPE       Command = "STOP_PIPE"
	START_PIPE      Command = "START_PIPE"
	IDENTIFY        Command = "IDENTITY"
	IDENTITY_EXISTS Command = "IDENTITY_EXISTS"
	IDENTITY_OK     Command = "IDENTITY_OK"
)

var Separator string = ":"

func CreateIdentification(name string) string {
	return string(IDENTIFY) + Separator + name
}

func GetIdentification(command string) string {
	name := strings.Split(command, Separator)[1]
	// Remove \n
	return name[:len(name)-1]
}

func CreatePipeCommand(name string) string {
	return string(START_PIPE) + Separator + name
}

func GetNameFromPipeCommand(command string) string {
	return strings.Split(command, Separator)[1]
}
