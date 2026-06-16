package asserts

import (
	"fmt"
	"time"

	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/core/saga"
	provSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/chirpstack/provisioning/steps"
	stackSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/chirpstack/stack/steps"
)

// uplinkTimeout bounds how long ChirpStack may take to register the fired uplink.
const uplinkTimeout = 30 * time.Second

// AssertUplinkReceived polls ChirpStack until the active device's last-seen
// timestamp is set, proving the data uplink fired through the simulator reached
// the LNS over the radio path.
//
// Reads (bag): stackSteps.BagKeyClient, provSteps.BagKeyActiveDevEUI.
func AssertUplinkReceived() saga.Assert {
	return saga.Assert{
		Name: "chirpstack.AssertUplinkReceived",
		Check: func(c *saga.Context) error {
			cli := stackSteps.ClientFromBag(c)
			devEUI := c.MustGetString(provSteps.BagKeyActiveDevEUI)
			deadline := time.Now().Add(uplinkTimeout)
			for {
				seen, err := cli.LastSeen(c.Stdctx, devEUI)
				if err == nil && seen {
					return nil
				}
				if time.Now().After(deadline) {
					return fmt.Errorf("device %s uplink not seen by ChirpStack within %s (last err %v)", devEUI, uplinkTimeout, err)
				}
				time.Sleep(time.Second)
			}
		},
	}
}
