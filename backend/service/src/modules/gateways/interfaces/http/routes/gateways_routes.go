package routes

import (
	"github.com/gofiber/fiber/v2"

	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"

	"simulator/service/src/modules/gateways/application/dtos"
	"simulator/service/src/modules/gateways/application/ports"
	"simulator/service/src/modules/gateways/interfaces/http/handlers"
)

// RegisterRoutes registers the gateway HTTP routes.
//
// Following Hexagonal Architecture, this function accepts the service port interface
// rather than a concrete service implementation. Each mutating route validates and
// binds its body / path params through the requestValidation middleware before the
// handler runs, so handlers read ready-made DTOs and never touch the raw request.
//
// Base path: /api/gateways
//
//	GET    /        - List gateways
//	POST   /        - Create gateway
//	PUT    /:id     - Replace gateway by id
//	DELETE /:id     - Delete gateway by id
//
// Parameters:
//   - group: the Fiber router group to register the routes on
//   - service: the GatewaysServicePort implementation
func RegisterRoutes(group fiber.Router, service ports.GatewaysServicePort) {
	// List gateways
	group.Get("/", handlers.ListGateways(service))

	// Create gateway (validate body)
	createValidation := validation.NewValidation(&dtos.GatewayInput{}, nil, nil)
	group.Post("/",
		validation.ValidationMiddleware(createValidation),
		handlers.CreateGateway(service),
	)

	// Replace gateway by id (validate body + path param)
	updateValidation := validation.NewValidation(&dtos.GatewayInput{}, nil, &dtos.GatewayIDParam{})
	group.Put("/:id",
		validation.ValidationMiddleware(updateValidation),
		handlers.UpdateGateway(service),
	)

	// Delete gateway by id (validate path param)
	deleteValidation := validation.NewValidation(nil, nil, &dtos.GatewayIDParam{})
	group.Delete("/:id",
		validation.ValidationMiddleware(deleteValidation),
		handlers.DeleteGateway(service),
	)
}
