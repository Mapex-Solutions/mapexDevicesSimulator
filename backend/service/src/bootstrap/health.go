package bootstrap

import (
	"github.com/gofiber/fiber/v2"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	response "github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"
)

// InitHealth registers GET /api/health, the liveness probe the Electron launcher
// polls before opening the window. It uses the standard envelope so the SPA's
// client treats it like every other endpoint; `data` is the HealthResponse.
func InitHealth(app *fiber.App) {
	app.Get("/api/health", func(c *fiber.Ctx) error {
		version, _ := config.GetStringValue("service_version")
		return response.Success(c, fiber.Map{
			"status":  "ok",
			"version": version,
		})
	})
}
