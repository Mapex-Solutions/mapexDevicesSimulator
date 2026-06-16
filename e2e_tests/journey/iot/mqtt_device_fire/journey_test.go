//go:build saga

package mqtt_device_fire

import "testing"

// TestJourney fires MQTT devices against an in-process broker over both auth
// modes — username/password and client certificate — and verifies the broker
// accepted each authenticated publish. Requires the simulator running on
// 127.0.0.1:5055; see the journey README.
func TestJourney(t *testing.T) {
	Run(t)
}
