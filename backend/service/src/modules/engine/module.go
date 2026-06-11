package engine

import (
	"github.com/gofiber/fiber/v2"

	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	"go.uber.org/dig"

	ports "simulator/service/src/modules/engine/application/ports"
	service "simulator/service/src/modules/engine/application/services"
	dispatch "simulator/service/src/modules/engine/infrastructure/dispatch"
	session "simulator/service/src/modules/engine/infrastructure/session"
	routes "simulator/service/src/modules/engine/interfaces/http/routes"
	"simulator/service/src/shared/reconcile"
)

// InitServices registers the dispatcher registry (one-shot sends), the connector
// registry (persistent sessions), and the engine service.
func InitServices() {
	c := container.GetContainer()
	provideReconcile(c)
	if err := c.Provide(dispatch.NewRegistry); err != nil {
		logger.Panic("[MODULE:Engine] provide registry: " + err.Error())
	}
	if err := c.Provide(session.NewConnectorRegistry); err != nil {
		logger.Panic("[MODULE:Engine] provide connector registry: " + err.Error())
	}
	if err := c.Provide(service.New); err != nil {
		logger.Panic("[MODULE:Engine] provide service: " + err.Error())
	}
	logger.Info("[MODULE:Engine] Services registered")
}

// provideReconcile registers the shared CRUD-change notifier and exposes its two
// sides: the writer (Signal) the devices and gateways services raise, and the
// subscriber (Listener) the engine binds its Reconcile to. It lives here because
// the engine is the listener of record; the type itself is neutral so neither
// the devices nor the gateways module imports the engine.
func provideReconcile(c *dig.Container) {
	if err := c.Provide(func() *reconcile.Notifier { return reconcile.New() }); err != nil {
		logger.Panic("[MODULE:Engine] provide notifier: " + err.Error())
	}
	if err := c.Provide(func(n *reconcile.Notifier) reconcile.Signal { return n }); err != nil {
		logger.Panic("[MODULE:Engine] provide reconcile signal: " + err.Error())
	}
	if err := c.Provide(func(n *reconcile.Notifier) reconcile.Listener { return n }); err != nil {
		logger.Panic("[MODULE:Engine] provide reconcile listener: " + err.Error())
	}
}

// InitInterfaces fires the engine OnMount (reads devices, starts the scheduler and
// the live sessions) and registers the on-demand fire route. This runs after every
// module's services are registered, so the device, log, and console ports it
// depends on are resolvable.
func InitInterfaces() {
	c := container.GetContainer()
	if err := c.Invoke(func(app *fiber.App, e ports.EnginePort) {
		e.OnMount()
		group := app.Group("/api/devices")
		routes.RegisterRoutes(group, e)
	}); err != nil {
		logger.Panic("[MODULE:Engine] mount engine: " + err.Error())
	}
	logger.Info("[MODULE:Engine] Interfaces registered (engine mounted, fire route)")
}
