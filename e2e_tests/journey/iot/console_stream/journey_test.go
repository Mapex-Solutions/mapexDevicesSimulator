//go:build saga

package console_stream

import "testing"

// TestJourney connects to the console WebSocket, fires an HTTP device, and
// verifies the uplink arrives live on the stream. Requires the simulator running
// on 127.0.0.1:5055; see the journey README.
func TestJourney(t *testing.T) {
	Run(t)
}
