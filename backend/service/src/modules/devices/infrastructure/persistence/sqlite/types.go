package sqlite

import (
	sqliteModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/sqlite/model"

	"simulator/service/src/modules/devices/domain/entities"
)

// adapter implements the DeviceRepository port over the shared sqlite model.
type adapter struct {
	model *sqliteModel.Model[entities.Device]
}
