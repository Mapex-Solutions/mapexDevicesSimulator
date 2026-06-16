package steps

import (
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/core/saga"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/devices/payloads"
	targetSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/targets/steps"
)

// CreateMQTTReceiveDevice creates an enabled MQTT device that subscribes to a
// downlink topic on the broker (receiving on). Enabling it opens a session that
// subscribes; a retained publish to the topic is then delivered as a downlink.
//
// Reads (bag): targetSteps.BagKeyMQTTPlainURL, BagKeyMQTTUsername,
// BagKeyMQTTPassword.
// Writes (bag): BagKeyDeviceID, BagKeyDeviceDeviceID, BagKeyMQTTReceiveDeviceID.
// Compensate: DELETE the device.
//
// The downlink topic the publisher targets is derived from the run id via
// payloads.MQTTDownlinkTopic, the same source the device's subscription uses.
func CreateMQTTReceiveDevice() saga.Step {
	return saga.Step{
		Name: "devices.CreateMQTTReceiveDevice",
		Do: func(c *saga.Context) error {
			spec := payloads.SagaMQTTReceiveDevice(
				c.RunID,
				c.MustGetString(targetSteps.BagKeyMQTTPlainURL),
				c.MustGetString(targetSteps.BagKeyMQTTUsername),
				c.MustGetString(targetSteps.BagKeyMQTTPassword),
			).Build()
			return createDevice(c, spec, BagKeyMQTTReceiveDeviceID)
		},
		Compensate: func(c *saga.Context) error {
			return deleteDevice(c, BagKeyMQTTReceiveDeviceID)
		},
	}
}
