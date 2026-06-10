# LoRaWAN simulation & realtime device sessions

> Status: IMPLEMENTED and verified end-to-end against ChirpStack v4 (Semtech UDP).
> Scope: how the engine gained live device sessions (inbound/downlinks) and a real,
> vendor-neutral LoRaWAN device+gateway that any LNS accepts. Captures the design,
> the code map, and every gotcha learned testing on the wire.

## 1. Why this exists

The first cut of the engine was send-only and stateless: a scheduler fired events
(uplinks) one-shot (`connect → publish → disconnect`) and slept. Two gaps drove this
work:

1. **Devices RECEIVE data.** MQTT devices subscribe to topics; LoRaWAN devices get
   downlinks in the RX window after an uplink. Receiving needs a *persistent
   connection* — and "disconnect after each send" is incompatible with it.
2. **The UI had no realtime command path.** "Fire/send" never reached the engine.

Resolution: the engine gained a second, always-on **session/listener half** beside
the scheduler; commands go over REST; all live results stream over `/ws`.

## 2. The engine's two halves

```
SCHEDULER (time-driven, outbound)        SESSION MANAGER (connection-driven, always-on)
min-heap of (device,event) jobs          one live connection per enabled session-capable device
fires uplinks THROUGH the live session   reconnect forever with backoff (floor 1s / ceiling 30s + jitter)
                                         inbound (downlink) -> console `down` frame + log
                                         every transition -> console `system/status` frame
```

Key rules (locked decisions):
- **`enabled` = live.** An enabled device opens a persistent connection; disabling
  disconnects. Broker/LNS offline ⇒ infinite reconnect, every attempt logged.
- **Commands via REST, events via `/ws`.** `/ws` stays server→UI only.
- **HTTP is uplink-only.** Only MQTT and LoRaWAN have true async inbound.

Code: `engine/application/ports/session_port.go` (Connector/Session/InboundSink),
`engine/application/services/engine_handler_sessions.go` (SessionManager, mirrors the
job reconcile), `engine/infrastructure/session/` (the connectors).

## 3. LoRaWAN fundamentals (the mental model that unlocks everything)

The single most important thing: **a LoRaWAN device and gateway are two layers.**

| Layer | What | Lives on |
|-------|------|----------|
| **Link / transport** (the always-open socket) | Basics Station (WS) or Semtech UDP to the LNS. Trafficks bytes, **decrypts nothing** | **Gateway** |
| **Device session** (crypto + state) | OTAA keys, `devAddr`, frame counters, decrypt, RX-window matching | **Device** |

- The gateway is a **dumb bridge**. One gateway link is shared by many devices.
- The device is the **brain**. It owns the keys and does all crypto.
- **Downlinks are tied to a prior uplink (Class A).** The device's radio is off most
  of the time; it only opens RX1 (~1s) and RX2 (~2s) windows AFTER it transmits. No
  uplink ⇒ the LNS literally cannot reach the device. The join accept itself is a
  downlink in that window.
- **OTAA join:** device sends JoinRequest (devEUI/joinEUI/devNonce, signed with the
  AppKey) → LNS validates the MIC → returns a JoinAccept (downlink) → device derives
  session keys. ABP skips this (keys preset).

Console consequence: link status (connecting/connected/reconnecting) is a **gateway**
concept; join/joined is a **device** concept; a device is "online" only if its gateway
link is up AND it has joined.

## 4. Crypto — vendored from The Things Stack, de-coupled

LoRaWAN crypto is byte-exact: one wrong bit in a MIC/key and any LNS rejects the
frame. We did **not** hand-roll it. Strategy (decided with the user):

- **Copy TTS `pkg/crypto` (6 self-contained files) verbatim** into
  `backend/packages/utils/lorawan/crypto/`, keeping the Apache-2.0 header. Only
  changes: repoint the `types` import to a local minimal `types` package (6 byte-array
  types: AES128Key/DevAddr/EUI64/DevNonce/JoinNonce/NetID), and replace the one
  `pkg/errors` use with a tiny stdlib shim (`errors.go`). Dep added:
  `github.com/jacobsa/crypto/cmac`.
- **Do NOT copy the codec** (`encoding/lorawan`) — it drags `pkg/ttnpb` (14 MB of
  generated protobuf). The PHYPayload byte layout is trivial; we wrote our own thin
  codec over plain structs.
- **Parity is guarded by static golden vectors** (`crypto/golden_test.go`) — zero TTS
  dependency, even in test. Provenance (verbatim copy) guarantees correctness; the
  golden test guards against future drift.

The crypto math is untouched, so its output is identical to TTS — proven later when
ChirpStack accepted our join MIC and decrypted our data uplinks.

## 5. Code map (`backend/packages/utils/lorawan/`)

| Package | Role |
|---------|------|
| `crypto/` | vendored TTS crypto (keys, MIC, encrypt/decrypt), 1.0.x legacy + 1.1 |
| `types/` | minimal LoRaWAN identifier/key types the crypto needs |
| `codec/` | our own PHYPayload marshal/unmarshal (join request, data uplink, downlink decode) |
| `device/` | stateful device session: OTAA join, derive 1.0.x keys, build uplink, decode downlink |
| `band/` | per-region RX1/RX2 + radio-meta tables (EU868, US915, …) |
| `udp/` | Semtech UDP packet-forwarder framing (PUSH_DATA, PULL_DATA, PULL_RESP, TX_ACK, **stat**) |
| `basicstation/` | Basics Station LNS-protocol framing (version, jreq/updf, dnmsg) |

Connector + transports: `engine/infrastructure/session/lorawan_connector.go`
(shared gateway link, ref-counted; OTAA join; devAddr downlink routing) and
`lorawan_link.go` (the real UDP + WebSocket clients).

## 6. Testing against ChirpStack v4 — the process and EVERY gotcha

This is the part worth re-reading before the next live test. We ran the full official
`chirpstack-docker` stack and pointed the simulator's gateway at it over Semtech UDP.

### 6.1 Bring up ChirpStack
- `git clone https://github.com/chirpstack/chirpstack-docker && docker compose up -d`.
- **Gotcha — image pull:** if `docker pull` fails with `personal access token is
  expired`, the stored Docker Hub token is dead. ChirpStack images are PUBLIC: run
  `docker logout` and anonymous pulls work. No login needed.
- **Gotcha — port conflicts:** the compose binds `1883` (mosquitto), `8080` (UI),
  `8090` (REST). If those are taken on the host, edit `docker-compose.yml`: mosquitto
  needs no host port (internal docker network), and remap the UI/REST host ports. We
  used UI `8088`, REST `8092`.
- Services that matter: gateway-bridge **Semtech UDP `:1700/udp`**, gateway-bridge
  **Basics Station `:3001`**, UI `:8088`, REST `:8092`. Default region EU868.

### 6.2 Authenticate to the API
- **Gotcha:** the ChirpStack `chirpstack-rest-api` gateway does **NOT** expose the
  internal login (`/api/internal/login` → 404). Authenticate the way the web UI does:
  the **`InternalService.Login` gRPC-web** call with real `admin/admin` on the
  ChirpStack server port (`:8088`), then use the returned JWT as `Bearer` on the REST
  API (`:8092`). (Do not forge JWTs — auth-bypass; use the real login RPC.)

### 6.3 Provision (gateway + device)
In ChirpStack: a Gateway (EUI matching the simulator's gateway), a Device Profile
(EU868, MAC 1.0.3, reg-params A, OTAA on), an Application, a Device (DevEUI/JoinEUI),
and the OTAA **Application key** (= the simulator's AppKey). In the simulator: a
gateway with a UDP link to `127.0.0.1:1700`, and a `lorawan` OTAA device referencing
it, EUIs/AppKey identical.

### 6.4 The two real bugs found ON THE WIRE (not in unit tests)
1. **TX_ACK must echo the PULL_RESP token.** Our first TX_ACK used a fresh token →
   ChirpStack logged `backend/semtechudp: could not handle packet ... no internal
   frame cache for token N`. Fix: echo the token from bytes 1-2 of the received
   PULL_RESP. (`udpTransport.receive` / `writeTxAck`, asserted in
   `lorawan_udp_integration_test`.)
2. **Gateway shows "never seen" without a periodic `stat`.** ChirpStack marks a
   gateway online from **gateway statistics** (a `stat` PUSH_DATA), not from uplinks.
   A real packet forwarder sends a `stat` ~every 30s. Fix: `udp.MarshalStat` +
   `udpTransport.stats` goroutine (sends on connect + every 30s). After that, the
   gateway's `last_seen_at` updates and it shows online.

### 6.5 Re-joining: the DevNonce gotcha
ChirpStack rejects a re-join with a **reused DevNonce** (replay protection). The
simulator's DevNonce is in-memory and resets to 1 on restart, so a fresh sim start
fails the join with `Handle join-request error ... validate dev-nonce`. To re-test,
clear it in ChirpStack's DB:

```
docker compose exec postgres psql -U chirpstack -d chirpstack \
  -c "update device_keys set dev_nonces='{}'::jsonb where dev_eui=decode('<deveui>','hex');"
```

**Gotcha within the gotcha:** `dev_nonces` is a JSONB **map** `{}` — NOT an array
`[]`. Setting `'[]'` corrupts it (`Error deserializing field 'dev_nonces': invalid
type: sequence, expected a map`) and every join then fails.

### 6.6 Payload codec: v3 vs v4
To decode a real sensor payload, ChirpStack runs the device profile's JS codec.
**Gotcha:** vendor "ChirpStack decoder" files (e.g. Dragino's) are often the **v3
format** — `function Decode(fPort, bytes, variables)`. ChirpStack **v4 requires
`decodeUplink(input)` returning `{ data: {...} }`**. With only `Decode`, ChirpStack
errors `decodeUplink is not defined` (event `code: UPLINK_CODEC`) and `object` is
null. Either wrap it (`function decodeUplink(input){ return { data:
Decode(input.fPort, input.bytes, input.variables) }; }`) or write a native
`decodeUplink`.

To watch decoded uplinks live:
```
docker compose exec mosquitto mosquitto_sub -t 'application/+/device/+/event/up' -C 1
```
The decoded values land in the `object` field.

### 6.7 The proof — a real sensor decoded
We simulated a **Dragino LHT65N**: payload `0BB809F6025D0000000000` on fPort 2
(battery 3.0 V, temp 25.5 °C, humidity 60.5 %). ChirpStack accepted the OTAA join
(valid MIC), decrypted the data uplinks (session keys matched), and the official
Dragino codec decoded the `object` to:

```json
{ "Node_type": "LHT65N", "BatV": 3.0, "TempC_SHT": 25.5, "Hum_SHT": 60.5,
  "Bat_status": 0, "Ext_sensor": "No external sensor" }
```

The whole chain is real and vendor-neutral: simulated end-device → vendored crypto →
Semtech UDP → ChirpStack gateway-bridge → OTAA join → decrypted uplinks → vendor codec
→ real values.

## 7. LoRaWAN 1.1 (also verified)

Both 1.0.x and 1.1 OTAA are implemented and accepted by ChirpStack live. The 1.1
differences, all in `device/device.go` + `codec/codec.go`:
- **Join request MIC** is signed with the **NwkKey** (1.0.x uses the AppKey).
- **Join accept** is encrypted with the **NwkKey** and the device derives three
  network session keys (`DeriveFNwkSIntKey`/`DeriveSNwkSIntKey`/`DeriveNwkSEncKey`
  from NwkKey) plus `AppSKey` from AppKey (1.0.x derives one legacy NwkSKey + AppSKey
  from the single AppKey).
- **Uplink MIC** uses `crypto.ComputeUplinkMIC(sNwkSIntKey, fNwkSIntKey, confFCnt,
  txDR, txCh, addr, fCnt, frame)` — it binds BOTH network keys AND the data-rate and
  channel index the frame was sent on. **Gotcha:** `txDR`/`txCh` must match what the
  LNS derives from the radio metadata (rxpk datr/freq), or the MIC fails. For EU868
  868.1 MHz SF7BW125 that is **DR 5, channel 0** (see `band.Region.UplinkDR/UplinkChannel`).
- FRMPayload is still encrypted with AppSKey.

ChirpStack provisioning for 1.1: device profile MAC version `LORAWAN_1_1_0`, and the
device keys set **both** `nwkKey` (the 1.1 NwkKey) and `appKey` (the 1.1 AppKey). In
the simulator the device carries `macVersion: "1.1.0"` plus both `appKey` and `nwkKey`.
(In ChirpStack's `device_keys` table, a 1.0.x device stores its single key in
`nwk_key`; a 1.1 device stores NwkKey in `nwk_key` and AppKey in `app_key`.)

ABP still uses the legacy MIC path (1.1 ABP would need three preset network keys we
don't model).

## 8. What's done vs pending

Done (verified live against ChirpStack over Semtech UDP): MQTT inbound; LoRaWAN OTAA
**1.0.x AND 1.1** join + data uplinks; gateway-online stats; full crypto parity; the
fire-event REST endpoint; a real Dragino LHT65N payload decoded by a vendor codec; the
frontend pieces (subscriptions step, real fire dialog, derived gateway status).

Pending: **Basics Station (`:3001`) live test** against ChirpStack (UDP fully done);
**dynamic sensor values** (payloads are fixed hex today — varying temp/hum needs a
hex-aware template); and the **backend↔frontend integration pass** (drive the whole
LoRaWAN flow from the Vue/Electron UI instead of curl/REST).
