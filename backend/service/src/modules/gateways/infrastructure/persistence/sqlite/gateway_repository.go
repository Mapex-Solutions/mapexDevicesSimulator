package sqlite

import (
	"context"
	"errors"

	sqliteManager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/sqlite/manager"
	sqliteModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/sqlite/model"

	"simulator/service/src/modules/gateways/domain/entities"
	"simulator/service/src/modules/gateways/domain/repositories"
)

// Compile-time proof the adapter satisfies the repository port.
var _ repositories.GatewayRepository = (*adapter)(nil)

// New builds the gateway repository over the sqlite model bound to the gateways
// table.
func New(mgr *sqliteManager.SQLiteManager) repositories.GatewayRepository {
	return &adapter{model: sqliteModel.New[entities.Gateway](mgr.DB(), tableGateways, sqliteModel.Config{})}
}

// Create inserts a gateway; the model fills id and created.
func (a *adapter) Create(ctx context.Context, g *entities.Gateway) (*entities.Gateway, error) {
	return a.model.CreateOne(ctx, g)
}

// GetByID returns a gateway or repositories.ErrNotFound.
func (a *adapter) GetByID(ctx context.Context, id string) (*entities.Gateway, error) {
	return notFound(a.model.FindByID(ctx, id))
}

// List returns every gateway.
func (a *adapter) List(ctx context.Context) ([]entities.Gateway, error) {
	return a.model.FindAll(ctx, nil, nil)
}

// Update replaces the mutable fields of a gateway and returns the updated row, or
// repositories.ErrNotFound.
func (a *adapter) Update(ctx context.Context, id string, g *entities.Gateway) (*entities.Gateway, error) {
	patch := sqliteModel.Map{
		"name":        g.Name,
		"eui":         g.EUI,
		"enabled":     g.Enabled,
		"region":      g.Region,
		"description": g.Description,
		"link":        g.Link,
	}
	return notFound(a.model.UpdateByID(ctx, id, patch))
}

// Delete removes a gateway or returns repositories.ErrNotFound.
func (a *adapter) Delete(ctx context.Context, id string) error {
	if err := a.model.DeleteByID(ctx, id); err != nil {
		if errors.Is(err, sqliteModel.ErrNotFound) {
			return repositories.ErrNotFound
		}
		return err
	}
	return nil
}

// notFound maps the model's not-found to the repository port's sentinel.
func notFound(g *entities.Gateway, err error) (*entities.Gateway, error) {
	if errors.Is(err, sqliteModel.ErrNotFound) {
		return nil, repositories.ErrNotFound
	}
	return g, err
}
