package ws

import (
	"context"
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

type channelMessage struct {
	data []byte
	mt   int
	conn *websocket.Conn
}

// Hub handles all websocket connections
type Hub struct {
	conns map[*websocket.Conn]struct{}
	ch    chan channelMessage
	mux   sync.RWMutex
}



// NewHub returns new hub instance
func NewHub(ctx context.Context) *Hub {
	h := &Hub{
		conns: make(map[*websocket.Conn]struct{}),
		ch:    make(chan channelMessage),
	}
	go h.sendLoop()
	go func() {
		<-ctx.Done()
		close(h.ch)
	}()
	return h
}

// Add conenction to map
func (h *Hub) Add(conn *websocket.Conn) {
	h.mux.Lock()
	defer h.mux.Unlock()

	h.conns[conn] = struct{}{}
}

// Remove connection from a map
func (h *Hub) Remove(conn *websocket.Conn) {
	h.mux.Lock()
	defer h.mux.Unlock()

	conn.Close()
	delete(h.conns, conn)
}

// Send sends message to a connection
func (h *Hub) Send(conn *websocket.Conn, data []byte, mt int) {
	h.ch <- channelMessage{
		data: data,
		conn: conn,
		mt:   mt,
	}
}

// SendToAll sends message to all clients
func (h *Hub) SendToAll(data []byte) {
	h.mux.Lock()
	defer h.mux.Unlock()
	for conn := range h.conns {
		h.Send(conn, data, websocket.TextMessage)
	}
}

func (h *Hub) sendLoop() {
	for env := range h.ch {
		if env.mt == websocket.TextMessage {
			err := env.conn.WriteMessage(
				env.mt,
				env.data)
			if err != nil {
				log.Printf("failed to write to a websocket: %v", err)
				continue
			}
		} else {
			err := env.conn.WriteControl(env.mt, env.data, time.Now().Add(time.Second*5))
			if err != nil {
				log.Printf("failed to write control message to a websocket: %v", err)
				continue
			}
		}
	}
}

// HandleBroadcast handles /broadcast http endpoint
// and send a message to every websocket connection
func (h *Hub) HandleBroadcast(w http.ResponseWriter, r *http.Request) {
	h.SendToAll([]byte("Hi ALL!"))
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
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("failed to read from websocket: %v", err)
			return
		}

		if mt == websocket.PingMessage {
			h.Send(conn, []byte{}, websocket.PongMessage)
		}

		if mt == websocket.TextMessage {
			log.Printf("got text message: '%s'", string(msg))

			if string(msg) == "/test" {
				h.Send(conn, []byte("Hello, Websocket!"), websocket.TextMessage)
			}
		}
	}
}
