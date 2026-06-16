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

// AssertFireErrorLogged polls GET /api/logs filtered by the device and asserts a
// frame landed with status "error" and a non-empty response — proof a failed
// send is surfaced to the user (and persisted), not silently dropped.
//
// Reads (bag): deviceSteps.BagKeyDeviceDeviceID.
func AssertFireErrorLogged() saga.Assert {
	return saga.Assert{
		Name: "logs.AssertFireErrorLogged",
		Check: func(c *saga.Context) error {
			deviceID := c.MustGetString(deviceSteps.BagKeyDeviceDeviceID)
			deadline := time.Now().Add(15 * time.Second)
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
					if l.Status == "error" && l.Response != "" {
						return nil
					}
				}
				if time.Now().After(deadline) {
					return fmt.Errorf("no error frame logged for device %q after timeout (%d logs seen)", deviceID, seen)
				}
				time.Sleep(300 * time.Millisecond)
			}
		},
	}
}
