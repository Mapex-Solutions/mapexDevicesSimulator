package ports

import "context"

// EnginePort is the simulation engine's lifecycle + control surface. OnMount
// reads the devices and starts the scheduler (fired during module init, not in
// main); Reconcile re-reads and re-aligns the running jobs; OnShutdown stops
// everything gracefully.
type EnginePort interface {
	OnMount()
	Reconcile()
	OnShutdown(ctx context.Context) error
}
