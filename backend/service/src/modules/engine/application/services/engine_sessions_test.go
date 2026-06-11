package services

import (
	"encoding/json"
	"testing"

	devicesDtos "simulator/service/src/modules/devices/application/dtos"
	"simulator/service/src/modules/engine/application/di"
	dispatch "simulator/service/src/modules/engine/infrastructure/dispatch"
	session "simulator/service/src/modules/engine/infrastructure/session"
	gatewaysDtos "simulator/service/src/modules/gateways/application/dtos"
	"simulator/service/src/shared/reconcile"
)

// mqttDevice builds an MQTT device with the given receive config.
func mqttDevice(id string, enabled, receive bool, subs string) devicesDtos.Device {
	cfg := `{"brokerUrl":"tcp://localhost:1883","clientId":"c1","baseTopic":"sim",` +
		`"receiveEnabled":` + boolStr(receive) + `,"subscriptions":` + subs + `}`
	return devicesDtos.Device{
		ID:         id,
		Enabled:    enabled,
		Name:       "MQTT " + id,
		DeviceID:   "dev-" + id,
		ProtocolID: "mqtt",
		Config:     json.RawMessage(cfg),
		Events:     json.RawMessage(`[]`),
	}
}

func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func newEngine(list []devicesDtos.Device, gateways ...gatewaysDtos.Gateway) *EngineService {
	return New(di.EngineServiceDI{
		Devices:    &fakeDevices{list: list},
		Gateways:   &fakeGateways{list: gateways},
		Logs:       &fakeLogWriter{},
		Console:    &fakePublisher{},
		Registry:   dispatch.NewRegistry(),
		Connectors: session.NewConnectorRegistry(),
		Reconcile:  reconcile.New(),
	}).(*EngineService)
}

func TestBuildDesiredSessions_MQTTReceiveOn(t *testing.T) {
	subs := `[{"name":"cmd","topic":"down","qos":1}]`
	eng := newEngine([]devicesDtos.Device{mqttDevice("a", true, true, subs)})

	desired := eng.buildDesiredSessions()
	d, ok := desired["a"]
	if !ok {
		t.Fatal("expected a session for the enabled mqtt device")
	}
	if len(d.spec.Subscriptions) != 1 {
		t.Fatalf("want 1 subscription, got %d", len(d.spec.Subscriptions))
	}
	// baseTopic ("sim") must prefix the subscription topic ("down").
	if got := d.spec.Subscriptions[0].Topic; got != "sim/down" {
		t.Fatalf("want topic sim/down, got %q", got)
	}
	if d.spec.Subscriptions[0].QoS != 1 {
		t.Fatalf("want qos 1, got %d", d.spec.Subscriptions[0].QoS)
	}
}

func TestBuildDesiredSessions_ReceiveOffHasNoSubscriptions(t *testing.T) {
	subs := `[{"name":"cmd","topic":"down","qos":1}]`
	eng := newEngine([]devicesDtos.Device{mqttDevice("a", true, false, subs)})

	desired := eng.buildDesiredSessions()
	d, ok := desired["a"]
	if !ok {
		t.Fatal("an enabled mqtt device still opens a session for uplinks")
	}
	if len(d.spec.Subscriptions) != 0 {
		t.Fatalf("receive off must yield no subscriptions, got %d", len(d.spec.Subscriptions))
	}
}

func TestBuildDesiredSessions_LoRaWANResolvesGatewayLink(t *testing.T) {
	dev := devicesDtos.Device{
		ID: "ld", Enabled: true, Name: "Field", DeviceID: "dev-lora", ProtocolID: "lorawan",
		Config: json.RawMessage(`{"gatewayId":"gw1","region":"EU868","macVersion":"1.0.3",` +
			`"activation":"otaa","devEui":"0011223344556677","joinEui":"70B3D57ED0000001",` +
			`"appKey":"00112233445566778899AABBCCDDEEFF"}`),
		Events: json.RawMessage(`[]`),
	}
	gw := gatewaysDtos.Gateway{
		ID: "gw1", EUI: "0016C001F1500001", Enabled: true, Region: "EU868",
		Link: json.RawMessage(`{"protocol":"udp","host":"127.0.0.1","port":1700}`),
	}
	eng := newEngine([]devicesDtos.Device{dev}, gw)

	desired := eng.buildDesiredSessions()
	d, ok := desired["ld"]
	if !ok || d.spec.LoRaWAN == nil {
		t.Fatal("expected a lorawan session spec")
	}
	if d.spec.LoRaWAN.LinkProtocol != "udp" || d.spec.LoRaWAN.LinkUDPPort != 1700 {
		t.Fatalf("link not resolved: %+v", d.spec.LoRaWAN)
	}
	if d.spec.LoRaWAN.GatewayEUI != "0016C001F1500001" {
		t.Fatalf("gateway EUI not carried: %s", d.spec.LoRaWAN.GatewayEUI)
	}
}

func TestBuildDesiredSessions_LoRaWANSkipsWhenGatewayMissing(t *testing.T) {
	dev := devicesDtos.Device{
		ID: "ld", Enabled: true, ProtocolID: "lorawan",
		Config: json.RawMessage(`{"gatewayId":"absent","region":"EU868","activation":"abp",` +
			`"appKey":"00112233445566778899AABBCCDDEEFF"}`),
		Events: json.RawMessage(`[]`),
	}
	eng := newEngine([]devicesDtos.Device{dev}) // no gateways
	if len(eng.buildDesiredSessions()) != 0 {
		t.Fatal("a lorawan device with a missing gateway must not open a session")
	}
}

func TestBuildDesiredSessions_DisabledAndHTTPExcluded(t *testing.T) {
	eng := newEngine([]devicesDtos.Device{
		mqttDevice("a", false, true, `[]`), // disabled
		{ID: "b", Enabled: true, ProtocolID: "http", Config: json.RawMessage(`{"url":"http://x"}`), Events: json.RawMessage(`[]`)}, // http: no connector
	})

	desired := eng.buildDesiredSessions()
	if len(desired) != 0 {
		t.Fatalf("disabled + http devices must not open sessions, got %d", len(desired))
	}
}

// bsDevice builds an enabled Basics Station device announcing the given gateway EUI.
func bsDevice(id, gatewayEUI string) devicesDtos.Device {
	cfg := `{"gatewayEui":"` + gatewayEUI + `","lnsUri":"ws://127.0.0.1:3001","region":"EU868",` +
		`"macVersion":"1.0.3","activation":"otaa","devEui":"0011223344556677",` +
		`"joinEui":"0000000000000000","appKey":"00112233445566778899AABBCCDDEEFF"}`
	return devicesDtos.Device{
		ID: id, Enabled: true, Name: "BS " + id, DeviceID: "dev-" + id,
		ProtocolID: "basicstation", Config: json.RawMessage(cfg), Events: json.RawMessage(`[]`),
	}
}

// bsGateway builds a Basics Station gateway entity with the given EUI / enabled state.
func bsGateway(eui string, enabled bool) gatewaysDtos.Gateway {
	return gatewaysDtos.Gateway{
		ID: "gw-" + eui, EUI: eui, Enabled: enabled, Region: "EU868",
		Link: json.RawMessage(`{"protocol":"basicstation","lnsUri":"ws://127.0.0.1:3001"}`),
	}
}

// A Basics Station device carries its own link, but if a gateway entity shares its
// EUI, the gateway's enabled flag gates the session — so disabling the gateway takes
// the device offline, the same as the UDP path.
func TestBuildDesiredSessions_BasicsStationRespectsGatewayFlag(t *testing.T) {
	const eui = "00AABBCCDDEEFF11" // device announces this (upper case)
	tests := []struct {
		name      string
		gateways  []gatewaysDtos.Gateway
		wantBuilt bool
	}{
		// lower-case gateway EUI proves the match is case-insensitive.
		{"matching gateway disabled drops the session", []gatewaysDtos.Gateway{bsGateway("00aabbccddeeff11", false)}, false},
		{"matching gateway enabled keeps the session", []gatewaysDtos.Gateway{bsGateway("00aabbccddeeff11", true)}, true},
		{"no matching gateway runs standalone", nil, true},
		{"a different disabled gateway is ignored", []gatewaysDtos.Gateway{bsGateway("1122334455667788", false)}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eng := newEngine([]devicesDtos.Device{bsDevice("bs", eui)}, tt.gateways...)
			_, ok := eng.buildDesiredSessions()["bs"]
			if ok != tt.wantBuilt {
				t.Fatalf("session built = %v, want %v", ok, tt.wantBuilt)
			}
		})
	}
}
