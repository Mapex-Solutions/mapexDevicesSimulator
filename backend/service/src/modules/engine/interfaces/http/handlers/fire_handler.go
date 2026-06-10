package handlers

import (
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v2"

	customErrors "github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"

	ports "simulator/service/src/modules/engine/application/ports"
)

// fireRequest is the POST /api/devices/:id/fire body: a pre-registered event id,
// or an inline ad-hoc event (the "Generic" path).
type fireRequest struct {
	EventID string          `json:"eventId"`
	Event   json.RawMessage `json:"event"`
}

// FireDevice returns a Fiber handler that fires one event for a device on demand.
// The render + send happen in the engine; the result (uplink echo, and any async
// downlink) streams to the console over the WebSocket.
//
// Returns:
//   - 200 OK with {"fired": true}
//   - 404 Not Found when the device or event does not exist
//   - 400 Bad Request when the device/event cannot be fired
func FireDevice(engine ports.EnginePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var req fireRequest
		if err := c.BodyParser(&req); err != nil && len(c.Body()) > 0 {
			return &customErrors.ServerCustomError{
				Code:   fiber.StatusBadRequest,
				Errors: []string{"invalid fire body: " + err.Error()},
			}
		}

		err := engine.Fire(c.UserContext(), id, ports.FireInput{EventID: req.EventID, Event: req.Event})
		switch {
		case err == nil:
			return response.Success(c, map[string]bool{"fired": true})
		case errors.Is(err, ports.ErrDeviceNotFound), errors.Is(err, ports.ErrEventNotFound):
			return &customErrors.ServerCustomError{Code: fiber.StatusNotFound, Errors: []string{err.Error()}}
		case errors.Is(err, ports.ErrFireUnsupported):
			return &customErrors.ServerCustomError{Code: fiber.StatusBadRequest, Errors: []string{err.Error()}}
		default:
			return err
		}
	}
}
