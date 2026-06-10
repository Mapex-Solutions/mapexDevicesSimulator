package ports

import (
	"context"
	"encoding/json"
	"errors"
)

// FireInput is an on-demand fire request: either a pre-registered event id, or an
// inline ad-hoc event (the "Generic" path) shaped like a device event.
type FireInput struct {
	EventID string          // fire a pre-registered event by id
	Event   json.RawMessage // ad-hoc event JSON (DeviceEvent shape); used when EventID == ""
}

// Fire error sentinels, mapped to HTTP status by the interface layer.
var (
	ErrDeviceNotFound  = errors.New("engine: device not found")
	ErrEventNotFound   = errors.New("engine: event not found")
	ErrFireUnsupported = errors.New("engine: device/event cannot be fired")
)

// EnginePort is the simulation engine's lifecycle + control surface. OnMount
// reads the devices and starts the scheduler (fired during module init, not in
// main); Reconcile re-reads and re-aligns the running jobs; Fire sends one event
// on demand (the console/REST "fire" action); OnShutdown stops everything.
type EnginePort interface {
	OnMount()
	Reconcile()
	Fire(ctx context.Context, deviceID string, in FireInput) error
	OnShutdown(ctx context.Context) error
}
