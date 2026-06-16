// Package mqtt_downlink exercises the MQTT INBOUND path end to end against the
// live simulator and an in-process broker: a device with receiving enabled
// subscribes to a topic, an external publish arrives on it, and the simulator
// surfaces it as a downlink in the logs.
//
// Outcome on PASS:
//   - An enabled MQTT device subscribes to its downlink topic on the broker.
//   - A retained publish on that topic is delivered to the device.
//   - GET /api/logs shows a down/downlink frame carrying the payload.
//   - The device and broker are torn down (compensation).
//
// Outcome on FAIL:
//   - The failing step / assert name surfaces in the saga log
//     (devices.CreateMQTTReceiveDevice, targets.PublishDownlink,
//     logs.AssertDownlinkLogged).
package mqtt_downlink

import (
	"context"
	"testing"

	"github.com/Mapex-Solutions/mapexGoKit/utils/random"

	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/common/constants"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/common/utils"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/core/saga"
	deviceSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/devices/steps"
	logAsserts "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/logs/asserts"
	targetSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/targets/steps"
)

// Items is the ordered saga the journey runs.
//
//  1. StartMQTTBroker            -> broker up, coords on the bag
//  2. CreateMQTTReceiveDevice    -> enabled device subscribes to its downlink topic
//  3. PublishDownlink            -> a retained message is injected on that topic
//  4. AssertDownlinkLogged       -> a down/downlink frame appears in the logs
//     (compensation) delete the device, close the broker
//
// PublishDownlink uses a retained message, so it is delivered whether it lands
// before or after the device's subscription becomes active — no settle needed.
func Items() []saga.Item {
	return []saga.Item{
		targetSteps.StartMQTTBroker(),
		deviceSteps.CreateMQTTReceiveDevice(),
		targetSteps.PublishDownlink(),
		logAsserts.AssertDownlinkLogged(),
	}
}

// Run executes the journey against the live simulator.
func Run(t *testing.T) {
	t.Helper()
	if err := utils.SetupE2EEnvironment(); err != nil {
		t.Skipf("simulator not ready: %v", err)
	}
	runID := random.NewRunID()
	clients := saga.NewClientSet(constants.SimURL)
	saga.Run(t, context.Background(), runID, clients, Items()...)
}
