package logs

import (
	"context"

	sqliteManager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/sqlite/manager"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	"github.com/gofiber/fiber/v2"

	"simulator/service/src/modules/logs/application/ports"
	service "simulator/service/src/modules/logs/application/services"
	"simulator/service/src/modules/logs/domain/repositories"
	sqliteRepo "simulator/service/src/modules/logs/infrastructure/persistence/sqlite"
	"simulator/service/src/modules/logs/interfaces/http/routes"
)

// InitRepositories migrates the logs table + index and registers the repository.
func InitRepositories() {
	c := container.GetContainer()
	if err := c.Invoke(func(mgr *sqliteManager.SQLiteManager) {
		if err := mgr.Migrate(context.Background(), sqliteRepo.Migrations...); err != nil {
			logger.Panic("[MODULE:Logs] migrate: " + err.Error())
		}
	}); err != nil {
		logger.Panic("[MODULE:Logs] resolve manager: " + err.Error())
	}
	if err := c.Provide(func(mgr *sqliteManager.SQLiteManager) repositories.LogRepository {
		return sqliteRepo.New(mgr)
	}); err != nil {
		logger.Panic("[MODULE:Logs] provide repository: " + err.Error())
	}
	logger.Info("[MODULE:Logs] Repositories registered")
}

// InitServices registers the logs service.
func InitServices() {
	c := container.GetContainer()
	if err := c.Provide(service.New); err != nil {
		logger.Panic("[MODULE:Logs] provide service: " + err.Error())
	}
	if err := c.Provide(service.NewWriter); err != nil {
		logger.Panic("[MODULE:Logs] provide writer: " + err.Error())
	}
	logger.Info("[MODULE:Logs] Services registered")
}

// InitInterfaces creates the log route group and registers its HTTP route over
// the resolved service port.
func InitInterfaces() {
	c := container.GetContainer()
	if err := c.Invoke(func(app *fiber.App, service ports.LogsServicePort) {
		group := app.Group("/api/logs")
		routes.RegisterRoutes(group, service)
		logger.Info("[MODULE:Logs] Routes registered")
	}); err != nil {
		logger.Panic("[MODULE:Logs] register interfaces: " + err.Error())
	}
}
