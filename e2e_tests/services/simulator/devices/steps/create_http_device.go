package steps

import (
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/core/saga"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/devices/payloads"
	targetSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/targets/steps"
)

// CreateHTTPDevice creates a send-only HTTP device targeting the in-process echo
// started by StartEcho, with one templated-body event and storeLogs on.
//
// Reads (bag): targetSteps.BagKeyEchoURL.
// Writes (bag): BagKeyDeviceID, BagKeyDeviceDeviceID, BagKeyHTTPDeviceID.
// Compensate: DELETE the device.
func CreateHTTPDevice() saga.Step {
	return saga.Step{
		Name: "devices.CreateHTTPDevice",
		Do: func(c *saga.Context) error {
			spec := payloads.SagaHTTPDevice(c.RunID, c.MustGetString(targetSteps.BagKeyEchoURL)).Build()
			return createDevice(c, spec, BagKeyHTTPDeviceID)
		},
		Compensate: func(c *saga.Context) error {
			return deleteDevice(c, BagKeyHTTPDeviceID)
		},
	}
}
