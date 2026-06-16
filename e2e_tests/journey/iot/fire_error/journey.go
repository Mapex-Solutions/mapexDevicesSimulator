// Package fire_error exercises the engine's send-error handling end to end: a
// device fired at an unreachable target must surface the failure as an error
// frame in the logs, not drop it silently.
//
// Outcome on PASS:
//   - An HTTP device targeting an unreachable address is created and fired.
//   - GET /api/logs shows a frame with status "error" and the failure reason in
//     its response.
//   - The device is torn down (compensation).
//
// Outcome on FAIL:
//   - The failing step / assert name surfaces in the saga log
//     (devices.CreateUnreachableHTTPDevice, logs.AssertFireErrorLogged).
package fire_error

import (
	"context"
	"testing"

	"github.com/Mapex-Solutions/mapexGoKit/utils/random"

	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/common/constants"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/common/utils"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/core/saga"
	deviceSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/devices/steps"
	engineSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/engine/steps"
	logAsserts "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/logs/asserts"
)

// Items is the ordered saga the journey runs.
//
//  1. CreateUnreachableHTTPDevice  -> device targeting a dead address
//  2. FireTelemetry                -> one send dispatched, which fails
//  3. AssertFireErrorLogged        -> an error frame appears in the logs
//     (compensation) delete the device
func Items() []saga.Item {
	return []saga.Item{
		deviceSteps.CreateUnreachableHTTPDevice(),
		engineSteps.FireTelemetry(),
		logAsserts.AssertFireErrorLogged(),
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
