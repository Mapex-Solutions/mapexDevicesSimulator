package ports

import (
	"context"

	"simulator/service/src/modules/devices/application/dtos"
)

// DevicesServicePort is the driving port for device CRUD. It speaks DTOs; the
// service maps to and from the domain entity internally.
type DevicesServicePort interface {
	List(ctx context.Context) ([]dtos.Device, error)
	Create(ctx context.Context, in *dtos.DeviceInput) (*dtos.Device, error)
	Update(ctx context.Context, id string, in *dtos.DeviceInput) (*dtos.Device, error)
	Delete(ctx context.Context, id string) (map[string]bool, error)
}
