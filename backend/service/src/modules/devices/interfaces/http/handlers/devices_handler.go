package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"

	"simulator/service/src/modules/devices/application/dtos"
	"simulator/service/src/modules/devices/application/ports"
)

// ListDevices returns a Fiber handler that lists every simulated device.
//
// Following Hexagonal Architecture, the handler accepts the service port interface
// and delegates the work to the service layer. Any error is returned to the global
// error handler, which renders it as the standard envelope.
//
// Parameters:
//   - service: the DevicesServicePort for device business operations
//
// Returns:
//   - 200 OK with the array of device DTOs
func ListDevices(service ports.DevicesServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		retData, err := service.List(c.UserContext())
		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// CreateDevice returns a Fiber handler that creates a device.
//
// It expects a validated DTO of type dtos.DeviceInput under the key "bodyDTO"
// (populated by the requestValidation middleware).
//
// Parameters:
//   - service: the DevicesServicePort for device business operations
//
// Returns:
//   - 201 Created with the stored device DTO
func CreateDevice(service ports.DevicesServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		bodyData, _ := requestValidation.GetDTO[*dtos.DeviceInput](c, "bodyDTO")
		retData, err := service.Create(c.UserContext(), bodyData)
		if err != nil {
			return err
		}
		return response.Created(c, retData)
	}
}

// UpdateDevice returns a Fiber handler that replaces a device by id.
//
// It expects two validated DTOs:
//   - dtos.DeviceIDParam under the key "paramsDTO"
//   - dtos.DeviceInput under the key "bodyDTO"
//
// Parameters:
//   - service: the DevicesServicePort for device business operations
//
// Returns:
//   - 200 OK with the updated device DTO
//   - 404 Not Found if the device does not exist (raised by the service)
func UpdateDevice(service ports.DevicesServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		params, _ := requestValidation.GetDTO[*dtos.DeviceIDParam](c, "paramsDTO")
		bodyData, _ := requestValidation.GetDTO[*dtos.DeviceInput](c, "bodyDTO")
		retData, err := service.Update(c.UserContext(), params.ID, bodyData)
		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// DeleteDevice returns a Fiber handler that deletes a device by id.
//
// It expects a validated DTO of type dtos.DeviceIDParam under the key "paramsDTO".
//
// Parameters:
//   - service: the DevicesServicePort for device business operations
//
// Returns:
//   - 200 OK with {"success": true}
//   - 404 Not Found if the device does not exist (raised by the service)
func DeleteDevice(service ports.DevicesServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		params, _ := requestValidation.GetDTO[*dtos.DeviceIDParam](c, "paramsDTO")
		retData, err := service.Delete(c.UserContext(), params.ID)
		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}
