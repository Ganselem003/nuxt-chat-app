package ws

import (
	"encoding/json"
	"sync"
)

type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	direct     chan DirectMessage
	mu         sync.Mutex
}

type DirectMessage struct {
	Raw []byte
	To  uint
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		direct:     make(chan DirectMessage),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			h.broadcastPresence()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			h.broadcastPresence()

		case dm := <-h.direct:
			h.mu.Lock()
			for c := range h.clients {
				if c.ID == dm.To {
					c.send <- dm.Raw
				}
			}
			h.mu.Unlock()
		}
	}
}

// connect hiisen buh usert username listiig ilgeene
func (h *Hub) broadcastPresence() {
	h.mu.Lock()
	defer h.mu.Unlock()
	users := []string{}
	for c := range h.clients {
		users = append(users, c.Username)
	}
	msg := map[string]interface{}{
		"type":  "presence",
		"users": users,
	}
	b, _ := json.Marshal(msg)
	for c := range h.clients {
		c.send <- b
	}
}
