package ports

import "context"

// DispatchRequest is a fully-resolved message to send: the engine renders the
// payload and fills the fields relevant to the device's protocol, and each
// dispatcher reads only what it needs.
type DispatchRequest struct {
	Payload string

	// HTTP
	URL     string
	Method  string
	Headers map[string]string

	// MQTT
	BrokerURL string
	ClientID  string
	Username  string
	Password  string
	Topic     string
	QoS       byte
	Retain    bool
}

// DispatchResult reports the outcome of one send, for the console/log line.
type DispatchResult struct {
	OK       bool
	Status   string // protocol status, e.g. "200" or "qos1"
	Response string // body the endpoint replied with (HTTP), capped; empty otherwise
	Err      error
}

// Dispatcher sends a resolved message over one protocol.
type Dispatcher interface {
	Protocol() string
	Dispatch(ctx context.Context, req DispatchRequest) DispatchResult
}

// Registry resolves a dispatcher by protocol id.
type Registry interface {
	For(protocol string) (Dispatcher, bool)
}
