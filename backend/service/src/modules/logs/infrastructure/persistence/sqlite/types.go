package sqlite

import (
	sqliteModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/sqlite/model"

	"simulator/service/src/modules/logs/domain/entities"
)

// adapter implements the LogRepository port. It uses the shared sqlite model for
// inserts and raw SQL (via the model's pool) for the filtered, paginated read,
// because the free-text `q` search needs LIKE/OR beyond the model's equality
// filters.
type adapter struct {
	model *sqliteModel.Model[entities.Log]
}
