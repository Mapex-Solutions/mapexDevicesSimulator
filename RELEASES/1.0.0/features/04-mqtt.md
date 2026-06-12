# 04 — MQTT protocol

An MQTT device keeps a **live connection** to a broker: it publishes uplinks and,
when receiving is enabled, stays subscribed to topics and streams every message it
gets back as a downlink. It is the simplest protocol that can both send and receive.

## Capabilities

- **Live broker session** — connects and holds the connection open (reconnects
  forever with backoff if the broker drops).
- **Publish uplinks** through the open client; **QoS 0/1/2** and **retain**.
- **Receive (downlinks)** — with receive on, the device subscribes to a list of
  topics and emits each received message as a `down` frame.
- **Auth** — none, **username/password**, or **TLS client certificate** (cert/key/CA
  PEM).
- **Base topic** — a per-device prefix that the engine prepends to every event and
  subscription topic.

## How it works

A persistent MQTT client per device lives in `infrastructure/session/mqtt_connector`
(over the GoKit `mqttclient`). `Open` connects, subscribes to each configured topic
and wires the handler to the engine's `InboundSink` (→ `down` frames); `Send`
publishes through the same live client. While the session is still connecting or
reconnecting, a fire falls back to the one-shot dispatcher so it is never dropped.

Config (`kind: "mqtt"`): `brokerUrl` (`tcp://host:1883`), `clientId`, `baseTopic`,
`authMode` (`none|userpass|tls`), `username`/`password`, `tlsCertPem`/`tlsKeyPem`/
`tlsCaPem`, `receiveEnabled`, `subscriptions` (`[{name, topic, qos}]`). Event (`mqtt`):
`topic`, `qos`, `retain`, body. The lifecycle (connecting → connected → subscribed,
reconnecting, disconnected) streams to the [console](./09-console.md).

## Notes

- **Topics are relative to `baseTopic`** — the engine prepends it to both publish and
  subscription topics, so an event topic `sensor/telemetry` under base `mapex/quick`
  publishes to `mapex/quick/sensor/telemetry`.
- Receive off → publish-only; the subscription list is ignored.
- The session reconnects with a bounded backoff (floor / ceiling / jitter), every
  transition shown on the console.

---
> Part of the [MapexOS ecosystem](../README.md#part-of-the-mapexos-ecosystem).
