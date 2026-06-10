package services

import "simulator/service/src/modules/gateways/application/di"

// GatewaysService orchestrates gateway CRUD over the repository port, mapping
// between the wire DTOs and the domain entity.
type GatewaysService struct {
	deps di.GatewaysServiceDI
}
