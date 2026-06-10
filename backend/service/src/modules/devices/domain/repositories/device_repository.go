package repositories

import (
	"context"
	"errors"

	"simulator/service/src/modules/devices/domain/entities"
)

// ErrNotFound is returned by the repository when a device id does not exist, so
// the application layer can map it to a 404 without importing the persistence
// technology's own error.
var ErrNotFound = errors.New("device not found")

// DeviceRepository is the persistence port for devices.
type DeviceRepository interface {
	Create(ctx context.Context, d *entities.Device) (*entities.Device, error)
	GetByID(ctx context.Context, id string) (*entities.Device, error)
	List(ctx context.Context) ([]entities.Device, error)
	Update(ctx context.Context, id string, d *entities.Device) (*entities.Device, error)
	Delete(ctx context.Context, id string) error
}
