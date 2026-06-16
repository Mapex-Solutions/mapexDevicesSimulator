// Package constants centralizes the e2e configuration as plain Go values — no
// environment variables, nothing to set in the shell. A run is just `go test`.
// The endpoints mirror the simulator sidecar's own defaults
// (service/src/shared/configuration/application/config.go) and the ChirpStack
// deployment (deployment/chirpstack/chirpstack.yml); keep them in sync if either
// changes.
package constants

const (
	// SimURL is the simulator sidecar base URL the e2e drives. It matches the
	// service's default http_address (127.0.0.1) and http_port (5055).
	SimURL = "http://127.0.0.1:5055"

	// ConsoleWSURL is the realtime console WebSocket the e2e consumes for live
	// frames (up/down/system) and connection-status events.
	ConsoleWSURL = "ws://127.0.0.1:5055/ws"

	// ChirpStackGRPCAddr is the gRPC API the LoRaWAN journey provisions over.
	ChirpStackGRPCAddr = "127.0.0.1:18080"
	// ChirpStackAdminUser and ChirpStackAdminPass are the default v4 credentials.
	ChirpStackAdminUser = "admin"
	ChirpStackAdminPass = "admin"
	// ChirpStackUDPHost and ChirpStackUDPPort are the Semtech UDP packet-forwarder
	// endpoint the simulator's UDP gateway transmits to.
	ChirpStackUDPHost = "127.0.0.1"
	ChirpStackUDPPort = 11700
	// ChirpStackBasicStationURL is the Basics Station WebSocket endpoint the
	// simulator's Basics Station device connects to.
	ChirpStackBasicStationURL = "ws://127.0.0.1:13001"
)
