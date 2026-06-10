package services

import "simulator/service/src/modules/devices/application/di"

// DevicesService orchestrates device CRUD over the repository port, mapping
// between the wire DTOs and the domain entity.
type DevicesService struct {
	deps di.DevicesServiceDI
}
