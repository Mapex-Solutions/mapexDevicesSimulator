# Quick Test — drive the simulator end to end

> 🇧🇷 Versão em português: [README_pt.md](./README_pt.md)

This folder is a hands-on tour of the **Mapex Devices Simulator**. Each protocol
has its own folder with **copy-paste values** and a step-by-step walk-through so
you can create a device, turn it on, fire an event, and watch the traffic live.

**📖 Guided tutorial with screenshots:**
[STEP-BY-STEP.en.md](./STEP-BY-STEP.en.md) · [STEP-BY-STEP.pt-BR.md](./STEP-BY-STEP.pt-BR.md)

| Protocol | Folder | What it shows |
|----------|--------|---------------|
| HTTP | [`http/`](./http/) | One-shot uplink to any HTTP endpoint |
| MQTT | [`mqtt/`](./mqtt/) | Live publish **and** subscribe (downlink) over a broker |
| LoRaWAN over Semtech UDP | [`lorawan-udp/`](./lorawan-udp/) | A gateway + sensor joining an LNS over the UDP packet-forwarder |
| LoRaWAN over Basics Station | [`lorawan-basic-station/`](./lorawan-basic-station/) | The same sensor carrying its own Basics Station WebSocket link |

---

## 1. Start the app

```bash
# from the repo root
cd frontend
npm install            # first time only
npm run dev:electron   # builds the Go sidecar, then launches the desktop app
```

`dev:electron` compiles the engine into the sidecar and the desktop window spawns
it automatically. To run in a browser instead:

```bash
cd frontend && npm run dev      # opens http://localhost:9100 (SPA)
# in another terminal, start the engine the UI talks to:
cd backend && go run ./service/src --addr 127.0.0.1 --port 5055
```

The header shows a **connection dot**: green = the engine is reachable, grey =
offline. With no engine, lists are simply empty (there is no fake/seed data).

---

## 2. The flow you repeat for every protocol

1. **Devices → New device** — name it, give it a Device ID, pick the protocol.
2. **Connection** — fill the protocol's connection fields (copy-paste from the
   folder's README).
3. **Events** — add an event with the payload to send (copy-paste from the folder).
4. **Save**, then toggle the device **Enabled** in the list.
5. **Console** — open it to watch frames stream live (`up`, `down`, `system`).
6. **Fire event** — from the device row or the console, send on demand and see
   the `up` echo (and any `down` reply) appear instantly.

> LoRaWAN also needs a **Gateway** (Gateways → New gateway) the sensor rides on —
> except Basics Station, where the sensor carries its own link.

---

## 3. Screenshots

The walk-throughs reference images under each folder's `images/`.

---

## 4. Need an LNS for LoRaWAN?

The LoRaWAN folders were verified against **ChirpStack v4** via `chirpstack-docker`.
The quickest local LNS:

```bash
git clone https://github.com/chirpstack/chirpstack-docker
cd chirpstack-docker && docker compose up -d
# UI http://localhost:8088 (admin/admin) · UDP :1700 · Basics Station :3001
```

Each LoRaWAN folder lists exactly what to register in ChirpStack so the keys match.
