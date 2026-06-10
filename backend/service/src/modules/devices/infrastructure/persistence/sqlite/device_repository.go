package sqlite

import (
	"context"
	"errors"

	sqliteManager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/sqlite/manager"
	sqliteModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/sqlite/model"

	"simulator/service/src/modules/devices/domain/entities"
	"simulator/service/src/modules/devices/domain/repositories"
)

// Compile-time proof the adapter satisfies the repository port.
var _ repositories.DeviceRepository = (*adapter)(nil)

// New builds the device repository over the sqlite model bound to the devices
// table.
func New(mgr *sqliteManager.SQLiteManager) repositories.DeviceRepository {
	return &adapter{model: sqliteModel.New[entities.Device](mgr.DB(), tableDevices, sqliteModel.Config{})}
}

// Create inserts a device; the model fills id and created.
func (a *adapter) Create(ctx context.Context, d *entities.Device) (*entities.Device, error) {
	return a.model.CreateOne(ctx, d)
}

// GetByID returns a device or repositories.ErrNotFound.
func (a *adapter) GetByID(ctx context.Context, id string) (*entities.Device, error) {
	return notFound(a.model.FindByID(ctx, id))
}

// List returns every device.
func (a *adapter) List(ctx context.Context) ([]entities.Device, error) {
	return a.model.FindAll(ctx, nil, nil)
}

// Update replaces the mutable fields of a device and returns the updated row, or
// repositories.ErrNotFound.
func (a *adapter) Update(ctx context.Context, id string, d *entities.Device) (*entities.Device, error) {
	patch := sqliteModel.Map{
		"name":        d.Name,
		"device_id":   d.DeviceID,
		"protocol_id": d.ProtocolID,
		"enabled":     d.Enabled,
		"store_logs":  d.StoreLogs,
		"config":      d.Config,
		"attributes":  d.Attributes,
		"events":      d.Events,
	}
	return notFound(a.model.UpdateByID(ctx, id, patch))
}

// Delete removes a device or returns repositories.ErrNotFound.
func (a *adapter) Delete(ctx context.Context, id string) error {
	if err := a.model.DeleteByID(ctx, id); err != nil {
		if errors.Is(err, sqliteModel.ErrNotFound) {
			return repositories.ErrNotFound
		}
		return err
	}
	return nil
}

// notFound maps the model's not-found to the repository port's sentinel so the
// persistence technology does not leak into the application layer.
func notFound(d *entities.Device, err error) (*entities.Device, error) {
	if errors.Is(err, sqliteModel.ErrNotFound) {
		return nil, repositories.ErrNotFound
	}
	return d, err
}
