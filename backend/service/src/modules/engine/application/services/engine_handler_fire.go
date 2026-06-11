package services

import (
	"context"
	"encoding/json"

	devicescontract "simulator/packages/contracts/devices"
	enginePorts "simulator/service/src/modules/engine/application/ports"
	"simulator/service/src/modules/engine/domain/entities"
	domainsvc "simulator/service/src/modules/engine/domain/services"
)

// findDevice loads a device by its server id.
func (s *EngineService) findDevice(ctx context.Context, deviceID string) (*devicescontract.Device, error) {
	devices, err := s.deps.Devices.List(ctx)
	if err != nil {
		return nil, err
	}
	for i := range devices {
		if devices[i].ID == deviceID {
			return &devices[i], nil
		}
	}
	return nil, enginePorts.ErrDeviceNotFound
}

// resolveFireEvent picks the pre-registered event by id, or parses the inline
// ad-hoc event when no id is given.
func resolveFireEvent(dev devicescontract.Device, in enginePorts.FireInput) (entities.DeviceEvent, error) {
	if in.EventID != "" {
		events, err := domainsvc.ParseEvents(dev.Events)
		if err != nil {
			return entities.DeviceEvent{}, err
		}
		for _, e := range events {
			if e.ID == in.EventID {
				return e, nil
			}
		}
		return entities.DeviceEvent{}, enginePorts.ErrEventNotFound
	}
	if len(in.Event) > 0 {
		var e entities.DeviceEvent
		if err := json.Unmarshal(in.Event, &e); err != nil {
			return entities.DeviceEvent{}, err
		}
		return e, nil
	}
	return entities.DeviceEvent{}, enginePorts.ErrEventNotFound
}
