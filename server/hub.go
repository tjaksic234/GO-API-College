package server

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn     *websocket.Conn
	Username string
	RoomName string
	UserID   uint
	RoomID   uint
}

type Hub struct {
	rooms map[string]map[*Client]bool
	lock  sync.RWMutex

	broadcast chan OutgoingMessage
}

func NewHub() *Hub {
	return &Hub{
		rooms:     make(map[string]map[*Client]bool),
		broadcast: make(chan OutgoingMessage, 1024),
	}
}

func (h *Hub) Run() {
	for msg := range h.broadcast {
		h.lock.RLock()
		clients := h.rooms[msg.Room]
		for c := range clients {
			if err := c.conn.WriteJSON(msg); err != nil {
				log.Printf("write error: %v", err)
				c.conn.Close()
			}
		}
		h.lock.RUnlock()
	}
}

func (h *Hub) AddClient(c *Client) {
	h.lock.Lock()
	defer h.lock.Unlock()
	if h.rooms[c.RoomName] == nil {
		h.rooms[c.RoomName] = make(map[*Client]bool)
	}
	h.rooms[c.RoomName][c] = true
}

func (h *Hub) RemoveClient(c *Client) {
	h.lock.Lock()
	defer h.lock.Unlock()
	if set, ok := h.rooms[c.RoomName]; ok {
		delete(set, c)
		if len(set) == 0 {
			delete(h.rooms, c.RoomName)
		}
	}
}

func (h *Hub) ActiveRooms() []string {
	h.lock.RLock()
	defer h.lock.RUnlock()
	out := make([]string, 0, len(h.rooms))
	for name := range h.rooms {
		out = append(out, name)
	}
	return out
}

type IncomingMessage struct {
	Type     string `json:"type"`
	Room     string `json:"room"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

type OutgoingMessage struct {
	Type     string `json:"type"`
	Room     string `json:"room"`
	Username string `json:"username"`
	Message  string `json:"message"`
}
