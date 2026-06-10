package entities

import (
	"encoding/json"
	"time"
)

// Device is the persisted simulated device. Scalar columns are the queryable
// fields; the protocol config, attributes and events are nested and stored as
// JSON text (the `,json` tag on the sqlite model serializes them).
type Device struct {
	ID         string            `db:"id,pk"`
	Name       string            `db:"name"`
	DeviceID   string            `db:"device_id"`
	ProtocolID string            `db:"protocol_id"`
	Enabled    bool              `db:"enabled"`
	StoreLogs  bool              `db:"store_logs"`
	Config     json.RawMessage   `db:"config,json"`
	Attributes map[string]string `db:"attributes,json"`
	Events     json.RawMessage   `db:"events,json"`
	Created    time.Time         `db:"created"`
}

// GetCreated exposes the creation time for the entity-to-DTO mapper.
func (d *Device) GetCreated() time.Time { return d.Created }
