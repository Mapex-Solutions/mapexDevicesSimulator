package routes

import (
	"github.com/gofiber/fiber/v2"

	ports "simulator/service/src/modules/engine/application/ports"
	"simulator/service/src/modules/engine/interfaces/http/handlers"
)

// RegisterRoutes registers the engine's on-demand control routes on the device
// group. The engine owns this route (not the devices module) because firing is the
// engine's job and the engine already depends on devices, keeping the module
// dependency one-directional.
//
// Base path: /api/devices
//
//	POST /:id/fire - fire one event for a device on demand
func RegisterRoutes(group fiber.Router, engine ports.EnginePort) {
	group.Post("/:id/fire", handlers.FireDevice(engine))
}
