package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"

	"simulator/service/src/modules/gateways/application/dtos"
	"simulator/service/src/modules/gateways/application/ports"
)

// ListGateways returns a Fiber handler that lists every simulated gateway.
//
// Following Hexagonal Architecture, the handler accepts the service port interface
// and delegates the work to the service layer. Any error is returned to the global
// error handler, which renders it as the standard envelope.
//
// Parameters:
//   - service: the GatewaysServicePort for gateway business operations
//
// Returns:
//   - 200 OK with the array of gateway DTOs
func ListGateways(service ports.GatewaysServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		retData, err := service.List(c.UserContext())
		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// CreateGateway returns a Fiber handler that creates a gateway.
//
// It expects a validated DTO of type dtos.GatewayInput under the key "bodyDTO"
// (populated by the requestValidation middleware).
//
// Parameters:
//   - service: the GatewaysServicePort for gateway business operations
//
// Returns:
//   - 201 Created with the stored gateway DTO
func CreateGateway(service ports.GatewaysServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		bodyData, _ := requestValidation.GetDTO[*dtos.GatewayInput](c, "bodyDTO")
		retData, err := service.Create(c.UserContext(), bodyData)
		if err != nil {
			return err
		}
		return response.Created(c, retData)
	}
}

// UpdateGateway returns a Fiber handler that replaces a gateway by id.
//
// It expects two validated DTOs:
//   - dtos.GatewayIDParam under the key "paramsDTO"
//   - dtos.GatewayInput under the key "bodyDTO"
//
// Parameters:
//   - service: the GatewaysServicePort for gateway business operations
//
// Returns:
//   - 200 OK with the updated gateway DTO
//   - 404 Not Found if the gateway does not exist (raised by the service)
func UpdateGateway(service ports.GatewaysServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		params, _ := requestValidation.GetDTO[*dtos.GatewayIDParam](c, "paramsDTO")
		bodyData, _ := requestValidation.GetDTO[*dtos.GatewayInput](c, "bodyDTO")
		retData, err := service.Update(c.UserContext(), params.ID, bodyData)
		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// DeleteGateway returns a Fiber handler that deletes a gateway by id.
//
// It expects a validated DTO of type dtos.GatewayIDParam under the key "paramsDTO".
//
// Parameters:
//   - service: the GatewaysServicePort for gateway business operations
//
// Returns:
//   - 200 OK with {"success": true}
//   - 404 Not Found if the gateway does not exist (raised by the service)
func DeleteGateway(service ports.GatewaysServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		params, _ := requestValidation.GetDTO[*dtos.GatewayIDParam](c, "paramsDTO")
		retData, err := service.Delete(c.UserContext(), params.ID)
		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}
