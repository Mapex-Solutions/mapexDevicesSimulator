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
	deviceSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/devices/steps"
)

// AssertDownlinkLogged polls GET /api/logs filtered by the device and asserts a
// down/downlink frame landed carrying the published payload — proof the device's
// subscription received the message and the engine surfaced it on the inbound
// path. The payload carries the run id, so the match is unambiguous.
//
// Reads (bag): deviceSteps.BagKeyDeviceDeviceID.
func AssertDownlinkLogged() saga.Assert {
	return saga.Assert{
		Name: "logs.AssertDownlinkLogged",
		Check: func(c *saga.Context) error {
			deviceID := c.MustGetString(deviceSteps.BagKeyDeviceDeviceID)
			deadline := time.Now().Add(10 * time.Second)
			seen := 0
			for {
				q := url.Values{}
				q.Set("limit", "20")
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
					if l.Direction == "down" && l.Kind == "downlink" && strings.Contains(l.Payload, c.RunID) {
						return nil
					}
				}
				if time.Now().After(deadline) {
					return fmt.Errorf("no downlink frame logged for device %q after timeout (%d logs seen)", deviceID, seen)
				}
				time.Sleep(300 * time.Millisecond)
			}
		},
	}
}
