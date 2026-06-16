// Package lorawan_downlink exercises the LoRaWAN INBOUND path end to end against
// the live simulator and a pinned ChirpStack stack: a downlink queued on the LNS
// is delivered to the device in its RX window after an uplink, and the simulator
// surfaces it as a downlink in the logs.
//
// Outcome on PASS:
//   - A UDP OTAA device is provisioned and joins ChirpStack.
//   - A downlink is queued on the LNS; an uplink opens the RX window and the LNS
//     sends it.
//   - GET /api/logs shows a down/downlink frame carrying the queued bytes (hex).
//   - Everything is torn down, stack removed with its volumes.
//
// Outcome on FAIL:
//   - The failing step / assert name surfaces in the saga log
//     (chirpstack.EnqueueDownlink, logs.AssertLoRaWANDownlinkReceived, ...).
package lorawan_downlink

import (
	"context"
	"testing"

	"github.com/Mapex-Solutions/mapexGoKit/utils/random"

	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/common/constants"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/common/utils"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/core/saga"
	provAsserts "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/chirpstack/provisioning/asserts"
	provSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/chirpstack/provisioning/steps"
	stackSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/chirpstack/stack/steps"
	deviceSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/devices/steps"
	engineSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/engine/steps"
	gatewaySteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/gateways/steps"
	logAsserts "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/logs/asserts"
)

// Items is the ordered saga the journey runs.
//
//  1. StartStack                        -> ChirpStack up + client connected
//  2. EnsureApplicationContext          -> tenant + application + OTAA profile
//  3. ProvisionUDPDevice                -> LNS gateway + device + keys
//  4. CreateUDPGateway                  -> simulator UDP gateway
//  5. CreateLoRaWANDevice               -> simulator device (joins on enable)
//  6. AssertJoinAccepted                -> ChirpStack assigned a DevAddr
//  7. EnqueueDownlink                   -> a downlink is queued on the LNS
//  8. FireTelemetry                     -> an uplink opens the RX window
//  9. AssertLoRaWANDownlinkReceived     -> the simulator logged the downlink
//     (compensation) delete device/gateway, then stack down -v
func Items() []saga.Item {
	return []saga.Item{
		stackSteps.StartStack(),
		provSteps.EnsureApplicationContext(),
		provSteps.ProvisionUDPDevice(),
		gatewaySteps.CreateUDPGateway(),
		deviceSteps.CreateLoRaWANDevice(),
		provAsserts.AssertJoinAccepted(),
		provSteps.EnqueueDownlink(),
		engineSteps.FireTelemetry(),
		logAsserts.AssertLoRaWANDownlinkReceived(),
	}
}

// Run executes the journey against the live simulator, bringing the ChirpStack
// stack up and down inside the saga chain.
func Run(t *testing.T) {
	t.Helper()
	if err := utils.SetupE2EEnvironment(); err != nil {
		t.Skipf("simulator not ready: %v", err)
	}
	runID := random.NewRunID()
	clients := saga.NewClientSet(constants.SimURL)
	saga.Run(t, context.Background(), runID, clients, Items()...)
}
