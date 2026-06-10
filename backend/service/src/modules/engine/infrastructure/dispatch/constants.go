package dispatch

import "time"

const (
	// httpTimeout caps one outbound HTTP send.
	httpTimeout = 15 * time.Second

	// mqttConnectTimeout caps the MQTT handshake per send.
	mqttConnectTimeout = 10 * time.Second

	// mqttQuiesceMillis is how long Disconnect waits for in-flight work.
	mqttQuiesceMillis = 250
)
