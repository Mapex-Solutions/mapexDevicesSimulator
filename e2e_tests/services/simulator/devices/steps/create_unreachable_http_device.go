package steps

import (
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/core/saga"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/devices/payloads"
)

// unreachableHTTPURL is a loopback port nothing listens on, so a fire against it
// fails fast with a connection error — the engine's send-error path under test.
const unreachableHTTPURL = "http://127.0.0.1:1/ingest"

// CreateUnreachableHTTPDevice creates an enabled HTTP device whose target is an
// unreachable address, so firing it produces a send error the engine reports.
//
// Writes (bag): BagKeyDeviceID, BagKeyDeviceDeviceID, BagKeyHTTPDeviceID.
// Compensate: DELETE the device.
func CreateUnreachableHTTPDevice() saga.Step {
	return saga.Step{
		Name: "devices.CreateUnreachableHTTPDevice",
		Do: func(c *saga.Context) error {
			spec := payloads.SagaHTTPDevice(c.RunID, unreachableHTTPURL).Build()
			return createDevice(c, spec, BagKeyHTTPDeviceID)
		},
		Compensate: func(c *saga.Context) error {
			return deleteDevice(c, BagKeyHTTPDeviceID)
		},
	}
}
