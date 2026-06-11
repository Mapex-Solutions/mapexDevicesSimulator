package services

import (
	"context"
	"errors"

	customErrors "github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	status "github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
	"github.com/Mapex-Solutions/mapexGoKit/utils/mapper"

	"simulator/service/src/modules/devices/application/di"
	"simulator/service/src/modules/devices/application/dtos"
	"simulator/service/src/modules/devices/application/ports"
	"simulator/service/src/modules/devices/domain/entities"
	"simulator/service/src/modules/devices/domain/repositories"
)

// notFound is the not-found error carrying the 404 code the global error handler
// renders as an envelope. It translates the domain's repository sentinel into the
// HTTP-aware application error.
func notFound() error {
	return &customErrors.ServerCustomError{Code: status.NOT_FOUND, Errors: []string{"device not found"}}
}

// Compile-time check that the service satisfies its port.
var _ ports.DevicesServicePort = (*DevicesService)(nil)

// New builds the devices service over its injected repository port.
func New(deps di.DevicesServiceDI) ports.DevicesServicePort {
	return &DevicesService{deps: deps}
}

// List returns every stored device mapped to its wire DTO.
func (s *DevicesService) List(ctx context.Context) ([]dtos.Device, error) {
	stored, err := s.deps.Repo.List(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]dtos.Device, 0, len(stored))
	for i := range stored {
		dto, err := mapper.EntityToDto[entities.Device, dtos.Device](&stored[i])
		if err != nil {
			return nil, err
		}
		out = append(out, *dto)
	}
	return out, nil
}

// Create maps the body to an entity, persists it, and returns the stored DTO.
func (s *DevicesService) Create(ctx context.Context, in *dtos.DeviceInput) (*dtos.Device, error) {
	entity, err := mapper.DtoToEntity[dtos.DeviceInput, entities.Device](in)
	if err != nil {
		return nil, err
	}
	created, err := s.deps.Repo.Create(ctx, entity)
	if err != nil {
		return nil, err
	}
	s.deps.Signal.Notify()
	return mapper.EntityToDto[entities.Device, dtos.Device](created)
}

// Update maps the body to an entity, replaces the device, and returns the DTO.
func (s *DevicesService) Update(ctx context.Context, id string, in *dtos.DeviceInput) (*dtos.Device, error) {
	entity, err := mapper.DtoToEntity[dtos.DeviceInput, entities.Device](in)
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
	s.deps.Signal.Notify()
	return mapper.EntityToDto[entities.Device, dtos.Device](updated)
}

// Delete removes a device by id and reports the outcome.
func (s *DevicesService) Delete(ctx context.Context, id string) (map[string]bool, error) {
	if err := s.deps.Repo.Delete(ctx, id); err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return nil, notFound()
		}
		return nil, err
	}
	s.deps.Signal.Notify()
	return map[string]bool{"success": true}, nil
}
