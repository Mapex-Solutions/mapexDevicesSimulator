/** TYPE IMPORTS */
import type { Device } from '@services/sim';

/**
 * Sample devices used as a fallback while the simulator engine is offline, so
 * the UI is usable in development. Once the engine serves real data this seed is
 * ignored. One device per enabled protocol.
 * @returns {Device[]} the seed devices
 */
export function seedDevices(): Device[] {
	return [
		{
			id: 'dev-seed-http',
			name: 'Edge Sensor 01',
			deviceId: '9f2c1a4e-http-0001',
			created: '2026-05-12T09:24:00.000Z',
			protocolId: 'http',
			enabled: true,
			storeLogs: true,
			config: {
				kind: 'http',
				url: 'http://localhost:5001/v1/ingest',
				method: 'POST',
				headers: [{ key: 'Content-Type', value: 'application/json' }],
				authMode: 'apiKey',
				apiKeyHeader: 'X-API-Key',
				apiKey: 'demo-key',
				basicUser: '',
				basicPass: '',
			},
			attributes: {},
			events: [
				{
					id: 'evt-1',
					name: 'Telemetry',
					http: {
						method: 'POST',
						path: '/v1/ingest',
						headers: [],
						bodyMode: 'raw',
						bodyFields: [],
						body: '{\n  "deviceId": "{{deviceId}}",\n  "temp": {{randInt(10,30)}}\n}',
					},
					schedule: { enabled: true, every: 30, unit: 'seconds' },
				},
			],
		},
		{
			id: 'dev-seed-mqtt',
			name: 'Greenhouse Hub',
			deviceId: '9f2c1a4e-mqtt-0002',
			created: '2026-05-28T14:10:00.000Z',
			protocolId: 'mqtt',
			enabled: true,
			storeLogs: true,
			config: {
				kind: 'mqtt',
				brokerUrl: 'mqtt://127.0.0.1:1883',
				clientId: 'greenhouse-hub',
				baseTopic: 'greenhouse',
				authMode: 'userpass',
				username: 'demo',
				password: 'demo',
				tlsCertPem: '',
				tlsKeyPem: '',
				tlsCaPem: '',
				receiveEnabled: false,
				subscriptions: [],
			},
			attributes: {},
			events: [
				{
					id: 'evt-1',
					name: 'Humidity',
					mqtt: {
						topic: 'greenhouse/{{deviceId}}/humidity',
						qos: 1,
						retain: false,
						bodyMode: 'raw',
						bodyFields: [],
						body: '{\n  "humidity": {{randInt(40,80)}}\n}',
					},
					schedule: { enabled: false, every: 60, unit: 'seconds' },
				},
			],
		},
		{
			id: 'dev-seed-lora',
			name: 'Field Node LoRa',
			deviceId: '9f2c1a4e-lora-0003',
			created: '2026-06-03T07:45:00.000Z',
			protocolId: 'lorawan',
			enabled: false,
			storeLogs: true,
			config: {
				kind: 'lorawan',
				gatewayId: 'gw-seed-1',
				region: 'EU868',
				macVersion: '1.0.4',
				activation: 'otaa',
				devEui: '70B3D57ED0000001',
				joinEui: '0000000000000000',
				appKey: '00112233445566778899AABBCCDDEEFF',
				nwkKey: '',
				devAddr: '',
				nwkSKey: '',
				appSKey: '',
			},
			attributes: {},
			events: [
				{
					id: 'evt-1',
					name: 'Uplink',
					lorawan: { fport: 10, confirmed: false, payloadHex: '01{{randInt(0,255)}}' },
					schedule: { enabled: true, every: 5, unit: 'minutes' },
				},
			],
		},
	];
}
