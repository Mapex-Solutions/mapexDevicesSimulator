package devices

import (
	"encoding/json"

	timeUtil "github.com/Mapex-Solutions/mapexGoKit/utils/time"
)

// Device is the wire shape returned for a simulated device. `created` serializes
// to an ISO-8601 string via NullTime; `config` and `events` are passed through as
// raw JSON (their shape varies by protocol and the engine, not this module,
// interprets them).
type Device struct {
	ID         string             `json:"id"`
	Created    *timeUtil.NullTime `json:"created"`
	Name       string             `json:"name"`
	DeviceID   string             `json:"deviceId"`
	ProtocolID string             `json:"protocolId"`
	Enabled    bool               `json:"enabled"`
	StoreLogs  bool               `json:"storeLogs"`
	Config     json.RawMessage    `json:"config"`
	Attributes map[string]string  `json:"attributes"`
	Events     json.RawMessage    `json:"events"`
}

// SetCreated lets mapper.EntityToDto populate the creation timestamp from the
// entity's GetCreated().
func (d *Device) SetCreated(t *timeUtil.NullTime) { d.Created = t }

// DeviceIDParam binds the :id path parameter on the device by-id routes.
type DeviceIDParam struct {
	ID string `params:"id" validate:"required"`
}

// DeviceInput is the create/update body. id and created are server-assigned and
// not part of it.
type DeviceInput struct {
	Name       string            `json:"name" validate:"required"`
	DeviceID   string            `json:"deviceId" validate:"required"`
	ProtocolID string            `json:"protocolId" validate:"required,oneof=http mqtt lorawan basicstation"`
	Enabled    bool              `json:"enabled"`
	StoreLogs  bool              `json:"storeLogs"`
	Config     json.RawMessage   `json:"config"`
	Attributes map[string]string `json:"attributes"`
	Events     json.RawMessage   `json:"events"`
}
