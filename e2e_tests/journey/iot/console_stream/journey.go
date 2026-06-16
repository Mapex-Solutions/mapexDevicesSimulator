// Package console_stream exercises the realtime console WebSocket end to end: a
// fired uplink must be broadcast live on /ws, not only persisted to the logs.
//
// Outcome on PASS:
//   - The console WebSocket is connected before any frame is produced.
//   - An HTTP device is created against an in-test echo and fired.
//   - An up/data frame for the device arrives on the live stream.
//   - Device, echo, and stream are torn down (compensation).
//
// Outcome on FAIL:
//   - The failing step / assert name surfaces in the saga log
//     (console.AssertConsoleUpFrame, ...).
package console_stream

import (
	"context"
	"testing"

	"github.com/Mapex-Solutions/mapexGoKit/utils/random"

	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/common/constants"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/common/utils"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/core/saga"
	consoleAsserts "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/console/asserts"
	deviceSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/devices/steps"
	engineSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/engine/steps"
	targetSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/targets/steps"
)

// Items is the ordered saga the journey runs.
//
//  1. StartConsoleStream    -> subscribe to /ws before any frame is produced
//  2. StartEcho             -> in-test echo target
//  3. CreateHTTPDevice      -> device targeting the echo
//  4. FireTelemetry         -> one uplink dispatched
//  5. AssertConsoleUpFrame  -> the up/data frame arrives live on /ws
//     (compensation) delete the device, close the echo and the stream
func Items() []saga.Item {
	return []saga.Item{
		targetSteps.StartConsoleStream(),
		targetSteps.StartEcho(),
		deviceSteps.CreateHTTPDevice(),
		engineSteps.FireTelemetry(),
		consoleAsserts.AssertConsoleUpFrame(),
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
