package steps

import (
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/core/saga"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/devices/payloads"
)

// unreachableMQTTBrokerURL is a loopback port nothing listens on, so the engine's
// session for this device never connects and cycles through the connecting /
// reconnecting status the test observes on the console stream.
const unreachableMQTTBrokerURL = "tcp://127.0.0.1:1"

// CreateUnreachableMQTTDevice creates an enabled MQTT device pointing at an
// unreachable broker. Enabling it opens a session that fails to connect, so the
// engine emits connecting / reconnecting status frames on the console.
//
// Writes (bag): BagKeyDeviceID, BagKeyDeviceDeviceID, BagKeyMQTTUserPassDeviceID.
// Compensate: DELETE the device.
func CreateUnreachableMQTTDevice() saga.Step {
	return saga.Step{
		Name: "devices.CreateUnreachableMQTTDevice",
		Do: func(c *saga.Context) error {
			spec := payloads.SagaMQTTDeviceUserPass(c.RunID, payloads.MQTTBrokerTarget{
				BrokerURL: unreachableMQTTBrokerURL,
				Username:  "e2e",
				Password:  "e2e",
			}).Build()
			return createDevice(c, spec, BagKeyMQTTUserPassDeviceID)
		},
		Compensate: func(c *saga.Context) error {
			return deleteDevice(c, BagKeyMQTTUserPassDeviceID)
		},
	}
}
