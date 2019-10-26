package global

type State int

const (
	Offline State = 0
	Open    State = 1
	InPipe  State = 2
)
