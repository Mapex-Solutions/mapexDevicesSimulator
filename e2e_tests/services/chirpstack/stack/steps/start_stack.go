package steps

import (
	"context"
	"fmt"
	"time"

	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/common/constants"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/common/utils"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/core/saga"
	csclient "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/chirpstack/client"
)

// stackUpTimeout bounds how long the compose stack may take to come up; the
// chirpstack login poll inside Connect provides the real readiness gate.
const (
	stackUpTimeout = 3 * time.Minute
	loginTimeout   = 90 * time.Second
)

// StartStack brings the pinned ChirpStack docker stack up and connects the gRPC
// client, publishing both to the bag. Its Compensate, running last in the saga
// rollback, closes the client and takes the stack down with its volumes — a
// fresh LNS every run.
//
// Writes (bag):
//   - BagKeyStack   *utils.ChirpStackStack
//   - BagKeyClient  *client.Client (logged in)
func StartStack() saga.Step {
	return saga.Step{
		Name: "chirpstack.StartStack",
		Do: func(c *saga.Context) error {
			stack := utils.NewChirpStackStack()
			c.Set(BagKeyStack, stack)

			upCtx, cancel := context.WithTimeout(c.Stdctx, stackUpTimeout)
			defer cancel()
			if err := stack.Up(upCtx); err != nil {
				return fmt.Errorf("chirpstack stack up: %w", err)
			}

			loginCtx, cancelLogin := context.WithTimeout(c.Stdctx, loginTimeout)
			defer cancelLogin()
			cli, err := csclient.Connect(loginCtx, constants.ChirpStackGRPCAddr, constants.ChirpStackAdminUser, constants.ChirpStackAdminPass)
			if err != nil {
				return fmt.Errorf("chirpstack connect: %w", err)
			}
			c.Set(BagKeyClient, cli)
			return nil
		},
		Compensate: func(c *saga.Context) error {
			if v, ok := c.Get(BagKeyClient); ok {
				if cli, ok := v.(*csclient.Client); ok {
					_ = cli.Close()
				}
			}
			v, ok := c.Get(BagKeyStack)
			if !ok {
				return nil
			}
			stack, ok := v.(*utils.ChirpStackStack)
			if !ok {
				return nil
			}
			downCtx, cancel := context.WithTimeout(context.Background(), stackUpTimeout)
			defer cancel()
			return stack.Down(downCtx)
		},
	}
}
