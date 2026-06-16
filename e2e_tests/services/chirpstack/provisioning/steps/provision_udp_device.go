package steps

import (
	"context"
	"fmt"

	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/core/saga"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/chirpstack/provisioning/payloads"
	stackSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/chirpstack/stack/steps"
)

// ProvisionUDPDevice registers the Semtech UDP transport's gateway, device, and
// OTAA keys on ChirpStack so the simulator's join over UDP is accepted. It
// publishes the identity to the bag for the simulator-side steps and the join
// assert.
//
// Reads (bag): stackSteps.BagKeyClient, BagKeyTenantID, BagKeyApplicationID,
// BagKeyDeviceProfileID.
// Writes (bag): BagKeyUDPGatewayEUI, BagKeyUDPDevEUI, BagKeyUDPAppKey.
// Compensate: delete the device and gateway.
func ProvisionUDPDevice() saga.Step {
	return saga.Step{
		Name: "chirpstack.ProvisionUDPDevice",
		Do: func(c *saga.Context) error {
			cli := stackSteps.ClientFromBag(c)
			id := payloads.SagaUDPIdentity(c.RunID)
			tenantID := c.MustGetString(BagKeyTenantID)
			appID := c.MustGetString(BagKeyApplicationID)
			profileID := c.MustGetString(BagKeyDeviceProfileID)

			if err := cli.CreateGateway(c.Stdctx, tenantID, id.GatewayEUI, "e2e-udp-gw"); err != nil {
				return fmt.Errorf("provision udp gateway: %w", err)
			}
			c.Set(BagKeyUDPGatewayEUI, id.GatewayEUI)

			if err := cli.CreateDevice(c.Stdctx, appID, profileID, id.DevEUI, id.JoinEUI, "e2e-udp-dev"); err != nil {
				return fmt.Errorf("provision udp device: %w", err)
			}
			c.Set(BagKeyUDPDevEUI, id.DevEUI)
			if err := cli.SetDeviceKeys(c.Stdctx, id.DevEUI, id.AppKey); err != nil {
				return fmt.Errorf("set udp device keys: %w", err)
			}
			c.Set(BagKeyUDPAppKey, id.AppKey)
			c.Set(BagKeyActiveDevEUI, id.DevEUI)
			return nil
		},
		Compensate: func(c *saga.Context) error {
			cli := stackSteps.ClientFromBag(c)
			if devEUI, ok := c.Get(BagKeyUDPDevEUI); ok {
				_ = cli.DeleteDevice(context.Background(), devEUI.(string))
			}
			if gwEUI, ok := c.Get(BagKeyUDPGatewayEUI); ok {
				_ = cli.DeleteGateway(context.Background(), gwEUI.(string))
			}
			return nil
		},
	}
}
