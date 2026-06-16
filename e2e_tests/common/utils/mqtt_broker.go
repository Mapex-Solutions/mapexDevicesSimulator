package utils

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log/slog"
	"sync"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/mochi-mqtt/server/v2/packets"
)

// MQTTPublish is one message the broker accepted from a device.
type MQTTPublish struct {
	ClientID string
	Topic    string
	Payload  string
}

// MQTTBroker is an in-process MQTT broker (mochi-mqtt) the journeys point a
// device at. It exposes two listeners so one fixture covers both auth modes:
//   - a plain TCP listener that authenticates by username/password;
//   - a TLS listener that requires and verifies a client certificate.
//
// Every accepted publish is captured so an assert can prove the device
// connected, authenticated, and delivered — the whole point of the auth journey.
type MQTTBroker struct {
	server    *mqtt.Server
	plainAddr string
	tlsAddr   string
	tls       TLSMaterial

	username string
	password string

	mu        sync.Mutex
	published []MQTTPublish
}

// MQTTCreds is the username/password the plain listener accepts.
type MQTTCreds struct {
	Username string
	Password string
}

// StartMQTTBroker boots the broker with both listeners on random loopback ports.
// creds are the credentials the plain listener accepts; the TLS listener trusts
// the generated CA and verifies the device's client cert against it. The caller
// stops it via t.Cleanup(broker.Close).
func StartMQTTBroker(creds MQTTCreds) (*MQTTBroker, error) {
	material, err := GenerateTLSMaterial()
	if err != nil {
		return nil, fmt.Errorf("mqtt broker tls material: %w", err)
	}

	b := &MQTTBroker{
		tls:      material,
		username: creds.Username,
		password: creds.Password,
	}

	b.server = mqtt.New(&mqtt.Options{
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		// InlineClient lets the fixture publish downlinks to subscribed devices
		// (see Publish), so the inbound path can be exercised.
		InlineClient: true,
	})
	if err := b.server.AddHook(&brokerHook{broker: b}, nil); err != nil {
		return nil, fmt.Errorf("mqtt broker hook: %w", err)
	}

	plain := listeners.NewTCP(listeners.Config{ID: "e2e-plain", Address: "127.0.0.1:0"})
	if err := b.server.AddListener(plain); err != nil {
		return nil, fmt.Errorf("mqtt broker plain listener: %w", err)
	}

	tlsCfg, err := b.serverTLSConfig()
	if err != nil {
		return nil, err
	}
	secure := listeners.NewTCP(listeners.Config{ID: "e2e-tls", Address: "127.0.0.1:0", TLSConfig: tlsCfg})
	if err := b.server.AddListener(secure); err != nil {
		return nil, fmt.Errorf("mqtt broker tls listener: %w", err)
	}

	if err := b.server.Serve(); err != nil {
		return nil, fmt.Errorf("mqtt broker serve: %w", err)
	}

	b.plainAddr = plain.Address()
	b.tlsAddr = secure.Address()
	return b, nil
}

// serverTLSConfig builds the broker's mutual-TLS config: it presents the server
// cert and requires a client cert signed by the run's CA.
func (b *MQTTBroker) serverTLSConfig() (*tls.Config, error) {
	cert, err := tls.X509KeyPair([]byte(b.tls.ServerCertPEM), []byte(b.tls.ServerKeyPEM))
	if err != nil {
		return nil, fmt.Errorf("mqtt broker server keypair: %w", err)
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM([]byte(b.tls.CAPEM)) {
		return nil, fmt.Errorf("mqtt broker: invalid CA pem")
	}
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    pool,
		MinVersion:   tls.VersionTLS12,
	}, nil
}

// PlainURL is the tcp:// broker URL the username/password device connects to.
func (b *MQTTBroker) PlainURL() string { return "tcp://" + b.plainAddr }

// TLSURL is the ssl:// broker URL the certificate device connects to.
func (b *MQTTBroker) TLSURL() string { return "ssl://" + b.tlsAddr }

// CAPEM is the CA the device trusts to verify the broker.
func (b *MQTTBroker) CAPEM() string { return b.tls.CAPEM }

// ClientCertPEM is the client cert the certificate device presents.
func (b *MQTTBroker) ClientCertPEM() string { return b.tls.ClientCertPEM }

// ClientKeyPEM is the client key paired with ClientCertPEM.
func (b *MQTTBroker) ClientKeyPEM() string { return b.tls.ClientKeyPEM }

// Published returns a copy of every message the broker has accepted so far.
func (b *MQTTBroker) Published() []MQTTPublish {
	b.mu.Lock()
	defer b.mu.Unlock()
	out := make([]MQTTPublish, len(b.published))
	copy(out, b.published)
	return out
}

// Publish injects a message on topic as if an external party sent it, so a
// subscribed simulator device receives it as a downlink. Publish it retained so
// it is delivered even if the device's subscription is not active yet (the
// session connects and subscribes asynchronously), removing the timing race.
func (b *MQTTBroker) Publish(topic string, payload []byte, retain bool, qos byte) error {
	return b.server.Publish(topic, payload, retain, qos)
}

// Close shuts the broker down.
func (b *MQTTBroker) Close() {
	if b.server != nil {
		_ = b.server.Close()
	}
}

// record appends an accepted publish under the lock.
func (b *MQTTBroker) record(p MQTTPublish) {
	b.mu.Lock()
	b.published = append(b.published, p)
	b.mu.Unlock()
}

// authenticate accepts a connection when it presented a verified client cert
// (the TLS listener already enforced that at the handshake, so username is
// empty) or when the username/password match the plain listener's credentials.
func (b *MQTTBroker) authenticate(username, password string) bool {
	if username == "" {
		return true
	}
	return username == b.username && password == b.password
}

// brokerHook bridges mochi-mqtt callbacks to the MQTTBroker: it authenticates
// connections, permits all topics, and captures accepted publishes.
type brokerHook struct {
	mqtt.HookBase
	broker *MQTTBroker
}

// ID names the hook in mochi-mqtt's registry.
func (h *brokerHook) ID() string { return "e2e-broker-hook" }

// Provides advertises the callbacks this hook implements.
func (h *brokerHook) Provides(b byte) bool {
	switch b {
	case mqtt.OnConnectAuthenticate, mqtt.OnACLCheck, mqtt.OnPublished:
		return true
	default:
		return false
	}
}

// OnConnectAuthenticate gates the CONNECT against the broker's credentials.
func (h *brokerHook) OnConnectAuthenticate(cl *mqtt.Client, pk packets.Packet) bool {
	return h.broker.authenticate(string(pk.Connect.Username), string(pk.Connect.Password))
}

// OnACLCheck permits publish/subscribe on every topic — this is a test broker.
func (h *brokerHook) OnACLCheck(cl *mqtt.Client, topic string, write bool) bool { return true }

// OnPublished captures each accepted publish for the journey's assert.
func (h *brokerHook) OnPublished(cl *mqtt.Client, pk packets.Packet) {
	h.broker.record(MQTTPublish{ClientID: cl.ID, Topic: pk.TopicName, Payload: string(pk.Payload)})
}
