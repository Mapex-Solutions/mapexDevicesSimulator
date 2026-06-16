//go:build saga

package mqtt_downlink

import "testing"

// TestJourney exercises the MQTT inbound path: a receive-enabled device gets a
// retained downlink published on its subscribed topic and the simulator logs it
// as a down/downlink frame. Requires the simulator running on 127.0.0.1:5055;
// see the journey README.
func TestJourney(t *testing.T) {
	Run(t)
}
