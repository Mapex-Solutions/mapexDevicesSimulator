# Mapex Devices Simulator — step by step

A guided, screenshot-driven tour. By the end you will have sent real traffic over
**HTTP**, **MQTT**, and **LoRaWAN** (both Semtech UDP and Basics Station) and seen
it arrive — including a real LoRaWAN device decrypted by **ChirpStack**.

> 🇧🇷 Versão em português: [STEP-BY-STEP.pt-BR.md](./STEP-BY-STEP.pt-BR.md)

---

## 0. Start the app

```bash
cd frontend
npm install            # first run only
npm run dev:electron   # builds the Go sidecar and opens the desktop app
```

The header shows a connection dot: green = the engine is reachable. With no engine
the lists are simply empty — there is no fake/seed data.

The repeating flow for every protocol is the same:

1. **Devices → New device** — name, Device ID, protocol.
2. **Target** — the connection (endpoint, broker, or LoRaWAN keys).
3. **Events** — the payload to send.
4. **Save**, then toggle the device **Enabled**.
5. **Console** — watch frames stream live; **Fire event** to send on demand.

![New device](./http/images/01-new-device.png)

---

## 1. HTTP

HTTP is send-only: every fire is one request and the response status comes back on
the `up` frame.

**Target** — URL `https://httpbin.org`, Method `POST`, Auth `None`:

![HTTP target](./http/images/02-connection.png)

**Event** — Method `POST`, Path `/post`, Body (raw):

```json
{ "deviceId": "{{deviceId}}", "temperature": {{randInt(18,30)}}, "humidity": {{randInt(40,70)}} }
```

![HTTP event](./http/images/03-event.png)

Enable the device, open the **Console**, and **Fire event** — an `up` frame appears
with status `200`:

![HTTP console](./http/images/04-console.png)

Full reference + one-command API path: [`http/`](./http/).

---

## 2. MQTT

MQTT keeps a live broker connection: it publishes uplinks and, with **Receive** on,
streams every message on its subscribed topics as a `down` frame.

**Target** — Broker `tcp://broker.hivemq.com:1883`, Base topic `mapex/quicktest`,
Receive **on**, Subscription `mqtt-quick-01/cmd` (QoS 1):

> Topics are relative to the base topic — the engine prepends `mapex/quicktest`.

![MQTT target + subscription](./mqtt/images/01-connection.png)

**Event** — Topic `mqtt-quick-01/telemetry`, QoS 1, Body `{ "level": {{randInt(0,100)}} }`:

![MQTT event](./mqtt/images/02-event.png)

Enable + open the **Console**: you'll see `connecting → connected → subscribed`, and
firing publishes an `up` frame (`qos1`). Publish to the subscribed topic from any
client to see a live `down` frame:

```bash
mosquitto_pub -h broker.hivemq.com -t 'mapex/quicktest/mqtt-quick-01/cmd' -q 1 -m '{"cmd":"ping"}'
```

![MQTT console](./mqtt/images/03-console.png)

Full reference: [`mqtt/`](./mqtt/).

---

## 3. LoRaWAN

A LoRaWAN sensor rides on a **gateway** that carries the link to a LoRaWAN Network
Server (LNS). The sensor does an OTAA join, then sends a real **Dragino LHT65N**
uplink the LNS decrypts. This section uses **ChirpStack** as the LNS.

### 3.1 Bring up ChirpStack

```bash
git clone https://github.com/chirpstack/chirpstack-docker
cd chirpstack-docker && docker compose up -d
# UI http://localhost:8088 (admin/admin) · UDP :1700 · Basics Station :3001
```

![ChirpStack login](./chirpstack/images/01-login.png)

### 3.2 Provision ChirpStack (so the keys match)

In the ChirpStack UI, under the **ChirpStack** tenant:

1. **Device Profiles → Add** — region `EU868`, MAC version `LoRaWAN 1.0.3`, **OTAA**
   enabled. (Optionally paste the Dragino LHT65N codec to decode the payload.)
2. **Gateways → Add gateway** — Gateway ID `0102030405060708`, region `EU868`.
3. **Applications → (your app) → Add device** — DevEUI `0011223344556677`,
   Join EUI `0000000000000000`, select the profile above.
4. On the device, **OTAA keys** — **Application key** `00112233445566778899AABBCCDDEEFF`.

> ChirpStack stores a 1.0.x device's single key in its `nwk_key` slot — the UI calls
> it the **Application key**. If you re-run a join, clear the used nonces first:
> ```bash
> docker compose exec postgres psql -U chirpstack -d chirpstack \
>   -c "update device_keys set dev_nonces='{}'::jsonb where dev_eui=decode('0011223344556677','hex');"
> ```

The gateways list — the simulator's gateway shows **Online** once it forwards:

![ChirpStack gateways](./chirpstack/images/03-gateways.png)

Opening it shows the live link metrics (received/transmitted, frequency, DR):

![ChirpStack gateway detail](./chirpstack/images/04-gateway-detail.png)

### 3.3 LoRaWAN over Semtech UDP

In the simulator, create the **gateway** (Gateways → New gateway), step **LNS link**:
Protocol `Semtech UDP`, Host `127.0.0.1`, Port `1700`:

![Gateway LNS link](./lorawan-udp/images/01-gateway.png)

Then the **device** (Protocol `LoRaWAN`), **Target** step — Gateway = the one above,
Region `EU868`, MAC `1.0.3`, Activation `OTAA`, DevEUI `0011223344556677`,
JoinEUI `0000000000000000`, AppKey `00112233445566778899AABBCCDDEEFF`:

![LoRaWAN target](./lorawan-udp/images/02-connection.png)

**Event** — a real LHT65N uplink: FPort `2`, Payload (hex) `0BB809F6025D0000000000`:

![LoRaWAN event](./lorawan-udp/images/03-event.png)

Enable + open the **Console**: `join-request → join-accept → joined`, then `Uplink
FCnt …`. The simulator console shows the full multi-protocol flow live — HTTP `200`,
MQTT `qos1`, LoRaWAN UDP uplink + downlink, and Basics Station uplink + downlink:

![LoRaWAN console](./lorawan-udp/images/04-console.png)

In ChirpStack, the device dashboard shows the uplinks arriving (Received / RSSI /
SNR / frequency / DR):

![ChirpStack device](./chirpstack/images/07-device-dashboard.png)

The **Events** tab lists the live join and uplinks as they arrive:

![ChirpStack device events](./chirpstack/images/08-device-events.png)

Full reference + one-command API path: [`lorawan-udp/`](./lorawan-udp/).

### 3.4 LoRaWAN over Basics Station

Same sensor, but the device carries its **own Basics Station WebSocket link** — no
separate gateway. Protocol `Basics Station`, **Target** step — LNS URI
`ws://127.0.0.1:3001`, Gateway EUI `0102030405060708`, same region/keys as above:

![Basics Station target](./lorawan-basic-station/images/01-connection.png)

The same LHT65N event; enable + console shows the WebSocket connect, join, and uplink:

![Basics Station console](./lorawan-basic-station/images/03-console.png)

Full reference: [`lorawan-basic-station/`](./lorawan-basic-station/).

---

## One-command paths

Each folder has a `curl.sh` that creates its device through the engine REST API and
fires once:

```bash
bash quickTest/http/curl.sh
bash quickTest/mqtt/curl.sh
bash quickTest/lorawan-udp/curl.sh
bash quickTest/lorawan-basic-station/curl.sh
```

## Regenerating the screenshots

```bash
cd frontend && npm i -D playwright        # API only; uses your installed Chrome
npm run dev                                # SPA on :9100
node INTERNALS/capture-screenshots.mjs     # simulator screenshots
node INTERNALS/capture-chirpstack.mjs      # ChirpStack screenshots
```
