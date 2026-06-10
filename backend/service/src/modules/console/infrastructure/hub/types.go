package hub

import "sync"

// Hub is an in-memory fan-out broker: connected WebSocket clients register here,
// and Publish broadcasts a frame to all of them. It is the console module's
// Publisher implementation.
type Hub struct {
	mu      sync.RWMutex
	clients map[*Client]struct{}
}

// Client is one connected console subscriber. The WS handler reads frames off
// Out() and writes them to the socket; a slow client drops frames rather than
// stalling the broadcaster.
type Client struct {
	send chan []byte
}
