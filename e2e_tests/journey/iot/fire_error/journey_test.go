//go:build saga

package fire_error

import "testing"

// TestJourney fires a device at an unreachable target and verifies the engine
// logs the failure as an error frame. Requires the simulator running on
// 127.0.0.1:5055; see the journey README.
func TestJourney(t *testing.T) {
	Run(t)
}
