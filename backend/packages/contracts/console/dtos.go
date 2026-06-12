package console

// ConsoleMessage is one frame on the live console stream: every uplink, downlink
// and auth/join handshake the engine emits, across every protocol. `ts` is the
// frame's event time on the live stream (not a persisted entity's created).
type ConsoleMessage struct {
	ID         string            `json:"id"`
	TS         string            `json:"ts"`
	Protocol   string            `json:"protocol"`
	DeviceID   string            `json:"deviceId"`
	DeviceName string            `json:"deviceName"`
	Direction  string            `json:"direction"`
	Kind       string            `json:"kind"`
	Summary    string            `json:"summary"`
	Payload    string            `json:"payload"`
	Response   string            `json:"response,omitempty"` // endpoint's reply (HTTP body); empty otherwise
	Status     string            `json:"status,omitempty"`
	Meta       map[string]string `json:"meta,omitempty"`
}
