// Package steps holds saga steps that stand up the in-process targets a device
// fires against (an HTTP echo, an MQTT broker). They are modelled as steps so
// the fixture's lifecycle lives inside the saga chain: Do starts it and
// publishes its coordinates to the bag, Compensate shuts it down during
// rollback.
package steps

import (
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/common/utils"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/core/saga"
)

const (
	// Echo target, written by StartEcho.
	BagKeyEchoServer = "targets.echoServer" // *utils.Echo
	BagKeyEchoURL    = "targets.echoUrl"

	// MQTT broker, written by StartMQTTBroker.
	BagKeyMQTTBroker        = "targets.mqttBroker" // *utils.MQTTBroker
	BagKeyMQTTPlainURL      = "targets.mqttPlainUrl"
	BagKeyMQTTTLSURL        = "targets.mqttTlsUrl"
	BagKeyMQTTUsername      = "targets.mqttUsername"
	BagKeyMQTTPassword      = "targets.mqttPassword"
	BagKeyMQTTCAPem         = "targets.mqttCaPem"
	BagKeyMQTTClientCertPem = "targets.mqttClientCertPem"
	BagKeyMQTTClientKeyPem  = "targets.mqttClientKeyPem"

	// BagKeyDownlinkPayload is the payload PublishDownlink injected, so the
	// downlink assert can match it in the logs.
	BagKeyDownlinkPayload = "targets.downlinkPayload"

	// BagKeyConsoleStream holds the connected *utils.ConsoleStream the console
	// asserts read live frames from.
	BagKeyConsoleStream = "targets.consoleStream"
)

// ConsoleStreamFromBag fetches the console stream published by StartConsoleStream,
// failing the test fast when it is missing.
func ConsoleStreamFromBag(c *saga.Context) *utils.ConsoleStream {
	v, ok := c.Get(BagKeyConsoleStream)
	if !ok {
		c.T.Fatalf("[SAGA] missing console stream in bag (StartConsoleStream did not run?)")
	}
	s, ok := v.(*utils.ConsoleStream)
	if !ok {
		c.T.Fatalf("[SAGA] bag key %q is not a *utils.ConsoleStream (got %T)", BagKeyConsoleStream, v)
	}
	return s
}

// brokerUser and brokerPass are the credentials the broker's plain listener
// accepts and the username/password device presents.
const (
	brokerUser = "e2e-user"
	brokerPass = "e2e-pass"
)

// BrokerFromBag fetches the MQTT broker handle published by StartMQTTBroker,
// failing the test fast when it is missing.
func BrokerFromBag(c *saga.Context) *utils.MQTTBroker {
	v, ok := c.Get(BagKeyMQTTBroker)
	if !ok {
		c.T.Fatalf("[SAGA] missing MQTT broker in bag (StartMQTTBroker did not run?)")
	}
	b, ok := v.(*utils.MQTTBroker)
	if !ok {
		c.T.Fatalf("[SAGA] bag key %q is not a *utils.MQTTBroker (got %T)", BagKeyMQTTBroker, v)
	}
	return b
}
