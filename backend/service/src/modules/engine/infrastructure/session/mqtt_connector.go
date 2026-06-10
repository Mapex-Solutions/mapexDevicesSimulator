package session

import (
	"context"
	"fmt"

	mqttclient "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mqttclient"

	"simulator/service/src/modules/engine/application/ports"
)

// Compile-time proofs the adapter satisfies its ports.
var (
	_ ports.Connector = (*mqttConnector)(nil)
	_ ports.Session   = (*mqttSession)(nil)
)

// NewMQTT builds the MQTT connector.
func NewMQTT() ports.Connector { return &mqttConnector{} }

// Protocol identifies this connector in the registry.
func (c *mqttConnector) Protocol() string { return "mqtt" }

// Open connects to the broker and, when the device has downlink subscriptions,
// subscribes to each topic and wires received messages to the inbound sink. The
// returned session keeps the connection open for uplinks. The session manager owns
// the reconnect loop; Open is a single attempt.
func (c *mqttConnector) Open(ctx context.Context, spec ports.SessionSpec, in ports.InboundSink, status ports.StatusSink) (ports.Session, error) {
	cli, err := mqttclient.New(mqttclient.Config{
		BrokerURL: spec.BrokerURL,
		ClientID:  spec.ClientID,
		Username:  spec.Username,
		Password:  spec.Password,
	})
	if err != nil {
		return nil, err
	}
	if err := cli.Connect(ctx); err != nil {
		return nil, err
	}

	for _, sub := range spec.Subscriptions {
		if sub.Topic == "" {
			continue
		}
		qos := sub.QoS
		status("subscribing", sub.Topic)
		if err := cli.Subscribe(ctx, sub.Topic, qos, func(topic string, payload []byte) {
			in(ports.InboundMessage{
				Topic:   topic,
				Payload: string(payload),
				Status:  fmt.Sprintf("received qos%d", qos),
				Summary: topic,
			})
		}); err != nil {
			cli.Disconnect(mqttQuiesceMillis)
			return nil, err
		}
		status("subscribed", sub.Topic)
	}

	return &mqttSession{client: cli}, nil
}

// Send publishes an uplink through the live connection.
func (s *mqttSession) Send(ctx context.Context, msg ports.OutboundMessage) ports.SendResult {
	if err := s.client.Publish(ctx, msg.Topic, msg.QoS, msg.Retain, []byte(msg.Payload)); err != nil {
		return ports.SendResult{Err: err}
	}
	return ports.SendResult{OK: true, Status: fmt.Sprintf("qos%d", msg.QoS)}
}

// Close tears down the connection.
func (s *mqttSession) Close() error {
	s.client.Disconnect(mqttQuiesceMillis)
	return nil
}

// Connected reports whether the underlying client is currently connected.
func (s *mqttSession) Connected() bool { return s.client.IsConnected() }
