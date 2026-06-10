package ports

import (
	"context"

	"simulator/service/src/modules/gateways/application/dtos"
)

// GatewaysServicePort is the driving port for gateway CRUD.
type GatewaysServicePort interface {
	List(ctx context.Context) ([]dtos.Gateway, error)
	Create(ctx context.Context, in *dtos.GatewayInput) (*dtos.Gateway, error)
	Update(ctx context.Context, id string, in *dtos.GatewayInput) (*dtos.Gateway, error)
	Delete(ctx context.Context, id string) (map[string]bool, error)
}
