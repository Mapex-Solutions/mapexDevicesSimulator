package steps

import (
	"fmt"

	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/common/constants"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/common/utils"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/core/saga"
)

// StartConsoleStream connects to the simulator's realtime console WebSocket and
// starts collecting frames, so later asserts can read what reached the UI live.
// It runs early so it is subscribed before the steps that produce frames.
//
// Writes (bag): BagKeyConsoleStream (*utils.ConsoleStream).
func StartConsoleStream() saga.Step {
	return saga.Step{
		Name: "targets.StartConsoleStream",
		Do: func(c *saga.Context) error {
			stream, err := utils.StartConsoleStream(constants.ConsoleWSURL)
			if err != nil {
				return fmt.Errorf("connect console ws: %w", err)
			}
			c.Set(BagKeyConsoleStream, stream)
			return nil
		},
		Compensate: func(c *saga.Context) error {
			if v, ok := c.Get(BagKeyConsoleStream); ok {
				if stream, ok := v.(*utils.ConsoleStream); ok {
					stream.Close()
				}
			}
			return nil
		},
	}
}
