package sqlite

import (
	sqliteModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/sqlite/model"

	"simulator/service/src/modules/gateways/domain/entities"
)

// adapter implements the GatewayRepository port over the shared sqlite model.
type adapter struct {
	model *sqliteModel.Model[entities.Gateway]
}
