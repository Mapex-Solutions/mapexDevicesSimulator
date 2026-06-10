# Device protocols, pre-registered events & firing

> Status: APPROVED — building per protocol, HTTP first.
> Source of truth for the add-device / add-gateway screens and the fire-event flow.

## 1. Purpose

Complete the protocols by making device registration protocol-aware and giving each
device **pre-registered events** that can be fired on demand. Firing is per device
(the lightning button in the console device list); the form adapts to the device's
protocol, and values can be fixed or random via template placeholders.

## 2. Device model

```
Device = {
  id, name,
  protocolId,                 // http | mqtt | lorawan | basicstation
  config: ProtocolConfig,     // target/credentials, per protocol
  attributes: Record<string,string>,
  events: DeviceEvent[],      // pre-registered, per protocol
}
```

### 2.1 Target / credentials (`config`) per protocol

- **HTTP**: `url`, `method`, `authMode` (none|apiKey|basic), `apiKeyHeader`, `apiKey`,
  `basicUser`, `basicPass`, `contentType`. (already implemented)
- **MQTT**: `brokerUrl`, `clientId`, `authMode` (userpass|cert), `username`, `password`,
  `clientCertPem`, `clientKeyPem`, `caPem`, `baseTopic`.
- **LoRaWAN**: `devEui`, `joinEui`, `appKey`, `nwkKey` (1.1), `macVersion` (1.0.x|1.1),
  `region` (e.g. EU868), `activation` (otaa|abp), `devAddr` + session keys (abp),
  `gatewayId` (the gateway it transmits through).

### 2.2 Pre-registered events (`events`) per protocol

```
DeviceEvent = { id, name, ...protocol-specific fields..., payload }
```

- **HTTP event**: `method`, `path` (+ query), `query` (k/v), `headers` (k/v), `body` (template).
- **MQTT event**: `topic`, `qos` (0|1|2), `retained`, `payload` (template).
- **LoRaWAN event**: `fPort`, `confirmed`, `payload` (template; hex or json).

## 3. Fixed vs random values — template placeholders

Payloads (and templatable values) are **text templates** resolved at send time (engine)
and for preview (client). Helpers:

- `{{randInt(min,max)}}` — random integer in range
- `{{randFloat(min,max)}}` — random float in range
- `{{now}}` — current ISO timestamp · `{{nowMs}}` — epoch ms
- `{{deviceId}}`, `{{deviceName}}` — device identity
- `{{counter}}` — incrementing counter per send
- `{{uuid}}` — random uuid
- A fixed value is just the literal text.

The fire dialog shows a small **hint palette** of these helpers. A client-side renderer
provides a live preview; the engine is the source of truth at send time.

## 4. Fire-event flow

Trigger: the lightning button on a device row in the console (or device list).

```
┌── Fire event · <device> (<protocol>) ──────────┐
│ Event: ( pre-registered ▼ )  |  ( Generic )     │
│   pre-registered → fills the form               │
│   Generic        → blank protocol fields + hints│
│ <protocol-specific fields, templatable>         │
│ Preview: <rendered payload>      [ Send now ]   │
└────────────────────────────────────────────────┘
```

On send: resolve templates → (engine) deliver via the device's `config` target →
record on the console stream (down/uplink) and in the Logs history. While the engine
is offline, the message is recorded locally so the flow is visible.

## 5. Gateway model (LoRaWAN)

```
Gateway = { id, name, eui, region/frequencyPlanId, transport (basicstation|semtech-udp),
            wsUrl | udpAddr, lnsUrl }
```
No events. Gateways carry sensors; a LoRaWAN device references a `gatewayId`.

## 6. Extensibility — ProtocolRegistry

Each protocol entry contributes (additive; new protocol = new entry):

- `configComponent` + `defaultConfig()` + `validate()` — the target/credentials form.
- `eventFields` (or `eventComponent`) + `defaultEvent()` — the event editor + fire form.
- `icon`, `labelKey`, `enabled`.

The add-device dialog and the fire dialog render whatever the active protocol exposes,
so neither has hardcoded protocol knowledge.

## 7. Build order (vertical slices)

1. **HTTP** — device `config` (have) + `events` (method/path/query/headers/body) +
   protocol-aware fire dialog (event picker + Generic) + template placeholders + preview.
2. **MQTT** — device `config` (broker/clientId/auth/baseTopic) + `events` (topic/qos/payload).
3. **LoRaWAN + Gateway** — device LoRa `config` (keys/version/region/otaa-abp/gateway) +
   `events` (fPort/payload), and the add-gateway screen.

Each slice is shippable; the registry + fire dialog are written once and grow by data.
