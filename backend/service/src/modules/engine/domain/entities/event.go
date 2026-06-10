package entities

// KeyValue is a header or form-field pair.
type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// RequestBody is the body authoring shared by HTTP and MQTT events. Values may
// carry {{placeholders}} resolved at send time.
type RequestBody struct {
	BodyMode   string     `json:"bodyMode"`
	BodyFields []KeyValue `json:"bodyFields"`
	Body       string     `json:"body"`
}

// HTTPEventConfig is the HTTP event payload.
type HTTPEventConfig struct {
	RequestBody
	Method  string     `json:"method"`
	Path    string     `json:"path"`
	Headers []KeyValue `json:"headers"`
}

// MQTTEventConfig is the MQTT event payload.
type MQTTEventConfig struct {
	RequestBody
	Topic  string `json:"topic"`
	QoS    int    `json:"qos"`
	Retain bool   `json:"retain"`
}

// LoRaWANEventConfig is the LoRaWAN uplink event.
type LoRaWANEventConfig struct {
	FPort      int    `json:"fport"`
	Confirmed  bool   `json:"confirmed"`
	PayloadHex string `json:"payloadHex"`
}

// EventSchedule is an event's optional auto-fire schedule.
type EventSchedule struct {
	Enabled bool   `json:"enabled"`
	Every   int    `json:"every"`
	Unit    string `json:"unit"`
}

// DeviceEvent is one pre-registered event on a device, parsed from the device's
// events JSON. It holds the protocol-specific config for the device's protocol
// plus an optional schedule.
type DeviceEvent struct {
	ID       string              `json:"id"`
	Name     string              `json:"name"`
	HTTP     *HTTPEventConfig    `json:"http,omitempty"`
	MQTT     *MQTTEventConfig    `json:"mqtt,omitempty"`
	LoRaWAN  *LoRaWANEventConfig `json:"lorawan,omitempty"`
	Schedule *EventSchedule      `json:"schedule,omitempty"`
}
