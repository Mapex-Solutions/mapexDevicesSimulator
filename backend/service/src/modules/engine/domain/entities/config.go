package entities

// HTTPConnectionConfig is the HTTP target on a device's config (the base URL plus
// auth). Combined with an HTTP event to build a request.
type HTTPConnectionConfig struct {
	URL          string     `json:"url"`
	Method       string     `json:"method"`
	Headers      []KeyValue `json:"headers"`
	AuthMode     string     `json:"authMode"`
	APIKeyHeader string     `json:"apiKeyHeader"`
	APIKey       string     `json:"apiKey"`
	BasicUser    string     `json:"basicUser"`
	BasicPass    string     `json:"basicPass"`
}

// MQTTConnectionConfig is the MQTT target on a device's config (broker + auth).
// Combined with an MQTT event to build a publish, and (when ReceiveEnabled) it
// drives the persistent session's downlink subscriptions.
type MQTTConnectionConfig struct {
	BrokerURL      string             `json:"brokerUrl"`
	ClientID       string             `json:"clientId"`
	BaseTopic      string             `json:"baseTopic"`
	AuthMode       string             `json:"authMode"`
	Username       string             `json:"username"`
	Password       string             `json:"password"`
	ReceiveEnabled bool               `json:"receiveEnabled"`
	Subscriptions  []MQTTSubscription `json:"subscriptions"`
}

// MQTTSubscription is one downlink topic the device subscribes to when receiving
// is enabled. The topic may carry a {{baseTopic}}/{{deviceId}} prefix convention,
// but is stored verbatim as authored.
type MQTTSubscription struct {
	Name  string `json:"name"`
	Topic string `json:"topic"`
	QoS   int    `json:"qos"`
}

// LoRaWANConnectionConfig is a LoRaWAN node that transmits through a shared gateway
// referenced by id. Keys are flat hex; OTAA fields drive the join, ABP fields a
// pre-provisioned session.
type LoRaWANConnectionConfig struct {
	GatewayID  string `json:"gatewayId"`
	Region     string `json:"region"`
	MACVersion string `json:"macVersion"`
	Activation string `json:"activation"`
	DevEUI     string `json:"devEui"`
	JoinEUI    string `json:"joinEui"`
	AppKey     string `json:"appKey"`
	NwkKey     string `json:"nwkKey"`
	DevAddr    string `json:"devAddr"`
	NwkSKey    string `json:"nwkSKey"`
	AppSKey    string `json:"appSKey"`
}

// BasicsStationConnectionConfig is a LoRaWAN node carrying its own Basics Station
// link (its embedded gateway), instead of attaching to a separate gateway.
type BasicsStationConnectionConfig struct {
	LNSURI     string `json:"lnsUri"`
	GatewayEUI string `json:"gatewayEui"`
	Region     string `json:"region"`
	MACVersion string `json:"macVersion"`
	Activation string `json:"activation"`
	DevEUI     string `json:"devEui"`
	JoinEUI    string `json:"joinEui"`
	AppKey     string `json:"appKey"`
	NwkKey     string `json:"nwkKey"`
	DevAddr    string `json:"devAddr"`
	NwkSKey    string `json:"nwkSKey"`
	AppSKey    string `json:"appSKey"`
}

// GatewayLink is the gateway's connection to the LNS, parsed from the gateway's
// link JSON. Only the fields relevant to the chosen protocol are used.
type GatewayLink struct {
	Protocol string `json:"protocol"`
	LNSURI   string `json:"lnsUri"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}
