package session

import "simulator/service/src/modules/engine/application/ports"

// Compile-time proof the registry satisfies its port.
var _ ports.ConnectorRegistry = (*connectorRegistry)(nil)

// NewConnectorRegistry builds the connector registry with the session-capable
// protocols: mqtt (per-device broker session) and lorawan/basicstation (a shared
// gateway link). One LoRaWAN connector serves both lorawan and basicstation ids.
func NewConnectorRegistry() ports.ConnectorRegistry {
	r := &connectorRegistry{byProtocol: make(map[string]ports.Connector)}
	r.add(NewMQTT())
	lora := NewLoRaWAN()
	r.byProtocol["lorawan"] = lora
	r.byProtocol["basicstation"] = lora
	return r
}

// add registers a connector under its protocol id.
func (r *connectorRegistry) add(c ports.Connector) {
	r.byProtocol[c.Protocol()] = c
}

// Connector resolves the connector for a protocol id.
func (r *connectorRegistry) Connector(protocol string) (ports.Connector, bool) {
	c, ok := r.byProtocol[protocol]
	return c, ok
}
