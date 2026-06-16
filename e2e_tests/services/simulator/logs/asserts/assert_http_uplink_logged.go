// Package asserts holds saga oracles for the simulator's logs module. They read
// the simulator's public GET /api/logs — never the SQLite file — so the journey
// validates the platform from the outside.
package asserts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	logscontract "simulator/packages/contracts/logs"

	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/common/types"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/core/saga"
	deviceSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/devices/steps"
)

// AssertHTTPUplinkLogged polls GET /api/logs filtered by the device and asserts
// a data uplink landed with status "200" and a non-empty response — proof the
// fire reached the echo target and the simulator persisted both the payload and
// the reply. A 200 can only come from the live echo, so it doubles as the
// round-trip check.
//
// Reads (bag): deviceSteps.BagKeyDeviceDeviceID.
func AssertHTTPUplinkLogged() saga.Assert {
	return saga.Assert{
		Name: "logs.AssertHTTPUplinkLogged",
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
					if l.Direction == "up" && l.Kind == "data" && l.Status == "200" && l.Response != "" {
						return nil
					}
				}
				if time.Now().After(deadline) {
					return fmt.Errorf("no 200 uplink with a response logged for device %q after timeout (%d logs seen)", deviceID, seen)
				}
				time.Sleep(300 * time.Millisecond)
			}
		},
	}
}
