package gateways

import (
	"context"

	sqliteManager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/sqlite/manager"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	"github.com/gofiber/fiber/v2"

	"simulator/service/src/modules/gateways/application/ports"
	service "simulator/service/src/modules/gateways/application/services"
	"simulator/service/src/modules/gateways/domain/repositories"
	sqliteRepo "simulator/service/src/modules/gateways/infrastructure/persistence/sqlite"
	"simulator/service/src/modules/gateways/interfaces/http/routes"
)

// InitRepositories migrates the gateways table and registers the repository port.
func InitRepositories() {
	c := container.GetContainer()
	if err := c.Invoke(func(mgr *sqliteManager.SQLiteManager) {
		if err := mgr.Migrate(context.Background(), sqliteRepo.DDL); err != nil {
			logger.Panic("[MODULE:Gateways] migrate: " + err.Error())
		}
	}); err != nil {
		logger.Panic("[MODULE:Gateways] resolve manager: " + err.Error())
	}
	if err := c.Provide(func(mgr *sqliteManager.SQLiteManager) repositories.GatewayRepository {
		return sqliteRepo.New(mgr)
	}); err != nil {
		logger.Panic("[MODULE:Gateways] provide repository: " + err.Error())
	}
	logger.Info("[MODULE:Gateways] Repositories registered")
}

// InitServices registers the gateways service.
func InitServices() {
	c := container.GetContainer()
	if err := c.Provide(service.New); err != nil {
		logger.Panic("[MODULE:Gateways] provide service: " + err.Error())
	}
	logger.Info("[MODULE:Gateways] Services registered")
}

// InitInterfaces creates the gateway route group and registers its HTTP routes
// over the resolved service port.
func InitInterfaces() {
	c := container.GetContainer()
	if err := c.Invoke(func(app *fiber.App, service ports.GatewaysServicePort) {
		group := app.Group("/api/gateways")
		routes.RegisterRoutes(group, service)
		logger.Info("[MODULE:Gateways] Routes registered")
	}); err != nil {
		logger.Panic("[MODULE:Gateways] register interfaces: " + err.Error())
	}
}
