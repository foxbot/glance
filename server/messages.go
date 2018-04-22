package server

const (
	// OpHello is the opcode for HelloMessage
	OpHello int = iota
)

type (
	// Message is the base outgoing message to any client
	Message struct {
		Op   int
		Data interface{}
	}
	// HelloMessage introduces initial state to a client
	HelloMessage struct {
		TotalShards int
	}
)
