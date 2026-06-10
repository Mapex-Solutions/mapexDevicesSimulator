package engine

import (
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	ports "simulator/service/src/modules/engine/application/ports"
	service "simulator/service/src/modules/engine/application/services"
	dispatch "simulator/service/src/modules/engine/infrastructure/dispatch"
)

// InitServices registers the dispatcher registry and the engine service.
func InitServices() {
	c := container.GetContainer()
	if err := c.Provide(dispatch.NewRegistry); err != nil {
		logger.Panic("[MODULE:Engine] provide registry: " + err.Error())
	}
	if err := c.Provide(service.New); err != nil {
		logger.Panic("[MODULE:Engine] provide service: " + err.Error())
	}
	logger.Info("[MODULE:Engine] Services registered")
}

// InitInterfaces fires the engine OnMount: it reads the devices and starts the
// scheduler. This runs after every module's services are registered, so the
// device, log, and console ports it depends on are resolvable.
func InitInterfaces() {
	c := container.GetContainer()
	if err := c.Invoke(func(e ports.EnginePort) {
		e.OnMount()
	}); err != nil {
		logger.Panic("[MODULE:Engine] mount engine: " + err.Error())
	}
	logger.Info("[MODULE:Engine] Interfaces registered (engine mounted)")
}
