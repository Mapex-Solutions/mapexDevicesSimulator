# Schemas do simulador — contrato REST + WebSocket completo

> 🇺🇸 English version: [README.md](./README.md)

Formatos de payload que são fonte da verdade e que o backend Go deve aceitar/retornar,
tirados literalmente do cliente tipado do frontend
(`frontend/src/services/sim/interfaces`). Use-os para derivar os DTOs Go (e qualquer
camada Zod/validação). Os tipos são mostrados em notação TypeScript; exemplos JSON seguem
cada um.

## Envelope de resposta (padrão Mapex)

Toda resposta é envolvida no envelope Mapex padrão. O tipo listado na coluna "Returns
(in `data`)" é o que vem em `data`; no erro, `errors` carrega as mensagens e `data` é
null. Corpos de requisição são o `*Input` cru (não envolvidos).

```json
{ "status": 200, "errors": null, "data": <payload> }
```

Endpoints REST cobertos:

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

Convenções:

- Todos os nomes de campo são camelCase no fio (JSON).
- Toda resposta é o envelope `{ status, errors, data }`; o DELETE também o retorna
  (`200` com `data: null`), não um `204` puro.
- `id` e `created` são **atribuídos pelo servidor** (não fazem parte dos corpos de
  criação `*Input`). `created` é o timestamp de criação em TODA entidade persistida, uma
  string ISO-8601. Não existe variação `ts` / `createdAt` / `timestamp`.
- `config` (device) é uma **união discriminada** por `kind`
  (`http` | `mqtt` | `lorawan` | `basicstation`). O `kind` sempre bate com o
  `protocolId` do dispositivo.
- Segredos (`apiKey`, `password`, blocos PEM TLS, chaves LoRaWAN) chegam nestes payloads
  mas **nunca devem ser logados** pela engine.

---

## Primitivos compartilhados

```ts
type ProtocolId = 'http' | 'mqtt' | 'lorawan' | 'basicstation';

interface KeyValue {
  key: string;
  value: string;
}
```

---

## 1) Log

`GET /api/logs` → `LogPage`. Um Log é uma mensagem de dispositivo persistida (o
histórico em SQLite por trás do stream de console ao vivo).

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

**Exemplo** `GET /api/logs?limit=20&offset=0&protocol=mqtt` (a resposta é o envelope;
`LogPage` é o `data`):

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

## 2) Criar Device

`POST /api/devices` com `DeviceInput` → retorna `Device` (`Input` + `id` + `created`).

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

### config — `ProtocolConfig` (união discriminada por `kind`)

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

Um evento guarda exatamente uma config específica de protocolo que bate com o protocolo
do dispositivo (`http`/`mqtt` usam o corpo compartilhado; `lorawan`/`basicstation` usam
o uplink LoRaWAN). `schedule` é o auto-disparo opcional.

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

### Tokens de placeholder (resolvidos pela engine no envio)

Usados dentro de `body`, `bodyFields[].value`, `topic` MQTT e `payloadHex` LoRaWAN:

```
{{randInt(min,max)}}    inteiro no intervalo, negativos permitidos
{{randFloat(min,max)}}  decimal no intervalo, negativos permitidos
{{now}}                 timestamp ISO-8601 atual
{{counter}}             contador auto-incremental
{{deviceId}}            o campo deviceId do dispositivo
{{uuid}}                UUID aleatório
```

### Exemplo — `POST /api/devices` (device HTTP)

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

### Exemplo — `POST /api/devices` (device LoRaWAN, OTAA)

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

## 3) Criar Gateway

`POST /api/gateways` com `GatewayInput` → retorna `Gateway` (`Input` + `id` + `created`).
Um gateway é o link Basics Station / UDP até o LNS que encaminha os frames dos
dispositivos LoRaWAN.

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

### Exemplo — `POST /api/gateways` (Basics Station)

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

### Exemplo — `POST /api/gateways` (Semtech UDP)

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

Um alvo nomeado e reutilizável mais a sua auth para um protocolo. Um `Run` aponta para
uma connection por id. `config` é a MESMA união discriminada `ProtocolConfig` da `config`
de um device (veja a seção 2).

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

**Exemplo** — `POST /api/connections` (alvo MQTT):

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

Um perfil de tráfego: um template de payload mais taxa e volume. `payloadTemplate` usa os
mesmos tokens de placeholder listados na seção 2.

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

Uma execução de simulação: um alvo de connection mais um conjunto de devices mais um
scenario, acompanhada ao longo do seu ciclo de vida. O progresso ao vivo é transmitido
pelo WebSocket (seção 8).

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

`GET /api/health` → `HealthResponse`. Probe de liveness; o launcher Electron faz polling
nele antes de abrir a janela. Retornado no envelope padrão (`data` carrega o payload)
como todo outro endpoint REST.

```ts
interface HealthResponse {
  status: string;   // "ok"
  version: string;
}
```

**Exemplo** `GET /api/health`:

```json
{ "status": 200, "errors": null, "data": { "status": "ok", "version": "0.1.0" } }
```

---

## 8) WebSocket — stream de run ao vivo

`GET /ws` (upgrade WebSocket). A engine empurra frames `RunEvent` conforme um run avança:
linhas de log, snapshots de métrica e mudanças de status de ciclo de vida compartilham o
mesmo canal. Os frames são **JSON cru** — o envelope `{ status, errors, data }` é só para
respostas REST. `ts` aqui é o tempo do evento do frame no stream ao vivo, NÃO o `created`
de uma entidade persistida.

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

**Frames de exemplo:**

```json
{ "kind": "status", "ts": "2026-06-09T03:00:00.000Z", "runId": "run-1", "state": "running" }
{ "kind": "log",    "ts": "2026-06-09T03:00:01.000Z", "runId": "run-1", "deviceId": "d1", "level": "info", "msg": "POST /v1/ingest 200" }
{ "kind": "metric", "ts": "2026-06-09T03:00:02.000Z", "runId": "run-1", "stats": { "sent": 42, "errors": 0, "lastLatencyMs": 12 } }
```
