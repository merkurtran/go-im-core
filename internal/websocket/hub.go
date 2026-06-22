package websocket

import "errors"

type Hub struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.userID] = client
		case client := <-h.unregister:
			if _, ok := h.clients[client.userID]; ok {
				delete(h.clients, client.userID)
				close(client.send)
			}
		}
	}
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

func (h *Hub) SendToUser(userID string, data []byte) error {
	client, ok := h.clients[userID]
	if !ok {
		return errors.New("user offline")
	}
	select {
	case client.send <- data:
		return nil
	default:
		close(client.send)
		delete(h.clients, userID)
		return errors.New("send channel full")
	}
}
