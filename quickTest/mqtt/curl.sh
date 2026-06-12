#!/usr/bin/env bash
# Create an MQTT device (publish + subscribe) through the engine REST API.
# Usage: SIM=http://127.0.0.1:5055 bash curl.sh
set -euo pipefail
SIM="${SIM:-http://127.0.0.1:5055}"

echo "engine: $SIM"
curl -s "$SIM/api/health" >/dev/null || { echo "engine not reachable at $SIM"; exit 1; }

ID=$(curl -s -X POST "$SIM/api/devices" -H 'Content-Type: application/json' -d '{
  "name": "Quick MQTT sensor",
  "deviceId": "mqtt-quick-01",
  "protocolId": "mqtt",
  "enabled": true,
  "storeLogs": true,
  "config": {
    "kind": "mqtt",
    "brokerUrl": "tcp://broker.hivemq.com:1883",
    "clientId": "mqtt-quick-01",
    "baseTopic": "mapex/quicktest",
    "authMode": "none",
    "username": "", "password": "",
    "tlsCertPem": "", "tlsKeyPem": "", "tlsCaPem": "",
    "receiveEnabled": true,
    "subscriptions": [{ "name": "commands", "topic": "mqtt-quick-01/cmd", "qos": 1 }]
  },
  "attributes": {},
  "events": [{
    "id": "evt-telemetry",
    "name": "Telemetry",
    "mqtt": {
      "topic": "mqtt-quick-01/telemetry",
      "qos": 1,
      "retain": false,
      "bodyMode": "raw",
      "bodyFields": [],
      "body": "{\n  \"deviceId\": \"{{deviceId}}\",\n  \"level\": {{randInt(0,100)}}\n}"
    }
  }]
}' | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

echo "created device: $ID"
curl -s -X POST "$SIM/api/devices/$ID/fire" -H 'Content-Type: application/json' \
  -d '{"eventId":"evt-telemetry"}'; echo
echo "downlink test:"
echo "  mosquitto_pub -h broker.hivemq.com -t mapex/quicktest/mqtt-quick-01/cmd -q 1 -m '{\"cmd\":\"ping\"}'"
