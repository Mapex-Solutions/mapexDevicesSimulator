package routes

import (
	"github.com/gofiber/fiber/v2"

	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"

	"simulator/service/src/modules/logs/application/dtos"
	"simulator/service/src/modules/logs/application/ports"
	"simulator/service/src/modules/logs/interfaces/http/handlers"
)

// RegisterRoutes registers the read-only log HTTP route.
//
// Following Hexagonal Architecture, this function accepts the service port interface
// rather than a concrete service implementation. The query string is validated and
// bound into a LogQuery DTO by the requestValidation middleware before the handler
// runs.
//
// Base path: /api/logs
//
//	GET / - List logs (paginated, filtered)
//
// Parameters:
//   - group: the Fiber router group to register the route on
//   - service: the LogsServicePort implementation
func RegisterRoutes(group fiber.Router, service ports.LogsServicePort) {
	listValidation := validation.NewValidation(nil, &dtos.LogQuery{}, nil)
	group.Get("/",
		validation.ValidationMiddleware(listValidation),
		handlers.ListLogs(service),
	)
}
