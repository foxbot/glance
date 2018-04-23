package server

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type socketServer struct {
	clients     []*websocket.Conn
	upgrader    websocket.Upgrader
	totalShards int
	shardCache  map[int]int
	updateList  chan ShardUpdate
}

func newSocketServer() *socketServer {
	totalShards, err := strconv.ParseInt(os.Getenv("BOT_SHARDS"), 10, 32)
	if err != nil {
		log.Panicln(err)
	}

	return &socketServer{
		clients:     make([]*websocket.Conn, 0),
		upgrader:    websocket.Upgrader{},
		totalShards: int(totalShards),
		shardCache:  make(map[int]int, totalShards),
		updateList:  make(chan ShardUpdate, 128),
	}
}

func (s *socketServer) socketUpgrader(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
	s.sayHello(conn)
	s.clients = append(s.clients, conn)

	closeHandler := func(code int, reason string) error {
		log.Println(conn.RemoteAddr(), "closed", code, "with", reason)
		s.die(conn)
		return nil
	}
	conn.SetCloseHandler(closeHandler)
}

func (s *socketServer) run() {
	demo := os.Getenv("ENABLE_DEMO")
	if demo != "" {
		demo, err := strconv.ParseBool(demo)
		if err != nil {
			log.Panicln(err)
		}
		if demo {
			log.Println("DEMO MODE enabled!")
			go s.testData()
		}
	}

	for {
		select {
		case u := <-s.updateList:
			s.sayUpdate(u)
			s.shardCache[u.ID] = u.Status
		}
	}
}

func (s *socketServer) sayHello(conn *websocket.Conn) {
	data := Message{
		Op: OpHello,
		Data: HelloMessage{
			TotalShards: s.totalShards,
			State:       s.shardCache,
		},
	}
	err := conn.WriteJSON(data)
	if err != nil {
		log.Println(err)
		s.die(conn)
	}
	log.Println(conn.RemoteAddr(), "opened and saluted")
}

func (s *socketServer) sayUpdate(u ShardUpdate) {
	data := Message{
		Op: OpUpdate,
		Data: UpdateMessage{
			Shard:  u.ID,
			Status: u.Status,
		},
	}

	b, err := json.Marshal(data)
	if err != nil {
		log.Panicln(err)
	}

	pm, err := websocket.NewPreparedMessage(websocket.TextMessage, b)
	if err != nil {
		log.Panicln(err)
	}

	for _, c := range s.clients {
		err = c.WritePreparedMessage(pm)
		if err != nil {
			log.Println(err)
			s.die(c)
		}
	}
}

func (s *socketServer) die(conn *websocket.Conn) {
	conn.Close()
	// find index
	i := -1
	for idx, c := range s.clients {
		if c == conn {
			i = idx
			break
		}
	}
	if i == -1 {
		log.Println("failed to find index of socket")
		return
	}

	// reassign s.clients for brevity
	c := s.clients
	// swap conn with the back of the array
	// [0 1 2 3 4] removing 1 turns into [0 4 2 3 1]
	c[len(c)-1], c[i] = c[i], c[len(c)-1]
	// set clients to the array - 1
	// [0 4 2 3]
	s.clients = c[:len(c)-1]
}

// TODO: remove this for prod
func (s *socketServer) testData() {
	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src)

	// populate initial list with data
	for i := 0; i < s.totalShards; i++ {
		// [1, 5) - don't want to send unknowns
		status := rng.Intn(4) + 1
		u := ShardUpdate{
			ID:     i,
			Status: status,
		}
		s.updateList <- u
	}

	for {
		shard := rng.Intn(s.totalShards)
		// [1, 5) - don't want to send unknowns
		status := rng.Intn(4) + 1

		u := ShardUpdate{
			ID:     shard,
			Status: status,
		}
		s.updateList <- u

		time.Sleep(time.Second * 1)
	}
}
