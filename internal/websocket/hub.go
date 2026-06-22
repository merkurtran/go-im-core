package websocket

import (
	"errors"
	"log/slog"
	"sync"
)

// Hub 管理所有 WebSocket 连接，支持一个用户多端在线。
type Hub struct {
	clients map[string]map[*Client]bool // userID -> set of clients
	mu      sync.RWMutex

	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte, 256),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)
		case client := <-h.unregister:
			h.unregisterClient(client)
		case msg := <-h.broadcast:
			h.broadcastToAll(msg)
		}
	}
}

func (h *Hub) registerClient(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	set, ok := h.clients[c.userID]
	if !ok {
		set = make(map[*Client]bool)
		h.clients[c.userID] = set
	}
	set[c] = true
	slog.Info("client registered", "user_id", c.userID, "total_clients", len(set))
}

func (h *Hub) unregisterClient(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if set, ok := h.clients[c.userID]; ok {
		if _, ok := set[c]; ok {
			delete(set, c)
			close(c.send)
			if len(set) == 0 {
				delete(h.clients, c.userID)
			}
		}
	}
}

func (h *Hub) broadcastToAll(data []byte) {
	h.mu.RLock()
	allClients := make([]*Client, 0)
	for _, set := range h.clients {
		for c := range set {
			allClients = append(allClients, c)
		}
	}
	h.mu.RUnlock()

	for _, client := range allClients {
		select {
		case client.send <- data:
		default:
			go h.Unregister(client)
		}
	}
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

// SendToUser 向指定用户的所有在线端推送消息。
func (h *Hub) SendToUser(userID string, data []byte) error {
	h.mu.RLock()
	set, ok := h.clients[userID]
	if !ok || len(set) == 0 {
		h.mu.RUnlock()
		return errors.New("user offline")
	}
	clients := make([]*Client, 0, len(set))
	for c := range set {
		clients = append(clients, c)
	}
	h.mu.RUnlock()

	var failed int
	for _, client := range clients {
		select {
		case client.send <- data:
		default:
			go h.Unregister(client)
			failed++
		}
	}
	if failed == len(clients) {
		return errors.New("all client send channels full")
	}
	return nil
}

// SendToUserAsync 异步发送，不阻塞调用方
func (h *Hub) SendToUserAsync(userID string, data []byte) {
	go func() {
		if err := h.SendToUser(userID, data); err != nil {
			slog.Warn("SendToUser failed", "user_id", userID, "error", err)
		}
	}()
}
