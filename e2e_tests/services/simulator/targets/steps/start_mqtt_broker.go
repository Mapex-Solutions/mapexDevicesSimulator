package steps

import (
	"fmt"

	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/common/utils"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/core/saga"
)

// StartMQTTBroker boots the in-process MQTT broker (plain username/password
// listener + mutual-TLS listener) and publishes its coordinates and a fresh CA /
// client certificate. Compensate shuts it down.
//
// Writes (bag): BagKeyMQTTBroker (*utils.MQTTBroker), BagKeyMQTTPlainURL,
// BagKeyMQTTTLSURL, BagKeyMQTTUsername, BagKeyMQTTPassword, BagKeyMQTTCAPem,
// BagKeyMQTTClientCertPem, BagKeyMQTTClientKeyPem.
func StartMQTTBroker() saga.Step {
	return saga.Step{
		Name: "targets.StartMQTTBroker",
		Do: func(c *saga.Context) error {
			broker, err := utils.StartMQTTBroker(utils.MQTTCreds{Username: brokerUser, Password: brokerPass})
			if err != nil {
				return fmt.Errorf("start mqtt broker: %w", err)
			}
			c.Set(BagKeyMQTTBroker, broker)
			c.Set(BagKeyMQTTPlainURL, broker.PlainURL())
			c.Set(BagKeyMQTTTLSURL, broker.TLSURL())
			c.Set(BagKeyMQTTUsername, brokerUser)
			c.Set(BagKeyMQTTPassword, brokerPass)
			c.Set(BagKeyMQTTCAPem, broker.CAPEM())
			c.Set(BagKeyMQTTClientCertPem, broker.ClientCertPEM())
			c.Set(BagKeyMQTTClientKeyPem, broker.ClientKeyPEM())
			return nil
		},
		Compensate: func(c *saga.Context) error {
			if v, ok := c.Get(BagKeyMQTTBroker); ok {
				if broker, ok := v.(*utils.MQTTBroker); ok {
					broker.Close()
				}
			}
			return nil
		},
	}
}
