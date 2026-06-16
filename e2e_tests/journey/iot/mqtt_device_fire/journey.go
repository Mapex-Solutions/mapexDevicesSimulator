// Package mqtt_device_fire exercises MQTT devices end to end against the live
// simulator and an in-process broker, covering BOTH auth modes the platform
// supports: username/password over plain tcp:// and a client certificate over
// mutual-TLS ssl://.
//
// Outcome on PASS:
//   - The broker accepts an authenticated publish from a username/password
//     device over tcp:// and records it.
//   - The broker accepts an authenticated publish from a client-certificate
//     device over ssl:// (RequireAndVerifyClientCert) and records it.
//   - Both devices and the broker are torn down (compensation).
//
// Outcome on FAIL:
//   - The failing step / assert name surfaces in the saga log. A rejected
//     CONNECT (wrong password, untrusted cert) surfaces as the publish never
//     arriving (engine.AssertMQTTPublished).
package mqtt_device_fire

import (
	"context"
	"testing"

	"github.com/Mapex-Solutions/mapexGoKit/utils/random"

	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/common/constants"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/common/utils"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/core/saga"
	engineAsserts "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/engine/asserts"
	engineSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/engine/steps"
	deviceSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/devices/steps"
	targetSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/targets/steps"
)

// Items is the ordered saga the journey runs. One broker serves both auth modes.
//
//  1. StartMQTTBroker             -> plain + TLS listeners, coords on the bag
//     username/password:
//  2. CreateMQTTUserPassDevice    -> device on tcp:// with user/pass
//  3. FireTelemetry               -> one uplink dispatched
//  4. AssertMQTTPublished         -> broker accepted the authenticated publish
//     certificate:
//  5. CreateMQTTTLSDevice         -> device on ssl:// with a client cert
//  6. FireTelemetry               -> one uplink dispatched
//  7. AssertMQTTPublished         -> broker accepted the authenticated publish
//     (compensation) delete both devices, close the broker
//
// The fire fires the instant the device is enabled — before its persistent
// session has connected — which exercises the engine's one-shot fallback. That
// path uses a distinct client id, so it no longer collides with the connecting
// session and the publish lands reliably.
func Items() []saga.Item {
	return []saga.Item{
		targetSteps.StartMQTTBroker(),

		deviceSteps.CreateMQTTUserPassDevice(),
		engineSteps.FireTelemetry(),
		engineAsserts.AssertMQTTPublished(),

		deviceSteps.CreateMQTTTLSDevice(),
		engineSteps.FireTelemetry(),
		engineAsserts.AssertMQTTPublished(),
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
