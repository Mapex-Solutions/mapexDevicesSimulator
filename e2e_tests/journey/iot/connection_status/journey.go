// Package connection_status exercises the engine's connection-status lifecycle
// over the realtime console WebSocket: an enabled device pointed at an
// unreachable broker must surface connecting / reconnecting frames live. These
// frames exist only on /ws — they are never written to the logs — so the console
// stream is the only way to observe them.
//
// Outcome on PASS:
//   - The console WebSocket is connected first.
//   - An MQTT device pointing at an unreachable broker is created and enabled.
//   - A system/status connecting or reconnecting frame for the device arrives
//     live on the stream.
//   - Device and stream are torn down (compensation).
//
// Outcome on FAIL:
//   - The failing step / assert name surfaces in the saga log
//     (console.AssertConsoleReconnecting, ...).
package connection_status

import (
	"context"
	"testing"

	"github.com/Mapex-Solutions/mapexGoKit/utils/random"

	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/common/constants"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/common/utils"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/core/saga"
	consoleAsserts "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/console/asserts"
	deviceSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/devices/steps"
	targetSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/targets/steps"
)

// Items is the ordered saga the journey runs.
//
//  1. StartConsoleStream           -> subscribe to /ws first
//  2. CreateUnreachableMQTTDevice  -> enabled device whose broker is unreachable
//  3. AssertConsoleReconnecting    -> a connecting/reconnecting frame arrives live
//     (compensation) delete the device, close the stream
func Items() []saga.Item {
	return []saga.Item{
		targetSteps.StartConsoleStream(),
		deviceSteps.CreateUnreachableMQTTDevice(),
		consoleAsserts.AssertConsoleReconnecting(),
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
