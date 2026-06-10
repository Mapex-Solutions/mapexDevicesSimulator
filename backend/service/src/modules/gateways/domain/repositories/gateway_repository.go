package repositories

import (
	"context"
	"errors"

	"simulator/service/src/modules/gateways/domain/entities"
)

// ErrNotFound is returned when a gateway id does not exist, so the application
// layer can map it to a 404 without importing the persistence technology's error.
var ErrNotFound = errors.New("gateway not found")

// GatewayRepository is the persistence port for gateways.
type GatewayRepository interface {
	Create(ctx context.Context, g *entities.Gateway) (*entities.Gateway, error)
	GetByID(ctx context.Context, id string) (*entities.Gateway, error)
	List(ctx context.Context) ([]entities.Gateway, error)
	Update(ctx context.Context, id string, g *entities.Gateway) (*entities.Gateway, error)
	Delete(ctx context.Context, id string) error
}
