package steps

import (
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/core/saga"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/devices/payloads"
	targetSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/targets/steps"
)

// CreateMQTTTLSDevice creates an MQTT device that authenticates to the broker
// with a client certificate over an ssl:// connection (mutual TLS).
//
// Reads (bag): targetSteps.BagKeyMQTTTLSURL, BagKeyMQTTClientCertPem,
// BagKeyMQTTClientKeyPem, BagKeyMQTTCAPem.
// Writes (bag): BagKeyDeviceID, BagKeyDeviceDeviceID, BagKeyMQTTTLSDeviceID.
// Compensate: DELETE the device.
func CreateMQTTTLSDevice() saga.Step {
	return saga.Step{
		Name: "devices.CreateMQTTTLSDevice",
		Do: func(c *saga.Context) error {
			spec := payloads.SagaMQTTDeviceTLS(c.RunID, payloads.MQTTBrokerTarget{
				BrokerURL:  c.MustGetString(targetSteps.BagKeyMQTTTLSURL),
				TLSCertPem: c.MustGetString(targetSteps.BagKeyMQTTClientCertPem),
				TLSKeyPem:  c.MustGetString(targetSteps.BagKeyMQTTClientKeyPem),
				TLSCaPem:   c.MustGetString(targetSteps.BagKeyMQTTCAPem),
			}).Build()
			return createDevice(c, spec, BagKeyMQTTTLSDeviceID)
		},
		Compensate: func(c *saga.Context) error {
			return deleteDevice(c, BagKeyMQTTTLSDeviceID)
		},
	}
}
