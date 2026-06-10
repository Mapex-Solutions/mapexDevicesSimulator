package devices

import (
	"context"

	sqliteManager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/sqlite/manager"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	"github.com/gofiber/fiber/v2"

	"simulator/service/src/modules/devices/application/ports"
	service "simulator/service/src/modules/devices/application/services"
	"simulator/service/src/modules/devices/domain/repositories"
	sqliteRepo "simulator/service/src/modules/devices/infrastructure/persistence/sqlite"
	"simulator/service/src/modules/devices/interfaces/http/routes"
)

// InitRepositories migrates the devices table and registers the repository port.
func InitRepositories() {
	c := container.GetContainer()
	if err := c.Invoke(func(mgr *sqliteManager.SQLiteManager) {
		if err := mgr.Migrate(context.Background(), sqliteRepo.DDL); err != nil {
			logger.Panic("[MODULE:Devices] migrate: " + err.Error())
		}
	}); err != nil {
		logger.Panic("[MODULE:Devices] resolve manager: " + err.Error())
	}
	if err := c.Provide(func(mgr *sqliteManager.SQLiteManager) repositories.DeviceRepository {
		return sqliteRepo.New(mgr)
	}); err != nil {
		logger.Panic("[MODULE:Devices] provide repository: " + err.Error())
	}
	logger.Info("[MODULE:Devices] Repositories registered")
}

// InitServices registers the devices service.
func InitServices() {
	c := container.GetContainer()
	if err := c.Provide(service.New); err != nil {
		logger.Panic("[MODULE:Devices] provide service: " + err.Error())
	}
	logger.Info("[MODULE:Devices] Services registered")
}

// InitInterfaces creates the device route group and registers its HTTP routes
// over the resolved service port.
func InitInterfaces() {
	c := container.GetContainer()
	if err := c.Invoke(func(app *fiber.App, service ports.DevicesServicePort) {
		group := app.Group("/api/devices")
		routes.RegisterRoutes(group, service)
		logger.Info("[MODULE:Devices] Routes registered")
	}); err != nil {
		logger.Panic("[MODULE:Devices] register interfaces: " + err.Error())
	}
}
