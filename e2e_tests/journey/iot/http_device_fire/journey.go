// Package http_device_fire exercises an HTTP device end to end against the live
// simulator: start an echo target, create the device, fire its event, and
// confirm the simulator logged a 200 uplink with the captured response.
//
// Outcome on PASS:
//   - An HTTP device is created targeting an in-process echo and fired once.
//   - GET /api/logs shows the up/data frame with status 200 and a response.
//   - The device and echo are torn down (compensation), leaving the simulator
//     clean.
//
// Outcome on FAIL:
//   - The failing step / assert name surfaces in the saga log
//     (devices.CreateHTTPDevice, logs.AssertHTTPUplinkLogged, ...).
package http_device_fire

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
	targetSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/targets/steps"
)

// Items is the ordered saga the journey runs.
//
//  1. StartEcho                -> in-process echo target, URL on the bag
//  2. CreateHTTPDevice         -> device persisted targeting the echo
//  3. FireTelemetry            -> one uplink dispatched
//  4. AssertHTTPUplinkLogged   -> a 200 uplink with a response is in the logs
//     (compensation) delete the device, close the echo
func Items() []saga.Item {
	return []saga.Item{
		targetSteps.StartEcho(),
		deviceSteps.CreateHTTPDevice(),
		engineSteps.FireTelemetry(),
		logAsserts.AssertHTTPUplinkLogged(),
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
