package steps

import (
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/core/saga"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/devices/payloads"
	gatewaySteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/gateways/steps"
	provSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/chirpstack/provisioning/steps"
)

// CreateLoRaWANDevice creates an OTAA LoRaWAN device riding the simulator's UDP
// gateway, using the DevEUI/AppKey provisioned on ChirpStack so the join keys
// match. Enabling the device starts the join; the join assert confirms it.
//
// Reads (bag): gatewaySteps.BagKeyUDPSimGatewayID, provSteps.BagKeyUDPDevEUI,
// provSteps.BagKeyUDPAppKey.
// Writes (bag): BagKeyDeviceID, BagKeyDeviceDeviceID, BagKeyLoRaWANDeviceID.
// Compensate: DELETE the device.
func CreateLoRaWANDevice() saga.Step {
	return saga.Step{
		Name: "devices.CreateLoRaWANDevice",
		Do: func(c *saga.Context) error {
			spec := payloads.SagaLoRaWANDevice(c.RunID, payloads.LoRaTarget{
				GatewayID: c.MustGetString(gatewaySteps.BagKeyUDPSimGatewayID),
				DevEUI:    c.MustGetString(provSteps.BagKeyUDPDevEUI),
				JoinEUI:   joinEUI,
				AppKey:    c.MustGetString(provSteps.BagKeyUDPAppKey),
			}).Build()
			return createDevice(c, spec, BagKeyLoRaWANDeviceID)
		},
		Compensate: func(c *saga.Context) error {
			return deleteDevice(c, BagKeyLoRaWANDeviceID)
		},
	}
}
