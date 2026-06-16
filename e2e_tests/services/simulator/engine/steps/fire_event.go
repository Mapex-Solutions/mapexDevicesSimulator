package steps

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/core/saga"
	deviceSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/devices/steps"
)

// FireTelemetry fires the canonical pre-registered event on the most recently
// created device and records the fire time. It dispatches one uplink through
// the engine.
//
// Reads (bag): deviceSteps.BagKeyDeviceID.
// Writes (bag): BagKeyFiredAt.
func FireTelemetry() saga.Step {
	return saga.Step{
		Name: "engine.FireTelemetry",
		Do: func(c *saga.Context) error {
			id := c.MustGetString(deviceSteps.BagKeyDeviceID)
			body := map[string]string{"eventId": fireEventID}
			resp, err := c.Clients.Sim.Raw(c.Stdctx, http.MethodPost, "/api/devices/"+id+"/fire", body)
			if err != nil {
				return fmt.Errorf("fire: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				return fmt.Errorf("fire: unexpected status %d", resp.StatusCode)
			}
			c.Set(BagKeyFiredAt, time.Now().UTC())
			return nil
		},
	}
}
