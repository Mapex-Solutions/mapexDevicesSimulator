// Package asserts holds the ChirpStack-side oracles for the LoRaWAN journey.
// They read the LNS over its gRPC API — never the simulator — so a passing
// assert proves the radio path reached ChirpStack and ChirpStack accepted it.
package asserts

import (
	"fmt"
	"time"

	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/core/saga"
	provSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/chirpstack/provisioning/steps"
	stackSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/chirpstack/stack/steps"
)

// joinTimeout bounds how long the OTAA join may take after the device is
// enabled in the simulator.
const joinTimeout = 30 * time.Second

// AssertJoinAccepted polls ChirpStack until the active device has an assigned
// DevAddr, proving the OTAA join-request reached the LNS and a join-accept was
// issued — the gateway link plus the keys lined up end to end.
//
// Reads (bag): stackSteps.BagKeyClient, provSteps.BagKeyActiveDevEUI.
func AssertJoinAccepted() saga.Assert {
	return saga.Assert{
		Name: "chirpstack.AssertJoinAccepted",
		Check: func(c *saga.Context) error {
			cli := stackSteps.ClientFromBag(c)
			devEUI := c.MustGetString(provSteps.BagKeyActiveDevEUI)
			deadline := time.Now().Add(joinTimeout)
			for {
				addr, err := cli.Activation(c.Stdctx, devEUI)
				if err == nil && addr != "" {
					return nil
				}
				if time.Now().After(deadline) {
					return fmt.Errorf("device %s did not join (no DevAddr) within %s (last err %v)", devEUI, joinTimeout, err)
				}
				time.Sleep(time.Second)
			}
		},
	}
}
