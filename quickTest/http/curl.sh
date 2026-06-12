#!/usr/bin/env bash
# Create an HTTP device through the engine REST API, then fire one event.
# Usage: SIM=http://127.0.0.1:5055 bash curl.sh
set -euo pipefail
SIM="${SIM:-http://127.0.0.1:5055}"

echo "engine: $SIM"
curl -s "$SIM/api/health" >/dev/null || { echo "engine not reachable at $SIM"; exit 1; }

ID=$(curl -s -X POST "$SIM/api/devices" -H 'Content-Type: application/json' -d '{
  "name": "Quick HTTP sensor",
  "deviceId": "http-quick-01",
  "protocolId": "http",
  "enabled": true,
  "storeLogs": true,
  "config": {
    "kind": "http",
    "url": "https://httpbin.org",
    "method": "POST",
    "headers": [{ "key": "Content-Type", "value": "application/json" }],
    "authMode": "none",
    "apiKeyHeader": "",
    "apiKey": "",
    "basicUser": "",
    "basicPass": ""
  },
  "attributes": {},
  "events": [{
    "id": "evt-telemetry",
    "name": "Telemetry",
    "http": {
      "method": "POST",
      "path": "/post",
      "headers": [],
      "bodyMode": "raw",
      "bodyFields": [],
      "body": "{\n  \"deviceId\": \"{{deviceId}}\",\n  \"temperature\": {{randInt(18,30)}},\n  \"humidity\": {{randInt(40,70)}}\n}"
    }
  }]
}' | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

echo "created device: $ID"
curl -s -X POST "$SIM/api/devices/$ID/fire" -H 'Content-Type: application/json' \
  -d '{"eventId":"evt-telemetry"}'; echo
echo "open the Console (or GET $SIM/api/logs) to see the up frame."
