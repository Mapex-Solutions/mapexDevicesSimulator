package services

import (
	"context"
	"errors"

	customErrors "github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	status "github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
	"github.com/Mapex-Solutions/mapexGoKit/utils/mapper"

	"simulator/service/src/modules/gateways/application/di"
	"simulator/service/src/modules/gateways/application/dtos"
	"simulator/service/src/modules/gateways/application/ports"
	"simulator/service/src/modules/gateways/domain/entities"
	"simulator/service/src/modules/gateways/domain/repositories"
)

// notFound is the not-found error carrying the 404 code the global error handler
// renders as an envelope. It translates the domain's repository sentinel into the
// HTTP-aware application error.
func notFound() error {
	return &customErrors.ServerCustomError{Code: status.NOT_FOUND, Errors: []string{"gateway not found"}}
}

// Compile-time check that the service satisfies its port.
var _ ports.GatewaysServicePort = (*GatewaysService)(nil)

// New builds the gateways service over its injected repository port.
func New(deps di.GatewaysServiceDI) ports.GatewaysServicePort {
	return &GatewaysService{deps: deps}
}

// List returns every stored gateway mapped to its wire DTO.
func (s *GatewaysService) List(ctx context.Context) ([]dtos.Gateway, error) {
	stored, err := s.deps.Repo.List(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]dtos.Gateway, 0, len(stored))
	for i := range stored {
		dto, err := mapper.EntityToDto[entities.Gateway, dtos.Gateway](&stored[i])
		if err != nil {
			return nil, err
		}
		out = append(out, *dto)
	}
	return out, nil
}

// Create maps the body to an entity, persists it, and returns the stored DTO.
func (s *GatewaysService) Create(ctx context.Context, in *dtos.GatewayInput) (*dtos.Gateway, error) {
	entity, err := mapper.DtoToEntity[dtos.GatewayInput, entities.Gateway](in)
	if err != nil {
		return nil, err
	}
	created, err := s.deps.Repo.Create(ctx, entity)
	if err != nil {
		return nil, err
	}
	return mapper.EntityToDto[entities.Gateway, dtos.Gateway](created)
}

// Update maps the body to an entity, replaces the gateway, and returns the DTO.
func (s *GatewaysService) Update(ctx context.Context, id string, in *dtos.GatewayInput) (*dtos.Gateway, error) {
	entity, err := mapper.DtoToEntity[dtos.GatewayInput, entities.Gateway](in)
	if err != nil {
		return nil, err
	}
	updated, err := s.deps.Repo.Update(ctx, id, entity)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return nil, notFound()
		}
		return nil, err
	}
	return mapper.EntityToDto[entities.Gateway, dtos.Gateway](updated)
}

// Delete removes a gateway by id and reports the outcome.
func (s *GatewaysService) Delete(ctx context.Context, id string) (map[string]bool, error) {
	if err := s.deps.Repo.Delete(ctx, id); err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return nil, notFound()
		}
		return nil, err
	}
	return map[string]bool{"success": true}, nil
}
