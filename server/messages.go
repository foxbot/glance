package server

const (
	// OpHello is the opcode for HelloMessage
	OpHello int = iota
	// OpUpdate is the opcode for UpdateMessage
	OpUpdate
	// StatusUnknown is a shard with unknown status
	StatusUnknown int = iota
	// StatusStarting is starting
	StatusStarting
	// StatusOnline is online
	StatusOnline
	// StatusStopping is stopping
	StatusStopping
	// StatusOffline is dead
	StatusOffline
)

type (
	shardUpdate struct {
		ID     int
		Status int
	}
	// Message is the base outgoing message to any client
	Message struct {
		Op   int
		Data interface{}
	}
	// HelloMessage introduces initial state to a client
	HelloMessage struct {
		TotalShards int
		State       map[int]int
	}
	// UpdateMessage is sent when a shard changes status
	UpdateMessage struct {
		Shard  int
		Status int
	}
)
