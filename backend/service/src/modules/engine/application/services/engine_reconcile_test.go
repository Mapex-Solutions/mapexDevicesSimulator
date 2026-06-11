package services

import (
	"encoding/json"
	"testing"

	devicesDtos "simulator/service/src/modules/devices/application/dtos"
	gatewaysDtos "simulator/service/src/modules/gateways/application/dtos"
)

// A scheduled LoRaWAN or Basics Station event must not be scheduled while its
// gateway is disabled: the job heap has to honor the gateway's enabled flag the
// same way the session manager does, otherwise the device keeps firing frames
// with no live link behind it.
func TestBuildDesired_LoRaWANRespectsGatewayFlag(t *testing.T) {
	schedEvent := `[{"id":"e1","name":"u","lorawan":{"fport":2,"confirmed":false,"payloadHex":"00"},` +
		`"schedule":{"enabled":true,"every":10,"unit":"seconds"}}]`
	udpDevice := devicesDtos.Device{
		ID: "ld", Enabled: true, ProtocolID: "lorawan",
		Config: json.RawMessage(`{"gatewayId":"gw1","region":"EU868","macVersion":"1.0.3","activation":"otaa",` +
			`"devEui":"0011223344556677","joinEui":"0000000000000000","appKey":"00112233445566778899AABBCCDDEEFF"}`),
		Events: json.RawMessage(schedEvent),
	}
	bsDevice := devicesDtos.Device{
		ID: "bd", Enabled: true, ProtocolID: "basicstation",
		Config: json.RawMessage(`{"gatewayEui":"0102030405060708","lnsUri":"ws://127.0.0.1:3001","region":"EU868",` +
			`"macVersion":"1.0.3","activation":"otaa","devEui":"0011223344556677","joinEui":"0000000000000000",` +
			`"appKey":"00112233445566778899AABBCCDDEEFF"}`),
		Events: json.RawMessage(schedEvent),
	}
	gw := func(enabled bool) gatewaysDtos.Gateway {
		return gatewaysDtos.Gateway{
			ID: "gw1", EUI: "0102030405060708", Enabled: enabled, Region: "EU868",
			Link: json.RawMessage(`{"protocol":"udp","host":"127.0.0.1","port":1700}`),
		}
	}

	tests := []struct {
		name     string
		device   devicesDtos.Device
		gateway  gatewaysDtos.Gateway
		wantJobs int
	}{
		{"udp gateway disabled drops the job", udpDevice, gw(false), 0},
		{"udp gateway enabled keeps the job", udpDevice, gw(true), 1},
		{"basics station gateway disabled drops the job", bsDevice, gw(false), 0},
		{"basics station gateway enabled keeps the job", bsDevice, gw(true), 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eng := newEngine([]devicesDtos.Device{tt.device}, tt.gateway)
			if n := len(eng.buildDesired()); n != tt.wantJobs {
				t.Fatalf("buildDesired jobs = %d, want %d", n, tt.wantJobs)
			}
		})
	}
}
