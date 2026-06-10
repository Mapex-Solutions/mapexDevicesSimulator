# Simulator schemas — full REST + WebSocket contract

Source-of-truth payload shapes the Go backend must accept/return, taken verbatim
from the frontend's typed client (`frontend/src/services/sim/interfaces`). Use
these to derive the Go DTOs (and any Zod/validation layer). Types are shown in
TypeScript notation; JSON examples follow each.

## Response envelope (Mapex standard)

Every response is wrapped in the standard Mapex envelope. The type listed in the
"Returns (in `data`)" column is what lands in `data`; on error, `errors` carries
the messages and `data` is null. Request bodies are the raw `*Input` (not wrapped).

```json
{ "status": 200, "errors": null, "data": <payload> }
```

REST endpoints these cover:

| Method | Path                 | Body          | Returns (in `data`) |
|--------|----------------------|---------------|---------------------|
| GET    | `/api/logs`          | — (query)     | `LogPage`           |
| GET    | `/api/devices`       | —             | `Device[]`          |
| POST   | `/api/devices`       | `DeviceInput` | `Device`            |
| PUT    | `/api/devices/:id`   | `DeviceInput` | `Device`            |
| DELETE | `/api/devices/:id`   | —             | `null`              |

| Method | Path                 | Body          | Returns (in `data`) |
|--------|----------------------|---------------|---------------------|
| GET    | `/api/gateways`      | —             | `Gateway[]`         |
| POST   | `/api/gateways`      | `GatewayInput`| `Gateway`           |
| PUT    | `/api/gateways/:id`  | `GatewayInput`| `Gateway`           |
| DELETE | `/api/gateways/:id`  | —             | `null`              |

Conventions:

- All field names are camelCase on the wire (JSON).
- Every response is the envelope `{ status, errors, data }`; DELETE returns it too
  (`200` with `data: null`), not a bare `204`.
- `id` and `created` are **server-assigned** (not part of the `*Input` create
  bodies). `created` is the creation timestamp on EVERY persisted entity, an
  ISO-8601 string. There is no `ts` / `createdAt` / `timestamp` variant.
- `config` (device) is a **discriminated union** on `kind`
  (`http` | `mqtt` | `lorawan` | `basicstation`). The `kind` always matches the
  device's `protocolId`.
- Secrets (`apiKey`, `password`, TLS PEM blocks, LoRaWAN keys) arrive in these
  payloads but must **never be logged** by the engine.

---

## Shared primitives

```ts
type ProtocolId = 'http' | 'mqtt' | 'lorawan' | 'basicstation';

interface KeyValue {
  key: string;
  value: string;
}
```

---

## 1) Log

`GET /api/logs` → `LogPage`. A Log is one persisted device message (the
SQLite-backed history behind the live console stream).

```ts
type LogDirection = 'up' | 'down' | 'system';
type LogKind      = 'data' | 'auth' | 'join' | 'downlink' | 'status';

interface Log {
  id: string;            // server-assigned
  created: string;       // ISO-8601 creation timestamp (server-assigned)
  protocol: ProtocolId;
  deviceId: string;      // the device's user-facing identifier
  deviceName: string;
  direction: LogDirection;
  kind: LogKind;
  summary: string;       // one-line description (e.g. "POST /v1/ingest")
  status?: string;       // optional protocol status (e.g. "200", "qos1", "FCnt 14")
  payload: string;       // raw rendered payload (text / JSON / hex)
}

// Query string for GET /api/logs
interface LogQuery {
  limit: number;         // page size
  offset: number;        // page offset
  protocol?: string;
  kind?: string;
  direction?: string;
  device?: string;       // filter by deviceId
  q?: string;            // free-text search over summary/payload/deviceName
}

// Response
interface LogPage {
  items: Log[];
  total: number;         // total matching rows (for pagination)
}
```

**Example** `GET /api/logs?limit=20&offset=0&protocol=mqtt` (response is the
envelope; `LogPage` is the `data`):

```json
{
  "status": 200,
  "errors": null,
  "data": {
    "items": [
      {
        "id": "log-1042",
        "created": "2026-06-09T02:44:10.000Z",
        "protocol": "mqtt",
        "deviceId": "9f2c1a4e-mqtt-0002",
        "deviceName": "Greenhouse Hub",
        "direction": "up",
        "kind": "data",
        "summary": "PUBLISH greenhouse/.../humidity",
        "status": "qos1",
        "payload": "{ \"humidity\": 63 }"
      }
    ],
    "total": 137
  }
}
```

---

## 2) Create Device

`POST /api/devices` with `DeviceInput` → returns `Device` (`Input` + `id` +
`created`).

```ts
interface DeviceInput {
  name: string;
  deviceId: string;                    // user-facing UUID, drives {{deviceId}}
  protocolId: ProtocolId;
  enabled: boolean;                    // On/Off
  storeLogs: boolean;                  // persist this device's messages to logs
  config: ProtocolConfig;              // discriminated by `kind` (== protocolId)
  attributes: Record<string, string>; // free-form payload variables
  events: DeviceEvent[];               // pre-registered events
}

interface Device extends DeviceInput {
  id: string;       // server-assigned
  created: string;  // ISO-8601
}
```

### config — `ProtocolConfig` (discriminated union on `kind`)

```ts
type ProtocolConfig =
  | HttpConnectionConfig
  | MqttConnectionConfig
  | LoraWanConnectionConfig
  | BasicsStationConnectionConfig;

type HttpMethod   = 'POST' | 'PUT';
type HttpAuthMode = 'none' | 'apiKey' | 'basic';

interface HttpConnectionConfig {
  kind: 'http';
  url: string;
  method: HttpMethod;
  headers: KeyValue[];      // includes Content-Type
  authMode: HttpAuthMode;
  apiKeyHeader: string;     // used when authMode = 'apiKey'
  apiKey: string;           // secret
  basicUser: string;        // used when authMode = 'basic'
  basicPass: string;        // secret
}

type MqttAuthMode = 'none' | 'userpass' | 'tls';

interface MqttConnectionConfig {
  kind: 'mqtt';
  brokerUrl: string;        // mqtt:// or mqtts://
  clientId: string;
  baseTopic: string;        // optional prefix
  authMode: MqttAuthMode;
  username: string;         // used when authMode = 'userpass'
  password: string;         // secret
  tlsCertPem: string;       // used when authMode = 'tls' (secret)
  tlsKeyPem: string;        // secret
  tlsCaPem: string;
}

type LoraWanActivation  = 'otaa' | 'abp';
type LoraWanMacVersion  = '1.0.2' | '1.0.3' | '1.0.4' | '1.1.0';
type GatewayRegion      = 'EU868' | 'US915' | 'AU915' | 'AS923'
                        | 'CN470' | 'IN865' | 'KR920' | 'RU864';

// LoRaWAN node attached to a Gateway (forwards through gateway.link)
interface LoraWanConnectionConfig {
  kind: 'lorawan';
  gatewayId: string;        // references a Gateway.id
  region: GatewayRegion;
  macVersion: LoraWanMacVersion;
  activation: LoraWanActivation;
  // OTAA fields:
  devEui: string;
  joinEui: string;
  appKey: string;           // secret
  nwkKey: string;           // secret, used when macVersion starts with 1.1
  // ABP fields:
  devAddr: string;
  nwkSKey: string;          // secret
  appSKey: string;          // secret
}

// LoRaWAN node carrying its OWN Basics Station link (no separate gateway)
interface BasicsStationConnectionConfig {
  kind: 'basicstation';
  lnsUri: string;           // wss://lns:1887
  gatewayEui: string;       // the embedded gateway's EUI
  region: GatewayRegion;
  macVersion: LoraWanMacVersion;
  activation: LoraWanActivation;
  devEui: string;
  joinEui: string;
  appKey: string;           // secret
  nwkKey: string;           // secret (1.1)
  devAddr: string;
  nwkSKey: string;          // secret
  appSKey: string;          // secret
}
```

### events — `DeviceEvent[]`

An event holds exactly one protocol-specific config matching the device's
protocol (`http`/`mqtt` use the shared body; `lorawan`/`basicstation` use the
LoRaWAN uplink). `schedule` is optional auto-fire.

```ts
type HttpBodyMode = 'none' | 'raw' | 'form';

// Shared body authoring; values may contain {{placeholders}} resolved at send time.
interface RequestBody {
  bodyMode: HttpBodyMode;   // none | raw (JSON text) | form (KeyValue[] -> JSON object)
  bodyFields: KeyValue[];   // used when bodyMode = 'form'
  body: string;             // used when bodyMode = 'raw'
}

interface HttpEventConfig extends RequestBody {
  method: HttpMethod;
  path: string;             // query params authored inline in the path
  headers: KeyValue[];
}

type MqttQoS = 0 | 1 | 2;

interface MqttEventConfig extends RequestBody {
  topic: string;
  qos: MqttQoS;
  retain: boolean;
}

interface LoraWanEventConfig {
  fport: number;            // 1..223
  confirmed: boolean;
  payloadHex: string;       // hex string, may contain {{placeholders}}
}

type EventScheduleUnit = 'seconds' | 'minutes' | 'hours' | 'days';

interface EventSchedule {
  enabled: boolean;
  every: number;            // positive integer
  unit: EventScheduleUnit;
}

interface DeviceEvent {
  id: string;
  name: string;
  http?: HttpEventConfig;       // present for protocolId 'http'
  mqtt?: MqttEventConfig;       // present for protocolId 'mqtt'
  lorawan?: LoraWanEventConfig; // present for protocolId 'lorawan' | 'basicstation'
  schedule?: EventSchedule;
}
```

### Placeholder tokens (resolved by the engine at send time)

Used inside `body`, `bodyFields[].value`, MQTT `topic`, and LoRaWAN `payloadHex`:

```
{{randInt(min,max)}}    integer in range, negatives allowed
{{randFloat(min,max)}}  decimal in range, negatives allowed
{{now}}                 current ISO-8601 timestamp
{{counter}}             auto-incrementing counter
{{deviceId}}            the device's deviceId field
{{uuid}}                random UUID
```

### Example — `POST /api/devices` (HTTP device)

```json
{
  "name": "Edge Sensor 01",
  "deviceId": "9f2c1a4e-http-0001",
  "protocolId": "http",
  "enabled": true,
  "storeLogs": true,
  "config": {
    "kind": "http",
    "url": "http://localhost:5001/v1/ingest",
    "method": "POST",
    "headers": [{ "key": "Content-Type", "value": "application/json" }],
    "authMode": "apiKey",
    "apiKeyHeader": "X-API-Key",
    "apiKey": "demo-key",
    "basicUser": "",
    "basicPass": ""
  },
  "attributes": {},
  "events": [
    {
      "id": "evt-1",
      "name": "Telemetry",
      "http": {
        "method": "POST",
        "path": "/v1/ingest",
        "headers": [],
        "bodyMode": "raw",
        "bodyFields": [],
        "body": "{ \"deviceId\": \"{{deviceId}}\", \"temp\": {{randInt(10,30)}} }"
      },
      "schedule": { "enabled": true, "every": 30, "unit": "seconds" }
    }
  ]
}
```

### Example — `POST /api/devices` (LoRaWAN device, OTAA)

```json
{
  "name": "Field Node LoRa",
  "deviceId": "9f2c1a4e-lora-0003",
  "protocolId": "lorawan",
  "enabled": false,
  "storeLogs": true,
  "config": {
    "kind": "lorawan",
    "gatewayId": "gw-seed-1",
    "region": "EU868",
    "macVersion": "1.0.4",
    "activation": "otaa",
    "devEui": "70B3D57ED0000001",
    "joinEui": "0000000000000000",
    "appKey": "00112233445566778899AABBCCDDEEFF",
    "nwkKey": "",
    "devAddr": "",
    "nwkSKey": "",
    "appSKey": ""
  },
  "attributes": {},
  "events": [
    {
      "id": "evt-1",
      "name": "Uplink",
      "lorawan": { "fport": 10, "confirmed": false, "payloadHex": "01{{randInt(0,255)}}" },
      "schedule": { "enabled": true, "every": 5, "unit": "minutes" }
    }
  ]
}
```

---

## 3) Create Gateway

`POST /api/gateways` with `GatewayInput` → returns `Gateway` (`Input` + `id` +
`created`). A gateway is the Basics Station / UDP link to the LNS that forwards
frames from LoRaWAN devices.

```ts
type GatewayLinkProtocol = 'basicstation' | 'udp';
// GatewayRegion: see the LoRaWAN section above.

// Flat shape covers either protocol; only the relevant fields are used.
interface GatewayLink {
  protocol: GatewayLinkProtocol;
  lnsUri: string;   // used when protocol = 'basicstation' (wss://lns:1887)
  host: string;     // used when protocol = 'udp'
  port: number;     // used when protocol = 'udp' (1700 by convention)
}

interface GatewayInput {
  name: string;
  eui: string;              // 16 hex chars
  enabled: boolean;         // On/Off
  region: GatewayRegion;
  description: string;
  link: GatewayLink;
}

interface Gateway extends GatewayInput {
  id: string;       // server-assigned
  created: string;  // ISO-8601
}
```

### Example — `POST /api/gateways` (Basics Station)

```json
{
  "name": "Rooftop Gateway",
  "eui": "0016C001F1500001",
  "enabled": true,
  "region": "EU868",
  "description": "Basics Station link to the local LNS.",
  "link": {
    "protocol": "basicstation",
    "lnsUri": "wss://127.0.0.1:1887",
    "host": "127.0.0.1",
    "port": 1700
  }
}
```

### Example — `POST /api/gateways` (Semtech UDP)

```json
{
  "name": "Warehouse UDP",
  "eui": "0016C001F1500002",
  "enabled": false,
  "region": "US915",
  "description": "Semtech UDP packet forwarder.",
  "link": {
    "protocol": "udp",
    "lnsUri": "",
    "host": "127.0.0.1",
    "port": 1700
  }
}
```

---

## 4) Connection

A reusable, named target plus its auth for one protocol. A `Run` points at a
connection by id. `config` is the SAME `ProtocolConfig` discriminated union as a
device's `config` (see section 2).

| Method | Path                          | Body              | Returns (in `data`)   |
|--------|-------------------------------|-------------------|-----------------------|
| GET    | `/api/connections`            | —                 | `Connection[]`        |
| POST   | `/api/connections`            | `ConnectionInput` | `Connection`          |
| PUT    | `/api/connections/:id`        | `ConnectionInput` | `Connection`          |
| DELETE | `/api/connections/:id`        | —                 | `null`                |
| POST   | `/api/connections/:id/test`   | —                 | `ConnectionTestResult`|

```ts
interface Connection {
  id: string;             // server-assigned
  created: string;        // ISO-8601 creation timestamp (server-assigned)
  name: string;
  protocolId: ProtocolId;
  config: ProtocolConfig; // discriminated on `kind` (== protocolId); see section 2
}

interface ConnectionInput {
  name: string;
  protocolId: ProtocolId;
  config: ProtocolConfig;
}

// POST /api/connections/:id/test
interface ConnectionTestResult {
  ok: boolean;
  message: string;
  latencyMs?: number;
}
```

**Example** — `POST /api/connections` (MQTT target):

```json
{
  "name": "Local Broker",
  "protocolId": "mqtt",
  "config": {
    "kind": "mqtt",
    "brokerUrl": "mqtt://127.0.0.1:1883",
    "clientId": "sim-01",
    "baseTopic": "sim",
    "authMode": "userpass",
    "username": "sim",
    "password": "secret",
    "tlsCertPem": "",
    "tlsKeyPem": "",
    "tlsCaPem": ""
  }
}
```

---

## 5) Scenario

A traffic profile: a payload template plus rate and volume. `payloadTemplate`
uses the same placeholder tokens listed in section 2.

| Method | Path                         | Body             | Returns (in `data`) |
|--------|------------------------------|------------------|---------------------|
| GET    | `/api/scenarios`             | —                | `Scenario[]`        |
| POST   | `/api/scenarios`             | `ScenarioInput`  | `Scenario`          |
| PUT    | `/api/scenarios/:id`         | `ScenarioInput`  | `Scenario`          |
| DELETE | `/api/scenarios/:id`         | —                | `null`              |
| POST   | `/api/scenarios/:id/preview` | `{ attributes }` | `ScenarioPreview`   |

```ts
interface Scenario {
  id: string;            // server-assigned
  created: string;       // ISO-8601 creation timestamp (server-assigned)
  name: string;
  payloadTemplate: string;
  sendRateHz: number;
  count: number;
  durationSec: number;
}

interface ScenarioInput {
  name: string;
  payloadTemplate: string;
  sendRateHz: number;
  count: number;
  durationSec: number;
}

// POST /api/scenarios/:id/preview
//   body: { attributes: Record<string, string> }   // sample device attributes
interface ScenarioPreview {
  rendered: string;      // the template rendered with the sample attributes
}
```

---

## 6) Run

A simulation run: a connection target plus a set of devices plus a scenario,
tracked through its lifecycle. Live progress streams over the WebSocket
(section 8).

| Method | Path                   | Body            | Returns (in `data`) |
|--------|------------------------|-----------------|---------------------|
| GET    | `/api/runs`            | —               | `Run[]`             |
| GET    | `/api/runs/:id`        | —               | `Run`               |
| POST   | `/api/runs`            | `StartRunInput` | `Run`               |
| POST   | `/api/runs/:id/stop`   | —               | `Run`               |

```ts
type RunState = 'pending' | 'running' | 'stopping' | 'stopped' | 'failed';

interface RunStats {
  sent: number;
  errors: number;
  lastLatencyMs?: number;
}

interface Run {
  id: string;            // server-assigned
  created: string;       // ISO-8601 creation timestamp (server-assigned)
  connectionId: string;
  deviceIds: string[];
  scenarioId: string;
  state: RunState;
  startedAt?: string;    // when the run actually began running
  stoppedAt?: string;    // when it stopped or finished
  stats: RunStats;
}

// POST /api/runs
interface StartRunInput {
  connectionId: string;
  deviceIds: string[];
  scenarioId: string;
}
```

---

## 7) Health

`GET /api/health` → `HealthResponse`. Liveness probe; the Electron launcher polls
it before opening the window. Returned in the standard envelope (`data` carries
the payload) like every other REST endpoint.

```ts
interface HealthResponse {
  status: string;   // "ok"
  version: string;
}
```

**Example** `GET /api/health`:

```json
{ "status": 200, "errors": null, "data": { "status": "ok", "version": "0.1.0" } }
```

---

## 8) WebSocket — live run stream

`GET /ws` (WebSocket upgrade). The engine pushes `RunEvent` frames as a run
progresses: log lines, metric snapshots, and lifecycle status changes share the
one channel. Frames are **raw JSON** — the `{ status, errors, data }` envelope is
for REST responses only. `ts` here is the frame's event time on the live stream,
NOT a persisted entity's `created`.

```ts
type RunEventKind = 'log' | 'metric' | 'status';
type RunLogLevel  = 'info' | 'warn' | 'error';

interface RunEvent {
  kind: RunEventKind;   // 'log' | 'metric' | 'status'
  ts: string;           // ISO-8601 frame time (live stream, not persisted)
  runId: string;
  deviceId?: string;
  level?: RunLogLevel;  // for kind = 'log'
  msg?: string;         // for kind = 'log'
  stats?: RunStats;     // for kind = 'metric'
  state?: RunState;     // for kind = 'status'
}
```

**Example frames:**

```json
{ "kind": "status", "ts": "2026-06-09T03:00:00.000Z", "runId": "run-1", "state": "running" }
{ "kind": "log",    "ts": "2026-06-09T03:00:01.000Z", "runId": "run-1", "deviceId": "d1", "level": "info", "msg": "POST /v1/ingest 200" }
{ "kind": "metric", "ts": "2026-06-09T03:00:02.000Z", "runId": "run-1", "stats": { "sent": 42, "errors": 0, "lastLatencyMs": 12 } }
```
