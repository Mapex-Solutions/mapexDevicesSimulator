# HTTP quick test

> 🇧🇷 Versão em português: [README_pt.md](./README_pt.md)

An HTTP device is **send-only**: every fire is a one-shot request to your endpoint,
and the response status comes back on the `up` frame. No downlinks.

The default target below is `https://httpbin.org` + path `/post`, which echoes back
a `200` from anywhere with internet. Swap it for your own ingest URL any time.

---

## Create the device (UI)

**Devices → New device**

| Step | Field | Paste this |
|------|-------|------------|
| Info | Name | `Quick HTTP sensor` |
| Info | Device ID | `http-quick-01` |
| Info | Protocol | `HTTP` |
| Connection | URL | `https://httpbin.org` |
| Connection | Method | `POST` |
| Connection | Header | `Content-Type` = `application/json` |
| Connection | Auth | `None` |

![New device](./images/01-new-device.png)
![HTTP connection](./images/02-connection.png)

### Add an event

**Events → Add event**

| Field | Paste this |
|-------|------------|
| Name | `Telemetry` |
| Method | `POST` |
| Path | `/post` |
| Body mode | `Raw` |
| Body | see below |

```json
{
  "deviceId": "{{deviceId}}",
  "temperature": {{randInt(18,30)}},
  "humidity": {{randInt(40,70)}}
}
```

`{{deviceId}}` and `{{randInt(a,b)}}` are rendered by the engine at send time, so
each fire carries fresh values.

![HTTP event](./images/03-event.png)

---

## Run it

1. **Save** the device, then flip **Enabled** on in the list.
2. Open the **Console**.
3. Click **Fire event** on the device row.
4. An `up` frame appears with status `200` and the rendered JSON body.

![Console up frame](./images/04-console.png)

---

## One-command alternative (API)

Creates the same device straight through the engine's REST API:

```bash
bash quickTest/http/curl.sh           # defaults to http://127.0.0.1:5055
# or point at another engine:
SIM=http://127.0.0.1:5080 bash quickTest/http/curl.sh
```
