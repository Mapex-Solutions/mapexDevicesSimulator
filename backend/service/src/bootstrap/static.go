package bootstrap

import (
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	"simulator/service/src/shared/web"
)

// InitStatic serves the SPA: from spa_dir when set (dev/testing), otherwise from
// the binary-embedded build. It must be registered LAST, after the API and WS
// routes, because it is a catch-all. Unmatched /api and /ws paths are skipped so
// they 404 as API calls instead of returning the SPA shell; every other path
// falls back to index.html for the history-mode router.
func InitStatic(app *fiber.App) {
	root := spaFS()
	if root == nil {
		logger.Warn("[INFRA:HTTP] no SPA assets embedded; serving API only")
		return
	}

	app.Use("/", filesystem.New(filesystem.Config{
		Next:         skipAPIRoutes,
		Root:         http.FS(root),
		Index:        "index.html",
		NotFoundFile: "index.html",
	}))
	logger.Info("[INFRA:HTTP] SPA static serving enabled")
}

// skipAPIRoutes keeps the SPA catch-all off the API and WS namespaces.
func skipAPIRoutes(c *fiber.Ctx) bool {
	p := c.Path()
	return strings.HasPrefix(p, "/api") || strings.HasPrefix(p, "/ws")
}

// spaFS returns the filesystem holding the SPA: the spa_dir override when it
// exists, otherwise the embedded build (nil when only the placeholder is present).
func spaFS() fs.FS {
	if dir, _ := config.GetStringValue("spa_dir"); dir != "" {
		if _, err := os.Stat(dir); err == nil {
			logger.Info("[INFRA:HTTP] serving SPA from spa_dir=" + dir)
			return os.DirFS(dir)
		}
		logger.Warn("[INFRA:HTTP] spa_dir not found, falling back to embedded: " + dir)
	}
	return web.SPA()
}
