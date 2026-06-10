package bootstrap

import (
	"context"

	sqliteManager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/sqlite/manager"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/shutdown"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/dig"

	enginePorts "simulator/service/src/modules/engine/application/ports"
)

// InitShutdown registers graceful shutdown hooks in order: stop accepting HTTP
// (priority 0), stop the simulation engine (priority 1, before the DB so its last
// log writes land), then close the database (priority 5).
func InitShutdown(c *dig.Container, sm *shutdown.ShutdownManager, app *fiber.App) {
	sm.RegisterFunc("http", 0, func(ctx context.Context) error {
		return app.ShutdownWithContext(ctx)
	})

	if err := c.Invoke(func(e enginePorts.EnginePort) {
		sm.RegisterFunc("engine", 1, e.OnShutdown)
	}); err != nil {
		logger.Error(err, "[INFRA:Shutdown] resolve engine for shutdown")
	}

	if err := c.Invoke(func(db *sqliteManager.SQLiteManager) {
		sm.RegisterFunc("sqlite", 5, func(_ context.Context) error {
			return db.Close()
		})
	}); err != nil {
		logger.Error(err, "[INFRA:Shutdown] resolve sqlite for shutdown")
	}
}
