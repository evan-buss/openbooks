package server

import "github.com/google/uuid"

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[uuid.UUID]*Client

	// Inbound messages from the clients.
	shutdown chan struct{}

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		shutdown:   make(chan struct{}),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[uuid.UUID]*Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.uuid] = client
		case client := <-h.unregister:
			if _, ok := h.clients[client.uuid]; ok {
				close(client.send)
				close(client.disconnect)
				delete(h.clients, client.uuid)
			}
		case <-h.shutdown:
			for _, client := range h.clients {
				close(client.send)
				close(client.disconnect)
				delete(h.clients, client.uuid)
			}
			return
		}
	}
}
