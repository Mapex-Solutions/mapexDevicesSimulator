package utils

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"

	consolecontract "simulator/packages/contracts/console"
)

// ConsoleStream is a connected reader of the simulator's realtime console
// WebSocket. It collects every frame the sidecar broadcasts so a test can assert
// on what reached the UI live — uplinks, downlinks, and connection-status events,
// some of which (the reconnect lifecycle) exist only on this stream, never in the
// logs.
type ConsoleStream struct {
	conn   *websocket.Conn
	mu     sync.Mutex
	frames []consolecontract.ConsoleMessage
}

// StartConsoleStream dials the console WebSocket and starts collecting frames in
// the background. The caller closes it when done.
func StartConsoleStream(wsURL string) (*ConsoleStream, error) {
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return nil, err
	}
	s := &ConsoleStream{conn: conn}
	go s.readLoop()
	return s, nil
}

// readLoop appends each received frame until the connection closes.
func (s *ConsoleStream) readLoop() {
	for {
		_, data, err := s.conn.ReadMessage()
		if err != nil {
			return
		}
		var msg consolecontract.ConsoleMessage
		if json.Unmarshal(data, &msg) != nil {
			continue
		}
		s.mu.Lock()
		s.frames = append(s.frames, msg)
		s.mu.Unlock()
	}
}

// Frames returns a copy of every frame received so far.
func (s *ConsoleStream) Frames() []consolecontract.ConsoleMessage {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]consolecontract.ConsoleMessage, len(s.frames))
	copy(out, s.frames)
	return out
}

// Close tears the WebSocket down.
func (s *ConsoleStream) Close() { _ = s.conn.Close() }
