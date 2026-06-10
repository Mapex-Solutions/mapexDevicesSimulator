package hub

import (
	"encoding/json"

	"simulator/service/src/modules/console/application/dtos"
	"simulator/service/src/modules/console/application/ports"
)

// Compile-time proof the hub satisfies the Publisher port.
var _ ports.Publisher = (*Hub)(nil)

// New builds an empty hub.
func New() *Hub {
	return &Hub{clients: make(map[*Client]struct{})}
}

// Register adds a subscriber and returns its client handle.
func (h *Hub) Register() *Client {
	cl := &Client{send: make(chan []byte, clientBuffer)}
	h.mu.Lock()
	h.clients[cl] = struct{}{}
	h.mu.Unlock()
	return cl
}

// Unregister removes a subscriber and closes its channel. Safe to call twice.
func (h *Hub) Unregister(cl *Client) {
	h.mu.Lock()
	if _, ok := h.clients[cl]; ok {
		delete(h.clients, cl)
		close(cl.send)
	}
	h.mu.Unlock()
}

// Publish broadcasts a frame to every client, dropping it for any client whose
// buffer is full (a slow consumer never stalls the engine).
func (h *Hub) Publish(msg dtos.ConsoleMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	for cl := range h.clients {
		select {
		case cl.send <- data:
		default:
		}
	}
}

// Out exposes a client's outbound frame channel for the WS write pump.
func (c *Client) Out() <-chan []byte {
	return c.send
}
