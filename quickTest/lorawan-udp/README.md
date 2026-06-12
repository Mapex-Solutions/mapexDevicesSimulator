# LoRaWAN quick test — Semtech UDP

A LoRaWAN sensor rides on a **gateway** that carries the link to your LNS. Here the
gateway speaks the classic **Semtech UDP packet-forwarder** to ChirpStack's
gateway-bridge on `:1700`. The sensor does an OTAA join, then sends a real
**Dragino LHT65N** uplink the LNS decodes to temperature/humidity/battery.

> Needs a running LNS. See the repo root [`../README.md`](../README.md) §4 for a
> one-line ChirpStack via `chirpstack-docker`.

---

## Register in ChirpStack (so the keys match)

- **Gateway** with EUI `0102030405060708`.
- **Device profile**: region `EU868`, MAC version `LoRaWAN 1.0.3`, OTAA enabled.
  (Optional: paste the Dragino LHT65N JS codec to decode the payload.)
- **Device**: DevEUI `0011223344556677`, JoinEUI `0000000000000000`.
- **Application key (OTAA)**: `00112233445566778899AABBCCDDEEFF`.

> Re-joining is rejected if the DevNonce repeats. To re-run from scratch, clear it:
> ```bash
> docker compose exec postgres psql -U chirpstack -d chirpstack \
>   -c "update device_keys set dev_nonces='{}'::jsonb where dev_eui=decode('0011223344556677','hex');"
> ```

---

## Create the gateway (UI)

**Gateways → New gateway**

| Field | Paste this |
|-------|------------|
| Name | `Quick UDP gateway` |
| EUI | `0102030405060708` |
| Region | `EU868` |
| Link protocol | `Semtech UDP` |
| Host | `127.0.0.1` |
| Port | `1700` |

![New gateway](./images/01-gateway.png)

Enable it — the console shows the gateway coming **online** (it sends periodic
stats so the LNS marks it seen).

## Create the device (UI)

**Devices → New device**

| Step | Field | Paste this |
|------|-------|------------|
| Info | Name | `Quick LoRa UDP` |
| Info | Device ID | `lora-udp-01` |
| Info | Protocol | `LoRaWAN` |
| Connection | Gateway | `Quick UDP gateway` |
| Connection | Region | `EU868` |
| Connection | MAC version | `1.0.3` |
| Connection | Activation | `OTAA` |
| Connection | DevEUI | `0011223344556677` |
| Connection | JoinEUI | `0000000000000000` |
| Connection | AppKey | `00112233445566778899AABBCCDDEEFF` |

![LoRaWAN connection](./images/02-connection.png)

### Add an event — a real LHT65N uplink

| Field | Paste this |
|-------|------------|
| Name | `LHT65N uplink` |
| FPort | `2` |
| Confirmed | off |
| Payload (hex) | `0BB809F6025D0000000000` |

That payload decodes to `BatV 3.0, TempC_SHT 25.5, Hum_SHT 60.5`.

![LoRaWAN event](./images/03-event.png)

---

## Run it

1. **Save**, flip **Enabled** on.
2. Open the **Console** — watch `join-request → join-accept → joined`, then on
   **Fire event** an `up` frame with `FCnt`.
3. In ChirpStack, the device shows the uplink (and the decoded object if you added
   the codec).

![LoRaWAN console](./images/04-console.png)

---

## One-command alternative (API)

Creates the gateway, then the device referencing it, then fires once:

```bash
bash quickTest/lorawan-udp/curl.sh    # defaults to http://127.0.0.1:5055
```
