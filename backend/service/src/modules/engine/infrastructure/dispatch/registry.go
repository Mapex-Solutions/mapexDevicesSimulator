package dispatch

import (
	"simulator/service/src/modules/engine/application/ports"
)

// Compile-time proof the registry satisfies its port.
var _ ports.Registry = (*registry)(nil)

// NewRegistry builds the dispatcher registry with the protocols enabled today
// (http, mqtt). LoRaWAN/Basics Station ship in a later phase.
func NewRegistry() ports.Registry {
	r := &registry{byProtocol: make(map[string]ports.Dispatcher)}
	r.add(NewHTTP())
	r.add(NewMQTT())
	return r
}

// add registers a dispatcher under its protocol id.
func (r *registry) add(d ports.Dispatcher) {
	r.byProtocol[d.Protocol()] = d
}

// For resolves the dispatcher for a protocol id.
func (r *registry) For(protocol string) (ports.Dispatcher, bool) {
	d, ok := r.byProtocol[protocol]
	return d, ok
}
