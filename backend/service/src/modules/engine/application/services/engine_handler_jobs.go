package services

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	devicescontract "simulator/packages/contracts/devices"
	"simulator/service/src/modules/engine/domain/entities"
)

// buildSendSpec resolves a device + event into the static parts of a send (the
// target and auth from the device config merged with the event), returning false
// when the protocol is not enabled yet or the config is unusable. The payload,
// url and topic templates are rendered later, at fire time.
func buildSendSpec(d devicescontract.Device, e entities.DeviceEvent) (sendSpec, bool) {
	spec := sendSpec{
		protocol:   d.ProtocolID,
		deviceKey:  d.ID,
		deviceID:   d.DeviceID,
		deviceName: d.Name,
		eventName:  e.Name,
		storeLogs:  d.StoreLogs,
	}
	switch d.ProtocolID {
	case "http":
		if e.HTTP == nil {
			return sendSpec{}, false
		}
		var cfg entities.HTTPConnectionConfig
		if err := json.Unmarshal(d.Config, &cfg); err != nil || cfg.URL == "" {
			return sendSpec{}, false
		}
		spec.url = joinURL(cfg.URL, e.HTTP.Path)
		spec.method = firstNonEmpty(e.HTTP.Method, cfg.Method, "POST")
		spec.headers = buildHTTPHeaders(cfg, e.HTTP.Headers)
		spec.payloadTemplate = eventBody(e.HTTP.RequestBody)
		return spec, true
	case "mqtt":
		if e.MQTT == nil {
			return sendSpec{}, false
		}
		var cfg entities.MQTTConnectionConfig
		if err := json.Unmarshal(d.Config, &cfg); err != nil || cfg.BrokerURL == "" {
			return sendSpec{}, false
		}
		spec.brokerURL = cfg.BrokerURL
		spec.clientID = cfg.ClientID
		spec.username = cfg.Username
		spec.password = cfg.Password
		spec.tlsCert = cfg.TLSCertPem
		spec.tlsKey = cfg.TLSKeyPem
		spec.tlsCa = cfg.TLSCaPem
		spec.topic = joinTopic(cfg.BaseTopic, e.MQTT.Topic)
		spec.qos = byte(e.MQTT.QoS)
		spec.retain = e.MQTT.Retain
		spec.payloadTemplate = eventBody(e.MQTT.RequestBody)
		return spec, true
	case "lorawan", "basicstation":
		if e.LoRaWAN == nil {
			return sendSpec{}, false
		}
		spec.fport = byte(e.LoRaWAN.FPort)
		spec.confirmed = e.LoRaWAN.Confirmed
		spec.payloadTemplate = e.LoRaWAN.PayloadHex
		return spec, true
	default:
		return sendSpec{}, false
	}
}

// joinURL appends an event path to a device base URL.
func joinURL(base, path string) string {
	if path == "" {
		return base
	}
	return strings.TrimRight(base, "/") + "/" + strings.TrimLeft(path, "/")
}

// joinTopic prefixes an event topic with the device base topic.
func joinTopic(base, topic string) string {
	switch {
	case base == "":
		return topic
	case topic == "":
		return base
	default:
		return strings.TrimRight(base, "/") + "/" + strings.TrimLeft(topic, "/")
	}
}

// firstNonEmpty returns the first non-empty value.
func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

// buildHTTPHeaders merges device headers, event headers, and the device's auth.
func buildHTTPHeaders(cfg entities.HTTPConnectionConfig, eventHeaders []entities.KeyValue) map[string]string {
	h := make(map[string]string)
	for _, kv := range cfg.Headers {
		h[kv.Key] = kv.Value
	}
	for _, kv := range eventHeaders {
		h[kv.Key] = kv.Value
	}
	switch cfg.AuthMode {
	case "apiKey":
		if cfg.APIKeyHeader != "" {
			h[cfg.APIKeyHeader] = cfg.APIKey
		}
	case "basic":
		token := base64.StdEncoding.EncodeToString([]byte(cfg.BasicUser + ":" + cfg.BasicPass))
		h["Authorization"] = "Basic " + token
	}
	return h
}

// eventBody resolves the body template from the event's body mode.
func eventBody(rb entities.RequestBody) string {
	switch rb.BodyMode {
	case "raw":
		return rb.Body
	case "form":
		return formToJSON(rb.BodyFields)
	default:
		return ""
	}
}

// formToJSON renders form fields as a JSON object template (values may still
// carry placeholders, resolved at fire time).
func formToJSON(fields []entities.KeyValue) string {
	m := make(map[string]string, len(fields))
	for _, kv := range fields {
		m[kv.Key] = kv.Value
	}
	b, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(b)
}
