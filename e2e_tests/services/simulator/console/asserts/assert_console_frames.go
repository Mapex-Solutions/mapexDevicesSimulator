// Package asserts holds oracles that read the simulator's realtime console
// WebSocket stream. They assert on what reached the UI live — including the
// connection-status lifecycle, which exists only on this stream and never in the
// logs.
package asserts

import (
	"fmt"
	"time"

	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/core/saga"
	deviceSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/devices/steps"
	targetSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/targets/steps"
)

// consoleTimeout bounds how long to wait for an expected frame to arrive.
const consoleTimeout = 15 * time.Second

// AssertConsoleUpFrame polls the console stream until an up/data frame for the
// device has arrived — proof the uplink was broadcast live to the UI.
//
// Reads (bag): targetSteps.BagKeyConsoleStream, deviceSteps.BagKeyDeviceDeviceID.
func AssertConsoleUpFrame() saga.Assert {
	return saga.Assert{
		Name: "console.AssertConsoleUpFrame",
		Check: func(c *saga.Context) error {
			stream := targetSteps.ConsoleStreamFromBag(c)
			deviceID := c.MustGetString(deviceSteps.BagKeyDeviceDeviceID)
			deadline := time.Now().Add(consoleTimeout)
			for {
				for _, f := range stream.Frames() {
					if f.DeviceID == deviceID && f.Direction == "up" && f.Kind == "data" {
						return nil
					}
				}
				if time.Now().After(deadline) {
					return fmt.Errorf("no up/data console frame for device %q within %s", deviceID, consoleTimeout)
				}
				time.Sleep(200 * time.Millisecond)
			}
		},
	}
}

// AssertConsoleReconnecting polls the console stream until a system/status frame
// reporting connecting or reconnecting for the device has arrived — proof the
// engine surfaces the connection lifecycle live when a broker is unreachable.
//
// Reads (bag): targetSteps.BagKeyConsoleStream, deviceSteps.BagKeyDeviceDeviceID.
func AssertConsoleReconnecting() saga.Assert {
	return saga.Assert{
		Name: "console.AssertConsoleReconnecting",
		Check: func(c *saga.Context) error {
			stream := targetSteps.ConsoleStreamFromBag(c)
			deviceID := c.MustGetString(deviceSteps.BagKeyDeviceDeviceID)
			deadline := time.Now().Add(consoleTimeout)
			for {
				for _, f := range stream.Frames() {
					if f.DeviceID == deviceID && f.Direction == "system" && f.Kind == "status" &&
						(f.Status == "connecting" || f.Status == "reconnecting") {
						return nil
					}
				}
				if time.Now().After(deadline) {
					return fmt.Errorf("no connecting/reconnecting console frame for device %q within %s", deviceID, consoleTimeout)
				}
				time.Sleep(200 * time.Millisecond)
			}
		},
	}
}
