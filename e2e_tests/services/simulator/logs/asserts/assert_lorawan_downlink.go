package asserts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	logscontract "simulator/packages/contracts/logs"

	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/common/types"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/core/saga"
	provSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/chirpstack/provisioning/steps"
	deviceSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/devices/steps"
)

// lorawanDownlinkTimeout bounds how long the queued downlink may take to arrive
// after the triggering uplink (RX window plus the poll).
const lorawanDownlinkTimeout = 20 * time.Second

// AssertLoRaWANDownlinkReceived polls GET /api/logs filtered by the device and
// asserts a down/downlink frame landed whose payload is the hex of the bytes
// enqueued on ChirpStack — proof the simulator received the LNS downlink in its
// RX window and surfaced it on the inbound path.
//
// Reads (bag): deviceSteps.BagKeyDeviceDeviceID, provSteps.BagKeyDownlinkHex.
func AssertLoRaWANDownlinkReceived() saga.Assert {
	return saga.Assert{
		Name: "logs.AssertLoRaWANDownlinkReceived",
		Check: func(c *saga.Context) error {
			deviceID := c.MustGetString(deviceSteps.BagKeyDeviceDeviceID)
			wantHex := c.MustGetString(provSteps.BagKeyDownlinkHex)
			deadline := time.Now().Add(lorawanDownlinkTimeout)
			seen := 0
			for {
				q := url.Values{}
				q.Set("limit", "30")
				q.Set("device", deviceID)
				resp, err := c.Clients.Sim.Raw(c.Stdctx, http.MethodGet, "/api/logs?"+q.Encode(), nil)
				if err != nil {
					return fmt.Errorf("list logs: %w", err)
				}
				var env types.Envelope
				decErr := json.NewDecoder(resp.Body).Decode(&env)
				resp.Body.Close()
				if decErr != nil {
					return fmt.Errorf("decode logs envelope: %w", decErr)
				}
				var page logscontract.LogPage
				if err := json.Unmarshal(env.Data, &page); err != nil {
					return fmt.Errorf("decode log page: %w", err)
				}
				seen = len(page.Items)
				for i := range page.Items {
					l := page.Items[i]
					if l.Direction == "down" && l.Kind == "downlink" && strings.Contains(strings.ToLower(l.Payload), wantHex) {
						return nil
					}
				}
				if time.Now().After(deadline) {
					return fmt.Errorf("no downlink frame with payload %q for device %q after timeout (%d logs seen)", wantHex, deviceID, seen)
				}
				time.Sleep(500 * time.Millisecond)
			}
		},
	}
}
