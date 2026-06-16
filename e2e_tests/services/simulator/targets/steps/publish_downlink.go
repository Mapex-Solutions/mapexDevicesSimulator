package steps

import (
	"fmt"

	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/core/saga"
	devicePayloads "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/devices/payloads"
)

// PublishDownlink injects a retained message on the MQTT receive device's
// subscribed topic, as an external party would, so the device receives it as a
// downlink. Retained delivery means it lands even if the device's subscription
// is not active yet (the session subscribes asynchronously), removing the race.
//
// Reads (bag): BagKeyMQTTBroker.
// Writes (bag): BagKeyDownlinkPayload.
func PublishDownlink() saga.Step {
	return saga.Step{
		Name: "targets.PublishDownlink",
		Do: func(c *saga.Context) error {
			broker := BrokerFromBag(c)
			topic := devicePayloads.MQTTDownlinkTopic(c.RunID)
			payload := fmt.Sprintf(`{"cmd":"open","run":%q}`, c.RunID)
			if err := broker.Publish(topic, []byte(payload), true, 1); err != nil {
				return fmt.Errorf("publish downlink: %w", err)
			}
			c.Set(BagKeyDownlinkPayload, payload)
			return nil
		},
	}
}
