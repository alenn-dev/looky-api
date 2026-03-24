package ws

import "sync"

type Hub struct {
	clients map[string]*Client
	mu      sync.RWMutex
}

var H = &Hub{
	clients: make(map[string]*Client),
}

func (h *Hub) Register(customerID string, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[customerID] = client
}

func (h *Hub) UnRegister(customerID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, customerID)
}

func (h *Hub) Send(customerID string, message []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if client, ok := h.clients[customerID]; ok {
		client.send <- message
	}
}
