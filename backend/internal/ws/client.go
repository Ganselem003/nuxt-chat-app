package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type WSMessageIn struct {
	To   string `json:"to"`
	Text string `json:"text"`
}

type WSMessageOut struct {
	Type string `json:"type"`
	From string `json:"from"`
	To   string `json:"to"`
	Text string `json:"text"`
}

type Client struct {
	Conn     *websocket.Conn
	send     chan []byte
	hub      *Hub
	Username string
	ID       uint
}

func (c *Client) writePump() {
	ticker := time.NewTicker(time.Second * 54)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Conn.WriteMessage(websocket.TextMessage, msg)
		case <-ticker.C:
			c.Conn.WriteMessage(websocket.PingMessage, []byte{})
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("ws read error: %v", err)
			}
			break
		}

		var in WSMessageIn
		if err := json.Unmarshal(message, &in); err != nil {
			continue
		}

		out := WSMessageOut{
			Type: "message",
			From: c.Username,
			To:   in.To,
			Text: in.Text,
		}
		b, _ := json.Marshal(out)

		// DM
		c.hub.mu.Lock()
		for cl := range c.hub.clients {
			if cl.Username == in.To {
				c.hub.direct <- DirectMessage{Raw: b, To: cl.ID}
				break
			}
		}
		c.hub.mu.Unlock()
	}
}

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

func ServeWS(h *Hub, w http.ResponseWriter, r *http.Request, username string, id uint) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	client := &Client{
		Conn:     conn,
		send:     make(chan []byte, 256),
		hub:      h,
		Username: username,
		ID:       id,
	}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}
