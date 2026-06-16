//go:build saga

package http_device_fire

import "testing"

// TestJourney fires an HTTP device against an in-test echo and verifies the
// simulator logged the 200 uplink with its response. Requires the simulator
// running on 127.0.0.1:5055; see the journey README.
func TestJourney(t *testing.T) {
	Run(t)
}
