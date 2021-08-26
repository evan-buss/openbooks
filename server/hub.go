package server

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

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
		// TODO: maybe use map[uuid.UUID]*Client so that we can look up users by
		// UUID, this would allow us to migrate some WS calls to HTTP because we
		// can send the UUID and the hub will take care of calling the correct
		// client IRC and WS connection.
		clients: make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				close(client.send)
				close(client.disconnect)
				delete(h.clients, client)
			}
		case <-h.shutdown:
			for client := range h.clients {
				close(client.send)
				close(client.disconnect)
				delete(h.clients, client)
			}
			return
		}
	}
}
