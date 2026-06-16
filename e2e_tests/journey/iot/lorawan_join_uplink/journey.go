// Package lorawan_join_uplink exercises LoRaWAN OTAA devices end to end against
// the live simulator and a pinned ChirpStack stack, over BOTH radio transports:
// Semtech UDP and Basics Station. One stack is brought up for the whole journey
// and torn down with its volumes at the end.
//
// Outcome on PASS:
//   - ChirpStack is provisioned (tenant, application, EU868/1.0.3 OTAA profile).
//   - UDP transport: a gateway + device are registered on the LNS and mirrored
//     in the simulator; enabling the device joins (ChirpStack assigns a DevAddr)
//     and a fired uplink is recorded by ChirpStack (last-seen advances).
//   - Basics Station transport: the same, with the device carrying its own
//     WebSocket link to the bridge instead of a separate gateway.
//   - Every resource is torn down in reverse, and the stack is removed with its
//     volumes (a fresh LNS each run, so OTAA DevNonces never collide).
//
// Outcome on FAIL:
//   - The failing step / assert name surfaces in the saga log
//     (chirpstack.ProvisionUDPDevice, chirpstack.AssertJoinAccepted,
//     devices.CreateBasicStationDevice, ...), pointing at the broken integration.
package lorawan_join_uplink

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
)

// Items is the ordered saga the journey runs. The stack and LNS context come up
// once; each transport then provisions, joins, fires, and verifies in turn.
//
//  1. StartStack                       -> ChirpStack up + client connected
//  2. EnsureApplicationContext         -> tenant + application + OTAA profile
//     UDP transport:
//  3. ProvisionUDPDevice               -> LNS gateway + device + keys
//  4. CreateUDPGateway                 -> simulator UDP gateway
//  5. CreateLoRaWANDevice              -> simulator device (joins on enable)
//  6. AssertJoinAccepted               -> ChirpStack assigned a DevAddr
//  7. FireTelemetry                    -> one uplink dispatched
//  8. AssertUplinkReceived             -> ChirpStack recorded the uplink
//     Basics Station transport:
//  9. ProvisionBasicStationDevice      -> LNS gateway + device + keys
// 10. CreateBasicStationDevice         -> simulator device with its own WS link
// 11. AssertJoinAccepted               -> joined
// 12. FireTelemetry                    -> uplink dispatched
// 13. AssertUplinkReceived             -> recorded
//     (compensation) delete devices/gateways, then stack down -v
func Items() []saga.Item {
	return []saga.Item{
		stackSteps.StartStack(),
		provSteps.EnsureApplicationContext(),

		provSteps.ProvisionUDPDevice(),
		gatewaySteps.CreateUDPGateway(),
		deviceSteps.CreateLoRaWANDevice(),
		provAsserts.AssertJoinAccepted(),
		engineSteps.FireTelemetry(),
		provAsserts.AssertUplinkReceived(),

		provSteps.ProvisionBasicStationDevice(),
		deviceSteps.CreateBasicStationDevice(),
		provAsserts.AssertJoinAccepted(),
		engineSteps.FireTelemetry(),
		provAsserts.AssertUplinkReceived(),
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
