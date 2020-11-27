package ws

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Hub handles all websocket connections
type Hub struct {
	conns map[*websocket.Conn]*sync.Mutex
	mux   sync.RWMutex
}

// NewHub returns new hub instance
func NewHub() *Hub {
	return &Hub{conns: make(map[*websocket.Conn]*sync.Mutex)}
}

// Add conenction to map
func (h *Hub) Add(conn *websocket.Conn) {
	h.mux.Lock()
	defer h.mux.Unlock()

	h.conns[conn] = &sync.Mutex{}
}

// Remove connection from a map
func (h *Hub) Remove(conn *websocket.Conn) {
	h.mux.Lock()
	defer h.mux.Unlock()

	conn.Close()
	delete(h.conns, conn)
}

func (h *Hub) lock(conn *websocket.Conn) {
	h.mux.Lock()
	defer h.mux.Unlock()

	h.conns[conn].Lock()
}

func (h *Hub) unlock(conn *websocket.Conn) {
	h.mux.Lock()
	defer h.mux.Unlock()

	h.conns[conn].Unlock()
}

// HandleBroadcast handles /broadcast http endpoint
// and sende a message to every websocket connection
func (h *Hub) HandleBroadcast(w http.ResponseWriter, r *http.Request) {
	h.mux.Lock()
	defer h.mux.Unlock()

	log.Println("got /broadcast")

	for conn, mux := range h.conns {
		log.Println("/broadcast loop")
		go func(conn *websocket.Conn, mux *sync.Mutex) {
			//h.lock(conn)
			log.Println("/broadcast sending")
			err := conn.WriteMessage(
				websocket.TextMessage,
				[]byte("Hi ALL!"))
			//h.unlock(conn)

			log.Printf("/broadcast send error: %v", err)

			if err != nil {
				log.Printf("failed to write to a websocket: %v", err)
				return
			}
		}(conn, mux)
	}
}

// HandleWS handles /ws endpoint and upgrades connection protocol
func (h *Hub) HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade connection failed: %v", err)
		return
	}

	h.Add(conn)
	defer h.Remove(conn)

	for {
		h.lock(conn)
		mt, msg, err := conn.ReadMessage()
		h.unlock(conn)
		if err != nil {
			log.Printf("failed to read from websocket: %v", err)
			return
		}

		if mt == websocket.PingMessage {
			h.lock(conn)
			err := conn.WriteControl(
				websocket.PongMessage,
				[]byte{},
				time.Now().Add(time.Second))
			h.unlock(conn)

			if err != nil {
				log.Printf("failed to read from websocket: %v", err)
				return
			}
		}

		if mt == websocket.TextMessage {
			log.Printf("got text message: '%s'", string(msg))

			if string(msg) == "/test" {
				h.lock(conn)
				err := conn.WriteMessage(
					websocket.TextMessage,
					[]byte("Hello, Websocket!"))
				h.unlock(conn)

				if err != nil {
					log.Printf("failed to write to a websocket: %v", err)
					return
				}
			}
		}
	}
}
