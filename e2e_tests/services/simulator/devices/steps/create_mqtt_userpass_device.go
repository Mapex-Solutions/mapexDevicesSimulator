package steps

import (
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/core/saga"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/devices/payloads"
	targetSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/targets/steps"
)

// CreateMQTTUserPassDevice creates an MQTT device that authenticates to the
// broker with a username and password over a plain tcp:// connection.
//
// Reads (bag): targetSteps.BagKeyMQTTPlainURL, BagKeyMQTTUsername,
// BagKeyMQTTPassword.
// Writes (bag): BagKeyDeviceID, BagKeyDeviceDeviceID, BagKeyMQTTUserPassDeviceID.
// Compensate: DELETE the device.
func CreateMQTTUserPassDevice() saga.Step {
	return saga.Step{
		Name: "devices.CreateMQTTUserPassDevice",
		Do: func(c *saga.Context) error {
			spec := payloads.SagaMQTTDeviceUserPass(c.RunID, payloads.MQTTBrokerTarget{
				BrokerURL: c.MustGetString(targetSteps.BagKeyMQTTPlainURL),
				Username:  c.MustGetString(targetSteps.BagKeyMQTTUsername),
				Password:  c.MustGetString(targetSteps.BagKeyMQTTPassword),
			}).Build()
			return createDevice(c, spec, BagKeyMQTTUserPassDeviceID)
		},
		Compensate: func(c *saga.Context) error {
			return deleteDevice(c, BagKeyMQTTUserPassDeviceID)
		},
	}
}
