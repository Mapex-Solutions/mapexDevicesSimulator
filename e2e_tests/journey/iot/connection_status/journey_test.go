//go:build saga

package connection_status

import "testing"

// TestJourney enables an MQTT device pointed at an unreachable broker and
// verifies the engine surfaces connecting/reconnecting status frames live on the
// console WebSocket. Requires the simulator running on 127.0.0.1:5055; see the
// journey README.
func TestJourney(t *testing.T) {
	Run(t)
}
