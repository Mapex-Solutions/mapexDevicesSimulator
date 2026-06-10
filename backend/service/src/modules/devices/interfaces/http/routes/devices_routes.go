package routes

import (
	"github.com/gofiber/fiber/v2"

	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"

	"simulator/service/src/modules/devices/application/dtos"
	"simulator/service/src/modules/devices/application/ports"
	"simulator/service/src/modules/devices/interfaces/http/handlers"
)

// RegisterRoutes registers the device HTTP routes.
//
// Following Hexagonal Architecture, this function accepts the service port interface
// rather than a concrete service implementation. Each mutating route validates and
// binds its body / path params through the requestValidation middleware before the
// handler runs, so handlers read ready-made DTOs and never touch the raw request.
//
// Base path: /api/devices
//
//	GET    /        - List devices
//	POST   /        - Create device
//	PUT    /:id     - Replace device by id
//	DELETE /:id     - Delete device by id
//
// Parameters:
//   - group: the Fiber router group to register the routes on
//   - service: the DevicesServicePort implementation
func RegisterRoutes(group fiber.Router, service ports.DevicesServicePort) {
	// List devices
	group.Get("/", handlers.ListDevices(service))

	// Create device (validate body)
	createValidation := validation.NewValidation(&dtos.DeviceInput{}, nil, nil)
	group.Post("/",
		validation.ValidationMiddleware(createValidation),
		handlers.CreateDevice(service),
	)

	// Replace device by id (validate body + path param)
	updateValidation := validation.NewValidation(&dtos.DeviceInput{}, nil, &dtos.DeviceIDParam{})
	group.Put("/:id",
		validation.ValidationMiddleware(updateValidation),
		handlers.UpdateDevice(service),
	)

	// Delete device by id (validate path param)
	deleteValidation := validation.NewValidation(nil, nil, &dtos.DeviceIDParam{})
	group.Delete("/:id",
		validation.ValidationMiddleware(deleteValidation),
		handlers.DeleteDevice(service),
	)
}
