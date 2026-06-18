# Device marketplace — mock catalog (POC)

> 🇧🇷 Versão em português: [README_pt.md](./README_pt.md)

This folder is a **proof-of-concept mock** of the online `mapexMarketplace` catalog,
served statically from the app's `public/` in development. It lets the marketplace UI
([feature 14](../../../../RELEASES/1.0.0/features/14-marketplace.md)) be browsed and
exercised without the real catalog service. In production the app is pointed at the
deployed `mapexMarketplace` service (`/api/v1/devices`) instead — the bundle format
below mirrors that service's responses.

The catalog is **read-only**. Installing a model writes a real device into the local
engine through the normal device-create path; nothing here is ever mutated.

## Layout

```
marketplace/
├── registry.json                 # the full catalog index (every card)
├── taxonomy.json                 # filter vocabularies (protocols, reading types)
└── vendors/
    └── <vendor>/<model>/
        ├── device_information.json   # the detail sheet
        ├── device_simulator.json     # the installable template
        └── codec/                    # optional payload codec(s)
            └── decoder.js
```

Current mock: **5 models across 3 vendors** (Milesight, Dragino, Generic), 3 of them
LoRaWAN with a reference codec.

## Files

### `registry.json` — the catalog index

One entry per device card; this is what the listing grid and filters read. Each
`item` carries identity (`id`, `vendor`, `vendorName`, `model`, `slug`, `name`,
`description`), classification (`protocol`, `readingTypes`, `tags`, `icon`), the
bundle `path` (`vendors/<vendor>/<model>`), and the `hasCodec` / `hasManual` flags
the UI uses to decide which detail actions to show.

```json
{
  "id": "milesight-em300-th",
  "vendor": "milesight",
  "vendorName": "Milesight",
  "model": "EM300-TH",
  "slug": "em300-th",
  "name": "EM300-TH Temperature & Humidity Sensor",
  "description": "Indoor LoRaWAN sensor reporting temperature and relative humidity…",
  "protocol": "lorawan",
  "readingTypes": ["temperature", "humidity"],
  "tags": ["indoor", "battery", "climate"],
  "icon": "mdi-thermometer",
  "path": "vendors/milesight/em300-th",
  "hasCodec": true,
  "hasManual": true
}
```

### `taxonomy.json` — filter vocabularies

The `protocol` and `readingType` filter options (value + display label, plus an icon
for reading types). It is the source for the marketplace facets; keep new values here
in sync with what the entries actually use. Manufacturer facets are derived from the
registry, so they are not listed here.

### `vendors/<vendor>/<model>/device_information.json` — the detail sheet

What the detail dialog renders: full description, `vendor` block (site/support),
`tags`, `version`, `images`, and the document links. Documents may be a **bundled
PDF** (`datasheet` / `manual` file) and/or an **external link** (`datasheetUrl` /
`manualUrl`) — the UI shows "view PDF" when a file exists, else "open online". An
optional `codec` block points at the model's reference decoder.

### `vendors/<vendor>/<model>/device_simulator.json` — the installable template

The payload of `install`: a device `DeviceInput` **minus the `deviceId`**, which the
engine mints fresh on install. It is exactly the shape the device-create API accepts
— `name`, `protocolId`, `enabled`, `storeLogs`, `config` (protocol-specific),
`attributes`, and a sample `events` list (templated, e.g.
`{{randInt(180,260)}}`, `{{deviceId}}`, `{{now}}`). After install, every value is a
plain device field the user can edit.

### `vendors/<vendor>/<model>/codec/` — payload codec (optional)

The vendor's reference decoder (and/or encoder). It is surfaced for the user to
**view and download**; the simulator does not execute it yet.

## Adding a model

1. Create `vendors/<vendor>/<model>/` with `device_information.json` and
   `device_simulator.json` (and `codec/` if it ships one).
2. Add a matching card to `registry.json` (`path` = the new folder, set
   `hasCodec` / `hasManual`).
3. If it introduces a new protocol or reading type, add it to `taxonomy.json` so the
   filters can offer it.
