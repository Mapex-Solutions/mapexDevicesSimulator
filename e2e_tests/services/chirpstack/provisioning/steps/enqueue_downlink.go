package steps

import (
	"encoding/hex"
	"fmt"

	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/core/saga"
	stackSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/chirpstack/stack/steps"
)

// EnqueueDownlink queues a downlink on ChirpStack for the active device. For a
// Class A device the LNS sends it in the RX window after the next uplink, so the
// journey fires an uplink afterwards to trigger delivery.
//
// Reads (bag): stackSteps.BagKeyClient, BagKeyActiveDevEUI.
// Writes (bag): BagKeyDownlinkHex (the payload's hex, for the inbound assert).
func EnqueueDownlink() saga.Step {
	return saga.Step{
		Name: "chirpstack.EnqueueDownlink",
		Do: func(c *saga.Context) error {
			cli := stackSteps.ClientFromBag(c)
			devEUI := c.MustGetString(BagKeyActiveDevEUI)
			if err := cli.EnqueueDownlink(c.Stdctx, devEUI, downlinkFPort, downlinkBytes); err != nil {
				return fmt.Errorf("enqueue downlink: %w", err)
			}
			c.Set(BagKeyDownlinkHex, hex.EncodeToString(downlinkBytes))
			return nil
		},
	}
}
