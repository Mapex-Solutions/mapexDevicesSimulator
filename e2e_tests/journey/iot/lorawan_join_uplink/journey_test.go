//go:build saga

package lorawan_join_uplink

import "testing"

// TestJourney provisions ChirpStack and drives LoRaWAN OTAA devices over both
// Semtech UDP and Basics Station, verifying the LNS records the join and the
// uplink for each. Requires the sidecar on 127.0.0.1:5055 and docker (the journey brings
// the ChirpStack stack up and tears it down); see the journey README.
func TestJourney(t *testing.T) {
	Run(t)
}
