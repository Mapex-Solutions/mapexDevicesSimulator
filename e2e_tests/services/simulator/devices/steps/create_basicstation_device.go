package steps

import (
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/common/constants"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/core/saga"
	provSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/chirpstack/provisioning/steps"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/devices/payloads"
)

// CreateBasicStationDevice creates an OTAA LoRaWAN device that carries its own
// Basics Station WebSocket link to the gateway bridge (no separate simulator
// gateway), using the DevEUI/AppKey provisioned on ChirpStack.
//
// Reads (bag): provSteps.BagKeyBSGatewayEUI, provSteps.BagKeyBSDevEUI,
// provSteps.BagKeyBSAppKey.
// Writes (bag): BagKeyDeviceID, BagKeyDeviceDeviceID, BagKeyBasicStationDeviceID.
// Compensate: DELETE the device.
func CreateBasicStationDevice() saga.Step {
	return saga.Step{
		Name: "devices.CreateBasicStationDevice",
		Do: func(c *saga.Context) error {
			spec := payloads.SagaBasicStationDevice(c.RunID, payloads.LoRaTarget{
				GatewayEUI: c.MustGetString(provSteps.BagKeyBSGatewayEUI),
				LNSURI:     constants.ChirpStackBasicStationURL,
				DevEUI:     c.MustGetString(provSteps.BagKeyBSDevEUI),
				JoinEUI:    joinEUI,
				AppKey:     c.MustGetString(provSteps.BagKeyBSAppKey),
			}).Build()
			return createDevice(c, spec, BagKeyBasicStationDeviceID)
		},
		Compensate: func(c *saga.Context) error {
			return deleteDevice(c, BagKeyBasicStationDeviceID)
		},
	}
}
