package session

import (
	mqttclient "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mqttclient"

	"simulator/service/src/modules/engine/application/ports"
)

// connectorRegistry maps a protocol id to its connector. A protocol absent from
// the map is not session-capable and falls back to the one-shot dispatcher.
type connectorRegistry struct {
	byProtocol map[string]ports.Connector
}

// mqttConnector opens persistent MQTT sessions (one live client per device).
type mqttConnector struct{}

// mqttSession is a live MQTT connection for one device: it publishes uplinks and
// receives downlinks on the same open client.
type mqttSession struct {
	client *mqttclient.Client
}
