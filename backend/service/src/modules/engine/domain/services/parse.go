package services

import (
	"encoding/json"

	"simulator/service/src/modules/engine/domain/entities"
)

// ParseEvents unmarshals a device's events JSON into structured events. An empty
// payload yields no events (a device without registered events).
func ParseEvents(raw json.RawMessage) ([]entities.DeviceEvent, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	var events []entities.DeviceEvent
	if err := json.Unmarshal(raw, &events); err != nil {
		return nil, err
	}
	return events, nil
}
