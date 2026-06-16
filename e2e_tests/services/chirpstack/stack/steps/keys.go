// Package steps owns the ChirpStack stack lifecycle as saga steps: bringing the
// pinned docker stack up and connecting the API client at the start of a
// journey, and tearing both down during rollback. Keeping the lifecycle inside
// the saga chain means "up on start, down on finish" falls out of the runner's
// reverse-order compensation for free.
package steps

import (
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/core/saga"
	csclient "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/chirpstack/client"
)

const (
	// BagKeyClient holds the connected *client.Client other ChirpStack steps and
	// asserts drive the API with.
	BagKeyClient = "chirpstack.client"
	// BagKeyStack holds the *utils.ChirpStackStack handle so Compensate can take
	// the stack down.
	BagKeyStack = "chirpstack.stack"
)

// ClientFromBag fetches the connected ChirpStack client published by StartStack,
// failing the test fast when it is missing (a step ran out of order).
func ClientFromBag(c *saga.Context) *csclient.Client {
	v, ok := c.Get(BagKeyClient)
	if !ok {
		c.T.Fatalf("[SAGA] missing ChirpStack client in bag (StartStack did not run?)")
	}
	cl, ok := v.(*csclient.Client)
	if !ok {
		c.T.Fatalf("[SAGA] bag key %q is not a *client.Client (got %T)", BagKeyClient, v)
	}
	return cl
}
