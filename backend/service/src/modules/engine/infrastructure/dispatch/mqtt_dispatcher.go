package dispatch

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"strconv"
	"sync/atomic"

	mqttclient "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mqttclient"

	"simulator/service/src/modules/engine/application/ports"
)

// Compile-time proof the adapter satisfies the Dispatcher port.
var _ ports.Dispatcher = (*mqttDispatcher)(nil)

// oneShotSeq makes each one-shot connection's client id unique within the
// process. A one-shot send can happen while the device's persistent session is
// still connecting with the SAME configured client id; MQTT brokers enforce a
// unique client id per connection and kick the older one, so without a distinct
// id the two would fight and the one-shot publish would be lost.
var oneShotSeq atomic.Uint64

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
	mc := mqttclient.Config{
		BrokerURL:      req.BrokerURL,
		ClientID:       oneShotClientID(req.ClientID),
		Username:       req.Username,
		Password:       req.Password,
		ConnectTimeout: mqttConnectTimeout,
	}
	if req.TLSCert != "" || req.TLSCa != "" {
		tlsCfg, err := buildTLSConfig(req.TLSCert, req.TLSKey, req.TLSCa)
		if err != nil {
			return ports.DispatchResult{Err: err}
		}
		mc.TLSConfig = tlsCfg
	}
	cli, err := mqttclient.New(mc)
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

// oneShotClientID derives a per-send client id from the device's configured one
// so a one-shot connection never collides with the device's persistent session
// (or another concurrent one-shot). An empty base is left empty so the client
// auto-generates a random id.
func oneShotClientID(base string) string {
	if base == "" {
		return ""
	}
	return base + "-oneshot-" + strconv.FormatUint(oneShotSeq.Add(1), 10)
}

// buildTLSConfig assembles the tls.Config for a certificate-authenticated MQTT
// broker from the device's PEM material: the client keypair it presents and the
// CA it trusts for the broker. Either may be empty (cert-only or CA-only), and
// the result is paired with an ssl:// broker URL by the caller.
func buildTLSConfig(certPem, keyPem, caPem string) (*tls.Config, error) {
	cfg := &tls.Config{MinVersion: tls.VersionTLS12}
	if certPem != "" && keyPem != "" {
		cert, err := tls.X509KeyPair([]byte(certPem), []byte(keyPem))
		if err != nil {
			return nil, fmt.Errorf("mqtt tls: client keypair: %w", err)
		}
		cfg.Certificates = []tls.Certificate{cert}
	}
	if caPem != "" {
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM([]byte(caPem)) {
			return nil, fmt.Errorf("mqtt tls: invalid CA pem")
		}
		cfg.RootCAs = pool
	}
	return cfg, nil
}
