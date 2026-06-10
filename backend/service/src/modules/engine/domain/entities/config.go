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
// Combined with an MQTT event to build a publish.
type MQTTConnectionConfig struct {
	BrokerURL string `json:"brokerUrl"`
	ClientID  string `json:"clientId"`
	BaseTopic string `json:"baseTopic"`
	AuthMode  string `json:"authMode"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}
