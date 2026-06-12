package services

import (
	"encoding/json"
	"testing"

	devicesDtos "simulator/service/src/modules/devices/application/dtos"
	gatewaysDtos "simulator/service/src/modules/gateways/application/dtos"
)

// The scheduler keeps a LoRaWAN device's jobs even when its gateway is offline: the
// fire still runs so the console can report the attempt with a gateway-offline
// status. The dispatcher -- not the scheduler -- is what skips the actual uplink
// when there is no live link.
func TestBuildDesired_KeepsLoRaWANJobWhenGatewayOffline(t *testing.T) {
	device := devicesDtos.Device{
		ID: "ld", Enabled: true, ProtocolID: "lorawan",
		Config: json.RawMessage(`{"gatewayId":"gw1","region":"EU868","macVersion":"1.0.3","activation":"otaa",` +
			`"devEui":"0011223344556677","joinEui":"0000000000000000","appKey":"00112233445566778899AABBCCDDEEFF"}`),
		Events: json.RawMessage(`[{"id":"e1","name":"u","lorawan":{"fport":2,"confirmed":false,"payloadHex":"00"},` +
			`"schedule":{"enabled":true,"every":10,"unit":"seconds"}}]`),
	}
	disabledGw := gatewaysDtos.Gateway{
		ID: "gw1", EUI: "0102030405060708", Enabled: false, Region: "EU868",
		Link: json.RawMessage(`{"protocol":"udp","host":"127.0.0.1","port":1700}`),
	}
	eng := newEngine([]devicesDtos.Device{device}, disabledGw)
	if n := len(eng.buildDesired()); n != 1 {
		t.Fatalf("scheduler must keep the job even with the gateway offline, got %d", n)
	}
}
