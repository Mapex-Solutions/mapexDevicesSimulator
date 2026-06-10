package bootstrap

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	customErrors "github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	"go.uber.org/dig"
)

// InitFiber builds the HTTP app and provides it to the container so modules can
// register their routes. The shared gokit error handler renders every error a
// handler returns as the standard {status, errors, data} envelope, honoring the
// HTTP code carried by a ServerCustomError. CORS is enabled for the dev SPA.
func InitFiber(c *dig.Container) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "simulatord",
		ErrorHandler: customErrors.FiberErrorHandler,
	})

	if origins, _ := config.GetStringValue("cors_origins"); origins != "" {
		app.Use(cors.New(cors.Config{
			AllowOrigins: origins,
			AllowMethods: strings.Join([]string{
				fiber.MethodGet, fiber.MethodPost, fiber.MethodPut,
				fiber.MethodDelete, fiber.MethodOptions,
			}, ","),
		}))
	}

	if err := c.Provide(func() *fiber.App { return app }); err != nil {
		logger.Panic("[INFRA:HTTP] provide fiber: " + err.Error())
	}
	return app
}

// ListenAddress resolves the bind address. A positive flagPort (passed by the
// Electron sidecar launcher) wins over the configured http_port.
func ListenAddress(flagPort int) string {
	host, _ := config.GetStringValue("http_address")
	port := flagPort
	if port <= 0 {
		port, _ = config.GetIntValue("http_port")
	}
	return fmt.Sprintf("%s:%d", host, port)
}
