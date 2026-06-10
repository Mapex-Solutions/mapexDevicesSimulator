package bootstrap

import (
	"context"
	"os"
	"path/filepath"

	sqliteManager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/sqlite/manager"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	"go.uber.org/dig"

	"simulator/service/src/shared/persistence"
)

// InitSQLite opens the on-disk database, runs the base schema migration, and
// provides the manager to the container for the modules' repositories.
func InitSQLite(c *dig.Container) {
	path, _ := config.GetStringValue("sqlite_path")
	if err := ensureParentDir(path); err != nil {
		logger.Panic("[INFRA:SQLITE] prepare data dir: " + err.Error())
	}

	mgr, err := sqliteManager.New(sqliteManager.Config{Path: path, ForeignKeys: true})
	if err != nil {
		logger.Panic("[INFRA:SQLITE] open: " + err.Error())
	}
	if err := mgr.Migrate(context.Background(), persistence.Schema...); err != nil {
		logger.Panic("[INFRA:SQLITE] migrate: " + err.Error())
	}
	if err := c.Provide(func() *sqliteManager.SQLiteManager { return mgr }); err != nil {
		logger.Panic("[INFRA:SQLITE] provide: " + err.Error())
	}
}

// ensureParentDir creates the directory holding the database file when needed.
func ensureParentDir(path string) error {
	dir := filepath.Dir(path)
	if dir == "" || dir == "." {
		return nil
	}
	return os.MkdirAll(dir, 0o755)
}
