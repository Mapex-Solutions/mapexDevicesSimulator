#!/usr/bin/env bash
# Create a LoRaWAN OTAA sensor carrying its own Basics Station link, then fire one uplink.
# Usage: SIM=http://127.0.0.1:5055 bash curl.sh
set -euo pipefail
SIM="${SIM:-http://127.0.0.1:5055}"

echo "engine: $SIM"
curl -s "$SIM/api/health" >/dev/null || { echo "engine not reachable at $SIM"; exit 1; }

ID=$(curl -s -X POST "$SIM/api/devices" -H 'Content-Type: application/json' -d '{
  "name": "Quick LoRa Basics",
  "deviceId": "lora-bs-01",
  "protocolId": "basicstation",
  "enabled": true,
  "storeLogs": true,
  "config": {
    "kind": "basicstation",
    "lnsUri": "ws://127.0.0.1:3001",
    "gatewayEui": "0102030405060708",
    "region": "EU868",
    "macVersion": "1.0.3",
    "activation": "otaa",
    "devEui": "0011223344556677",
    "joinEui": "0000000000000000",
    "appKey": "00112233445566778899AABBCCDDEEFF",
    "nwkKey": "",
    "devAddr": "", "nwkSKey": "", "appSKey": ""
  },
  "attributes": {},
  "events": [{
    "id": "evt-lht65n",
    "name": "LHT65N uplink",
    "lorawan": { "fport": 2, "confirmed": false, "payloadHex": "0BB809F6025D0000000000" }
  }]
}' | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
echo "created device: $ID"

sleep 2   # let the OTAA join complete
curl -s -X POST "$SIM/api/devices/$ID/fire" -H 'Content-Type: application/json' \
  -d '{"eventId":"evt-lht65n"}'; echo
echo "watch the Console for the WebSocket connect, join-accept -> joined -> up."
