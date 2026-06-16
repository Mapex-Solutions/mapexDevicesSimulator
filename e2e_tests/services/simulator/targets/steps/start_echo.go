package steps

import (
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/common/utils"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/core/saga"
)

// StartEcho boots the in-process HTTP echo target an HTTP device fires against
// and publishes its URL. Compensate shuts it down.
//
// Writes (bag): BagKeyEchoServer (*utils.Echo), BagKeyEchoURL.
func StartEcho() saga.Step {
	return saga.Step{
		Name: "targets.StartEcho",
		Do: func(c *saga.Context) error {
			echo := utils.StartEcho()
			c.Set(BagKeyEchoServer, echo)
			c.Set(BagKeyEchoURL, echo.URL())
			return nil
		},
		Compensate: func(c *saga.Context) error {
			if v, ok := c.Get(BagKeyEchoServer); ok {
				if echo, ok := v.(*utils.Echo); ok {
					echo.Close()
				}
			}
			return nil
		},
	}
}
