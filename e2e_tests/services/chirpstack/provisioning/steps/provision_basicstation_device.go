package steps

import (
	"context"
	"fmt"

	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/core/saga"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/chirpstack/provisioning/payloads"
	stackSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/chirpstack/stack/steps"
)

// ProvisionBasicStationDevice registers the Basics Station transport's gateway,
// device, and OTAA keys on ChirpStack so the simulator's join over the Basics
// Station WebSocket is accepted.
//
// Reads (bag): stackSteps.BagKeyClient, BagKeyTenantID, BagKeyApplicationID,
// BagKeyDeviceProfileID.
// Writes (bag): BagKeyBSGatewayEUI, BagKeyBSDevEUI, BagKeyBSAppKey.
// Compensate: delete the device and gateway.
func ProvisionBasicStationDevice() saga.Step {
	return saga.Step{
		Name: "chirpstack.ProvisionBasicStationDevice",
		Do: func(c *saga.Context) error {
			cli := stackSteps.ClientFromBag(c)
			id := payloads.SagaBasicStationIdentity(c.RunID)
			tenantID := c.MustGetString(BagKeyTenantID)
			appID := c.MustGetString(BagKeyApplicationID)
			profileID := c.MustGetString(BagKeyDeviceProfileID)

			if err := cli.CreateGateway(c.Stdctx, tenantID, id.GatewayEUI, "e2e-bs-gw"); err != nil {
				return fmt.Errorf("provision basicstation gateway: %w", err)
			}
			c.Set(BagKeyBSGatewayEUI, id.GatewayEUI)

			if err := cli.CreateDevice(c.Stdctx, appID, profileID, id.DevEUI, id.JoinEUI, "e2e-bs-dev"); err != nil {
				return fmt.Errorf("provision basicstation device: %w", err)
			}
			c.Set(BagKeyBSDevEUI, id.DevEUI)
			if err := cli.SetDeviceKeys(c.Stdctx, id.DevEUI, id.AppKey); err != nil {
				return fmt.Errorf("set basicstation device keys: %w", err)
			}
			c.Set(BagKeyBSAppKey, id.AppKey)
			c.Set(BagKeyActiveDevEUI, id.DevEUI)
			return nil
		},
		Compensate: func(c *saga.Context) error {
			cli := stackSteps.ClientFromBag(c)
			if devEUI, ok := c.Get(BagKeyBSDevEUI); ok {
				_ = cli.DeleteDevice(context.Background(), devEUI.(string))
			}
			if gwEUI, ok := c.Get(BagKeyBSGatewayEUI); ok {
				_ = cli.DeleteGateway(context.Background(), gwEUI.(string))
			}
			return nil
		},
	}
}
