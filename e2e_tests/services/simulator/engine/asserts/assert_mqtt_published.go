// Package asserts holds saga oracles for the simulator's engine send path. The
// MQTT oracle reads the in-process broker the device published to — the
// definitive proof the device connected, authenticated, and delivered.
package asserts

import (
	"fmt"
	"strings"
	"time"

	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/core/saga"
	deviceSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/devices/steps"
	targetSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/targets/steps"
)

// AssertMQTTPublished polls the broker until it has accepted a publish carrying
// the most recently created device's deviceId in its payload. Because the broker
// only records a publish after the CONNECT passed authentication, this single
// check proves the whole chain: the device reached the broker, satisfied its
// auth mode, and delivered the uplink.
//
// Reads (bag): targetSteps.BagKeyMQTTBroker, deviceSteps.BagKeyDeviceDeviceID.
func AssertMQTTPublished() saga.Assert {
	return saga.Assert{
		Name: "engine.AssertMQTTPublished",
		Check: func(c *saga.Context) error {
			broker := targetSteps.BrokerFromBag(c)
			deviceID := c.MustGetString(deviceSteps.BagKeyDeviceDeviceID)
			deadline := time.Now().Add(10 * time.Second)
			for {
				for _, p := range broker.Published() {
					if strings.Contains(p.Payload, deviceID) {
						return nil
					}
				}
				if time.Now().After(deadline) {
					return fmt.Errorf("broker received no publish carrying deviceId %q before timeout", deviceID)
				}
				time.Sleep(200 * time.Millisecond)
			}
		},
	}
}
