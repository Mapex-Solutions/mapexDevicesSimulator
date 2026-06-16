// Package steps holds saga steps that exercise the simulator's engine fire
// endpoint (POST /api/devices/{id}/fire).
package steps

// fireEventID is the canonical pre-registered event id every device fixture
// carries, so the fire step needs no parameter.
const fireEventID = "e1"

// BagKeyFiredAt is the time the fire returned, available to asserts as the lower
// bound of the "logs produced by this action" search window.
const BagKeyFiredAt = "engine.firedAt"
