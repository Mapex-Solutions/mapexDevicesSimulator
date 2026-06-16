package ports

import "context"

// SessionSpec is the resolved target + credentials to open ONE live connection for
// a device. Unlike a DispatchRequest (one-shot), a session stays open so the device
// can both publish uplinks through it and receive downlinks on it.
type SessionSpec struct {
	Protocol   string
	DeviceID   string
	DeviceName string
	StoreLogs  bool

	// mqtt
	BrokerURL string
	ClientID  string
	Username  string
	Password  string
	// TLS client-cert material, in PEM. Set for an ssl:// broker with
	// certificate auth; the connector builds the tls.Config from them.
	TLSCert       string
	TLSKey        string
	TLSCa         string
	Subscriptions []Subscription // downlink topics, when the device has receive enabled

	// lorawan (the device's crypto material + the gateway link it transmits through)
	LoRaWAN *LoRaWANSpec
}

// LoRaWANSpec is the per-device LoRaWAN material plus the shared gateway link it
// forwards through. The gateway fields are resolved from the device's gateway at
// reconcile time so the connector can share one link across the gateway's devices.
type LoRaWANSpec struct {
	Region     string
	MACVersion string
	Activation string // otaa | abp
	Class      string // A | C

	JoinEUI string
	DevEUI  string
	AppKey  string
	NwkKey  string
	DevAddr string
	NwkSKey string
	AppSKey string

	// Gateway link (shared by devices on the same gateway).
	GatewayEUI    string
	LinkProtocol  string // basicstation | udp
	LinkLNSURI    string // basicstation: wss://host:port
	LinkUDPHost   string // udp: host
	LinkUDPPort   int    // udp: port
}

// Subscription is one downlink topic the device listens on.
type Subscription struct {
	Name  string
	Topic string
	QoS   byte
}

// OutboundMessage is one uplink sent through a live session (the session already
// owns the connection, so only the per-message fields are needed). MQTT uses
// Topic/QoS/Retain; LoRaWAN uses FPort/Confirmed and a hex Payload.
type OutboundMessage struct {
	Topic     string
	QoS       byte
	Retain    bool
	Payload   string
	FPort     byte
	Confirmed bool
}

// SendResult reports the outcome of one session send, for the console/log line.
type SendResult struct {
	OK     bool
	Status string // protocol status, e.g. "qos1"
	Err    error
}

// InboundMessage is one received downlink the session surfaces to the engine.
type InboundMessage struct {
	Topic   string
	Payload string
	Status  string // protocol status, e.g. "received qos1"
	Summary string // one-line description, e.g. the topic
}

// InboundSink is called by a session on every received message. It runs on the
// transport's callback goroutine, so the implementation must be non-blocking.
type InboundSink func(InboundMessage)

// StatusSink is called by a connector/session on each connection-lifecycle
// transition (subscribing, subscribed, ...). The session manager emits the outer
// lifecycle (connecting/connected/reconnecting/disconnected) itself.
type StatusSink func(status, detail string)

// Session is a live, per-device connection. Send publishes an uplink through it;
// Close tears it down.
type Session interface {
	Send(ctx context.Context, msg OutboundMessage) SendResult
	Close() error
	Connected() bool
}

// Connector opens a live session for one session-capable protocol (mqtt today;
// lorawan, via a gateway link, in a later phase). It is the persistent counterpart
// of the one-shot Dispatcher.
type Connector interface {
	Protocol() string
	Open(ctx context.Context, spec SessionSpec, in InboundSink, status StatusSink) (Session, error)
}

// ConnectorRegistry resolves a connector by protocol id. A protocol with no
// connector is not session-capable (e.g. http), and uses the one-shot Dispatcher.
type ConnectorRegistry interface {
	Connector(protocol string) (Connector, bool)
}
