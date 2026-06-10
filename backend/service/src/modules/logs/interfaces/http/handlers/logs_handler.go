package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"

	"simulator/service/src/modules/logs/application/dtos"
	"simulator/service/src/modules/logs/application/ports"
)

// ListLogs returns a Fiber handler that serves the paginated, filterable log
// history behind the live console stream.
//
// Following Hexagonal Architecture, the handler accepts the service port interface
// and delegates the work to the service layer. Any error is returned to the global
// error handler, which renders it as the standard envelope.
//
// It expects a validated DTO of type dtos.LogQuery under the key "queryDTO"
// (populated by the requestValidation middleware) carrying the pagination and
// optional filters (protocol, kind, direction, device, q).
//
// Parameters:
//   - service: the LogsServicePort for the read-only log history
//
// Returns:
//   - 200 OK with a LogPage ({items, total})
func ListLogs(service ports.LogsServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		queryData, _ := requestValidation.GetDTO[*dtos.LogQuery](c, "queryDTO")
		retData, err := service.List(c.UserContext(), queryData)
		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}
