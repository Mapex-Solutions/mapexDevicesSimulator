package dispatch

import "time"

const (
	// httpTimeout caps one outbound HTTP send.
	httpTimeout = 15 * time.Second

	// httpResponseCaptureLimit caps how much of the HTTP response body is kept to
	// show on the console, so a large or streaming reply never blows up a frame.
	httpResponseCaptureLimit = 8 << 10 // 8 KiB

	// mqttConnectTimeout caps the MQTT handshake per send.
	mqttConnectTimeout = 10 * time.Second

	// mqttQuiesceMillis is how long Disconnect waits for in-flight work.
	mqttQuiesceMillis = 250
)
