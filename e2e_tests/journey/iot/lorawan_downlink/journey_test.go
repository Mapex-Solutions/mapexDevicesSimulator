//go:build saga

package lorawan_downlink

import "testing"

// TestJourney provisions ChirpStack, joins a UDP LoRaWAN device, queues a
// downlink on the LNS, fires an uplink to open the RX window, and verifies the
// simulator logs the received downlink. Requires the sidecar on 127.0.0.1:5055
// and docker (the journey brings the ChirpStack stack up and tears it down);
// see the journey README.
func TestJourney(t *testing.T) {
	Run(t)
}
