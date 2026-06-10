package console

import (
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	"github.com/gofiber/fiber/v2"

	"simulator/service/src/modules/console/application/ports"
	"simulator/service/src/modules/console/infrastructure/hub"
	"simulator/service/src/modules/console/interfaces/ws"
)

// InitServices registers the console hub and exposes it as the Publisher port the
// engine streams through.
func InitServices() {
	c := container.GetContainer()
	if err := c.Provide(hub.New); err != nil {
		logger.Panic("[MODULE:Console] provide hub: " + err.Error())
	}
	if err := c.Provide(func(h *hub.Hub) ports.Publisher { return h }); err != nil {
		logger.Panic("[MODULE:Console] provide publisher: " + err.Error())
	}
	logger.Info("[MODULE:Console] Services registered")
}

// InitInterfaces mounts the console WebSocket route.
func InitInterfaces() {
	c := container.GetContainer()
	if err := c.Invoke(func(app *fiber.App, h *hub.Hub) {
		ws.Register(app, h)
	}); err != nil {
		logger.Panic("[MODULE:Console] register interfaces: " + err.Error())
	}
	logger.Info("[MODULE:Console] Interfaces registered")
}
