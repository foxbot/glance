package server

const (
	// OpHello is the opcode for HelloMessage
	OpHello int = iota
	// OpUpdate is the opcode for UpdateMessage
	OpUpdate
	// OpTick is a keep-alive for the websocket, no payload
	OpTick
	// StatusUnknown is a shard with unknown status
	StatusUnknown int = 0
	// StatusWaiting is waiting for startup
	StatusWaiting = 1
	// StatusStarting is starting
	StatusStarting = 2
	// StatusOnline is online
	StatusOnline = 3
	// StatusStopping is stopping
	StatusStopping = 4
	// StatusOffline is dead
	StatusOffline = 5
)

type (
	// ShardUpdate is the carrier for a shard update
	ShardUpdate struct {
		ID     int
		Status int
		Bot    int
	}
	// Message is the base outgoing message to any client
	Message struct {
		Op   int
		Data interface{}
	}
	// HelloMessage introduces initial state to a client
	HelloMessage struct {
		TotalShards int
		// State: [bot [shard state]]
		State map[int]map[int]int
	}
	// UpdateMessage is sent when a shard changes status
	UpdateMessage struct {
		Shard  int
		Status int
		Bot    int
	}
)
