// Package payloads holds the LoRaWAN identities the journey provisions on both
// ChirpStack and the simulator. The two sides MUST agree on DevEUI/JoinEUI/AppKey
// for the OTAA join to succeed, so a single identity object is the source of
// truth for both.
//
// The DevEUI and gateway EUI are derived from the run id so every run uses fresh
// identities. This is essential: the simulator sidecar is a long-lived shared
// process that persists per-DevEUI join state, and ChirpStack remembers OTAA
// DevNonces per device — reusing a fixed DevEUI across runs would make the
// second run's join look like a replay and be dropped. Per-run EUIs sidestep
// both, independent of any stack teardown.
package payloads

import (
	"crypto/sha256"
	"encoding/hex"
)

// LoRaIdentity is the shared OTAA material for one LoRaWAN device plus the EUI of
// the gateway it transmits through.
type LoRaIdentity struct {
	GatewayEUI string
	DevEUI     string
	JoinEUI    string
	AppKey     string
}

// SagaUDPIdentity is the run-scoped identity for the Semtech UDP transport.
func SagaUDPIdentity(runID string) LoRaIdentity {
	return LoRaIdentity{
		GatewayEUI: eui(runID, "udp-gw"),
		DevEUI:     eui(runID, "udp-dev"),
		JoinEUI:    "0000000000000000",
		AppKey:     "00112233445566778899AABBCCDDEEFF",
	}
}

// SagaBasicStationIdentity is the run-scoped identity for the Basics Station
// transport, distinct from the UDP one so both coexist on a single stack.
func SagaBasicStationIdentity(runID string) LoRaIdentity {
	return LoRaIdentity{
		GatewayEUI: eui(runID, "bs-gw"),
		DevEUI:     eui(runID, "bs-dev"),
		JoinEUI:    "0000000000000000",
		AppKey:     "112233445566778899AABBCCDDEEFF00",
	}
}

// eui derives a stable 16-hex (8-byte) EUI from the run id and a role salt, so
// each run and each role (gateway vs device, UDP vs Basics Station) gets a
// distinct, valid identifier.
func eui(runID, salt string) string {
	sum := sha256.Sum256([]byte(runID + ":" + salt))
	return hex.EncodeToString(sum[:8])
}
