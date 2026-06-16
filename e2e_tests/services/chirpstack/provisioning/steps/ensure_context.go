package steps

import (
	"fmt"

	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/core/saga"
	stackSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/chirpstack/stack/steps"
)

// EnsureApplicationContext makes sure the LNS has a tenant, an application, and
// an EU868 / LoRaWAN 1.0.3 OTAA device profile to provision devices under, and
// publishes their ids. Compensate is a no-op: the stack is torn down with its
// volumes at the end, which removes everything this created.
//
// Reads (bag): stackSteps.BagKeyClient.
// Writes (bag): BagKeyTenantID, BagKeyApplicationID, BagKeyDeviceProfileID.
func EnsureApplicationContext() saga.Step {
	return saga.Step{
		Name: "chirpstack.EnsureApplicationContext",
		Do: func(c *saga.Context) error {
			cli := stackSteps.ClientFromBag(c)

			tenantID, err := cli.EnsureTenant(c.Stdctx, "e2e")
			if err != nil {
				return fmt.Errorf("ensure tenant: %w", err)
			}
			c.Set(BagKeyTenantID, tenantID)

			appID, err := cli.CreateApplication(c.Stdctx, tenantID, "e2e-"+c.RunID)
			if err != nil {
				return fmt.Errorf("create application: %w", err)
			}
			c.Set(BagKeyApplicationID, appID)

			profileID, err := cli.CreateDeviceProfile(c.Stdctx, tenantID, "e2e-eu868-otaa-"+c.RunID)
			if err != nil {
				return fmt.Errorf("create device profile: %w", err)
			}
			c.Set(BagKeyDeviceProfileID, profileID)
			return nil
		},
	}
}
