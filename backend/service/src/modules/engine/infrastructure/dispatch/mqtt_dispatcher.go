package dispatch

import (
	"context"
	"fmt"

	mqttclient "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mqttclient"

	"simulator/service/src/modules/engine/application/ports"
)

// Compile-time proof the adapter satisfies the Dispatcher port.
var _ ports.Dispatcher = (*mqttDispatcher)(nil)

// NewMQTT builds the MQTT dispatcher.
func NewMQTT() ports.Dispatcher {
	return &mqttDispatcher{}
}

// Protocol identifies this dispatcher in the registry.
func (d *mqttDispatcher) Protocol() string { return "mqtt" }

// Dispatch connects to the broker, publishes the rendered payload, and
// disconnects. Connecting per send is the simple v1; a per-broker connection
// cache is a later optimization.
func (d *mqttDispatcher) Dispatch(ctx context.Context, req ports.DispatchRequest) ports.DispatchResult {
	cli, err := mqttclient.New(mqttclient.Config{
		BrokerURL:      req.BrokerURL,
		ClientID:       req.ClientID,
		Username:       req.Username,
		Password:       req.Password,
		ConnectTimeout: mqttConnectTimeout,
	})
	if err != nil {
		return ports.DispatchResult{Err: err}
	}
	if err := cli.Connect(ctx); err != nil {
		return ports.DispatchResult{Err: err}
	}
	defer cli.Disconnect(mqttQuiesceMillis)

	if err := cli.Publish(ctx, req.Topic, req.QoS, req.Retain, []byte(req.Payload)); err != nil {
		return ports.DispatchResult{Err: err}
	}
	return ports.DispatchResult{OK: true, Status: fmt.Sprintf("qos%d", req.QoS)}
}
