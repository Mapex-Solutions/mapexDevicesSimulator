// Package payloads holds canonical device-create fixtures for the simulator's
// devices module. Builders wrap the create contract so steps override single
// fields without redeclaring the baseline. Builders are pure: a step reads the
// runtime values it needs from the bag and passes them in.
package payloads

import (
	"encoding/json"
	"fmt"

	devicescontract "simulator/packages/contracts/devices"
)

// DeviceInputBuilder wraps devicescontract.DeviceInput.
type DeviceInputBuilder struct {
	spec devicescontract.DeviceInput
}

// Build returns the contract payload ready for POST /api/devices.
func (b *DeviceInputBuilder) Build() devicescontract.DeviceInput { return b.spec }

// WithEnabled overrides the enabled flag.
func (b *DeviceInputBuilder) WithEnabled(v bool) *DeviceInputBuilder {
	b.spec.Enabled = v
	return b
}

// SagaHTTPDevice builds a send-only HTTP device targeting echoURL, with one
// pre-registered event carrying a templated JSON body. storeLogs is on so the
// fire lands in the logs the journey asserts on.
func SagaHTTPDevice(runID, echoURL string) *DeviceInputBuilder {
	config := mustJSON(map[string]any{
		"kind":     "http",
		"url":      echoURL,
		"method":   "POST",
		"headers":  []any{map[string]string{"key": "Content-Type", "value": "application/json"}},
		"authMode": "none",
	})
	events := mustJSON([]any{
		map[string]any{
			"id":   "e1",
			"name": "Telemetry",
			"http": map[string]any{
				"method":     "POST",
				"path":       "/ingest",
				"headers":    []any{},
				"bodyMode":   "raw",
				"bodyFields": []any{},
				"body":       `{"deviceId":"{{deviceId}}","temp":{{randInt(1,9)}}}`,
			},
		},
	})
	return &DeviceInputBuilder{spec: devicescontract.DeviceInput{
		Name:       fmt.Sprintf("e2e HTTP %s", runID),
		DeviceID:   fmt.Sprintf("e2e-http-%s", runID),
		ProtocolID: "http",
		Enabled:    true,
		StoreLogs:  true,
		Config:     config,
		Attributes: map[string]string{},
		Events:     events,
	}}
}

// MQTTBrokerTarget is the broker coordinates an MQTT device fixture connects to.
// The certificate variant fills the PEM fields; the username/password variant
// fills Username/Password and leaves the PEM fields empty.
type MQTTBrokerTarget struct {
	BrokerURL  string
	Username   string
	Password   string
	TLSCertPem string
	TLSKeyPem  string
	TLSCaPem   string
}

// mqttDevice builds an enabled MQTT device against target with one pre-registered
// event publishing a templated JSON body to a per-device topic.
func mqttDevice(runID, slug, authMode string, target MQTTBrokerTarget) *DeviceInputBuilder {
	config := mustJSON(map[string]any{
		"kind":           "mqtt",
		"brokerUrl":      target.BrokerURL,
		"clientId":       fmt.Sprintf("e2e-%s-%s", slug, runID),
		"baseTopic":      "e2e",
		"authMode":       authMode,
		"username":       target.Username,
		"password":       target.Password,
		"tlsCertPem":     target.TLSCertPem,
		"tlsKeyPem":      target.TLSKeyPem,
		"tlsCaPem":       target.TLSCaPem,
		"receiveEnabled": false,
		"subscriptions":  []any{},
	})
	events := mustJSON([]any{
		map[string]any{
			"id":   "e1",
			"name": "Telemetry",
			"mqtt": map[string]any{
				"topic":      "{{deviceId}}/telemetry",
				"qos":        1,
				"retain":     false,
				"bodyMode":   "raw",
				"bodyFields": []any{},
				"body":       `{"deviceId":"{{deviceId}}","temp":{{randInt(1,9)}}}`,
			},
		},
	})
	return &DeviceInputBuilder{spec: devicescontract.DeviceInput{
		Name:       fmt.Sprintf("e2e MQTT %s %s", slug, runID),
		DeviceID:   fmt.Sprintf("e2e-mqtt-%s-%s", slug, runID),
		ProtocolID: "mqtt",
		Enabled:    true,
		StoreLogs:  true,
		Config:     config,
		Attributes: map[string]string{},
		Events:     events,
	}}
}

// SagaMQTTDeviceUserPass builds an MQTT device that authenticates with a
// username and password over a plain tcp:// connection.
func SagaMQTTDeviceUserPass(runID string, target MQTTBrokerTarget) *DeviceInputBuilder {
	return mqttDevice(runID, "userpass", "userpass", target)
}

// SagaMQTTDeviceTLS builds an MQTT device that authenticates with a client
// certificate over an ssl:// connection (mutual TLS).
func SagaMQTTDeviceTLS(runID string, target MQTTBrokerTarget) *DeviceInputBuilder {
	return mqttDevice(runID, "tls", "certificate", target)
}

// MQTTDownlinkTopic is the concrete topic an MQTT receive device subscribes to,
// derived from the run id so each run is isolated. The downlink publisher and
// the device agree on it: it is joinTopic("e2e", "dl/<runID>").
func MQTTDownlinkTopic(runID string) string { return "e2e/dl/" + runID }

// SagaMQTTReceiveDevice builds an enabled MQTT device with receiving on: it
// subscribes to MQTTDownlinkTopic over a plain tcp:// connection, with storeLogs
// on so a received downlink lands in the logs the journey asserts on. It carries
// no send events — the inbound path is what is exercised.
func SagaMQTTReceiveDevice(runID, brokerURL, user, pass string) *DeviceInputBuilder {
	config := mustJSON(map[string]any{
		"kind":           "mqtt",
		"brokerUrl":      brokerURL,
		"clientId":       fmt.Sprintf("e2e-rx-%s", runID),
		"baseTopic":      "e2e",
		"authMode":       "userpass",
		"username":       user,
		"password":       pass,
		"tlsCertPem":     "",
		"tlsKeyPem":      "",
		"tlsCaPem":       "",
		"receiveEnabled": true,
		"subscriptions": []any{
			map[string]any{"name": "downlink", "topic": "dl/" + runID, "qos": 1},
		},
	})
	return &DeviceInputBuilder{spec: devicescontract.DeviceInput{
		Name:       fmt.Sprintf("e2e MQTT RX %s", runID),
		DeviceID:   fmt.Sprintf("e2e-mqtt-rx-%s", runID),
		ProtocolID: "mqtt",
		Enabled:    true,
		StoreLogs:  true,
		Config:     config,
		Attributes: map[string]string{},
		Events:     mustJSON([]any{}),
	}}
}

// LoRaTarget is the OTAA material plus the routing a LoRaWAN device needs. The
// keys must match what was provisioned on the LNS. GatewayID is the simulator
// gateway UUID for the UDP transport; GatewayEUI and LNSURI are set for the
// Basics Station transport instead.
type LoRaTarget struct {
	GatewayID  string
	GatewayEUI string
	LNSURI     string
	DevEUI     string
	JoinEUI    string
	AppKey     string
}

// SagaLoRaWANDevice builds an OTAA LoRaWAN device riding a simulator UDP gateway
// (GatewayID), firing one LHT65N-style uplink event.
func SagaLoRaWANDevice(runID string, target LoRaTarget) *DeviceInputBuilder {
	config := mustJSON(map[string]any{
		"kind":       "lorawan",
		"gatewayId":  target.GatewayID,
		"region":     "EU868",
		"macVersion": "1.0.3",
		"activation": "otaa",
		"devEui":     target.DevEUI,
		"joinEui":    target.JoinEUI,
		"appKey":     target.AppKey,
	})
	return &DeviceInputBuilder{spec: devicescontract.DeviceInput{
		Name:       fmt.Sprintf("e2e LoRa UDP %s", runID),
		DeviceID:   fmt.Sprintf("e2e-lora-udp-%s", runID),
		ProtocolID: "lorawan",
		Enabled:    true,
		StoreLogs:  true,
		Config:     config,
		Attributes: map[string]string{},
		Events:     loraUplinkEvents(),
	}}
}

// SagaBasicStationDevice builds an OTAA LoRaWAN device carrying its own Basics
// Station WebSocket link (LNSURI + GatewayEUI), firing one uplink event.
func SagaBasicStationDevice(runID string, target LoRaTarget) *DeviceInputBuilder {
	config := mustJSON(map[string]any{
		"kind":       "basicstation",
		"lnsUri":     target.LNSURI,
		"gatewayEui": target.GatewayEUI,
		"region":     "EU868",
		"macVersion": "1.0.3",
		"activation": "otaa",
		"devEui":     target.DevEUI,
		"joinEui":    target.JoinEUI,
		"appKey":     target.AppKey,
	})
	return &DeviceInputBuilder{spec: devicescontract.DeviceInput{
		Name:       fmt.Sprintf("e2e LoRa BS %s", runID),
		DeviceID:   fmt.Sprintf("e2e-lora-bs-%s", runID),
		ProtocolID: "basicstation",
		Enabled:    true,
		StoreLogs:  true,
		Config:     config,
		Attributes: map[string]string{},
		Events:     loraUplinkEvents(),
	}}
}

// loraUplinkEvents is the one pre-registered LoRaWAN uplink both transports fire.
func loraUplinkEvents() json.RawMessage {
	return mustJSON([]any{
		map[string]any{
			"id":   "e1",
			"name": "LHT65N uplink",
			"lorawan": map[string]any{
				"fport":      2,
				"confirmed":  false,
				"payloadHex": "0BB809F6025D0000000000",
			},
		},
	})
}

// mustJSON marshals a fixture fragment, panicking on the impossible error so
// builders stay expression-friendly.
func mustJSON(v any) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}
