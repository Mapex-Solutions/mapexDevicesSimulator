/** TYPE IMPORTS */
import type { ConsoleMessage } from './messages.types';
import type { ConsoleStream } from '@services/sim';

/** VUE IMPORTS */
import { ref } from 'vue';

/** SERVICES */
import { defineStore } from 'pinia';
import { createConsoleStream } from '@services/sim';

/** LOCAL IMPORTS */
import { createMessagesGetters } from './messages.getters';

// The console keeps only the most recent events as a FIFO ring buffer: once it
// reaches this many, the oldest are dropped as new ones arrive (newest shown on
// top). Older history lives in the persisted Logs page.
const MAX_MESSAGES = 1000;

let seq = 0;
function nextId(): string {
	seq += 1;
	return `m${seq}`;
}

/**
 * The console message stream: every uplink, downlink and auth/join handshake
 * across all protocols. Fed by the engine over the websocket (when connected)
 * and by locally fired events; seeded with examples while the engine is offline
 * so the console conveys the flow.
 */
export const useMessagesStore = defineStore('messages', () => {
	/** STATE */
	const items = ref<ConsoleMessage[]>([]);
	const selectedId = ref<string | null>(null);
	const deviceFilter = ref<string | null>(null);
	const connected = ref(false);

	// The live stream handle lives outside reactive state; it is opened on the
	// console page mount and closed on unmount.
	let stream: ConsoleStream | null = null;

	/**
	 * Append a message to the stream and trim history.
	 * @param {Omit<ConsoleMessage, 'id'>} message - the message without an id
	 */
	function add(message: Omit<ConsoleMessage, 'id'>): ConsoleMessage {
		const entry: ConsoleMessage = { ...message, id: nextId() };
		const next = items.value.concat(entry);
		if (next.length > MAX_MESSAGES) next.splice(0, next.length - MAX_MESSAGES);
		items.value = next;
		return entry;
	}

	function select(id: string | null): void {
		selectedId.value = id;
	}

	function setDeviceFilter(deviceId: string | null): void {
		deviceFilter.value = deviceId;
	}

	function clear(): void {
		items.value = [];
		selectedId.value = null;
	}

	/**
	 * Load example messages spanning protocols, directions and handshakes so the
	 * console is explorable before the engine is wired.
	 */
	function seedSamples(): void {
		if (items.value.length) return;

		const samples: Omit<ConsoleMessage, 'id'>[] = [
			{
				ts: '10:02:11', protocol: 'http', deviceId: 'dev-http-1', deviceName: 'HTTP sensor 01',
				direction: 'up', kind: 'data', status: '200', summary: 'POST /v1/ingest',
				payload: '{\n  "temperature": 21.5,\n  "humidity": 60,\n  "battery": 98\n}',
				meta: { endpoint: '/v1/ingest', auth: 'X-API-Key' },
			},
			{
				ts: '10:02:09', protocol: 'mqtt', deviceId: 'dev-mqtt-1', deviceName: 'MQTT sensor 01',
				direction: 'system', kind: 'auth', status: 'OK', summary: 'CONNECT accepted (user/pass)',
				payload: '{\n  "clientId": "dev-mqtt-1",\n  "username": "device",\n  "result": "authorized"\n}',
				meta: { authType: 'basic', keepAlive: '60s' },
			},
			{
				ts: '10:02:08', protocol: 'mqtt', deviceId: 'dev-mqtt-1', deviceName: 'MQTT sensor 01',
				direction: 'up', kind: 'data', status: 'qos1', summary: 'PUBLISH telemetry/dev-mqtt-1',
				payload: '{\n  "level": 0.72,\n  "flow": 12.4\n}',
				meta: { topic: 'telemetry/dev-mqtt-1', qos: '1' },
			},
			{
				ts: '10:02:07', protocol: 'mqtt', deviceId: 'dev-mqtt-1', deviceName: 'MQTT sensor 01',
				direction: 'down', kind: 'downlink', status: 'sent', summary: 'command/dev-mqtt-1',
				payload: '{\n  "cmd": "set-interval",\n  "seconds": 30\n}',
				meta: { topic: 'command/dev-mqtt-1', qos: '1' },
			},
			{
				ts: '10:02:03', protocol: 'lorawan', deviceId: 'sensor-1', deviceName: 'LoRa sensor 01',
				direction: 'system', kind: 'join', status: 'accepted', summary: 'OTAA join-accept (1.0.3)',
				payload: 'JoinAccept\n  DevAddr: 26 0B AC 12\n  DevEUI:  00 80 E1 15 00 12 34 56\n  MIC ok',
				meta: { version: '1.0.3', gateway: 'gw-1', region: 'EU868' },
			},
			{
				ts: '10:02:01', protocol: 'lorawan', deviceId: 'sensor-1', deviceName: 'LoRa sensor 01',
				direction: 'up', kind: 'data', status: 'FCnt 14', summary: 'Uplink FPort 2',
				payload: 'FRMPayload (decoded)\n{\n  "temp": 19.8,\n  "door": "closed"\n}\n\nraw: 01 8E 00',
				meta: { fPort: '2', fCnt: '14', gateway: 'gw-1', rssi: '-87 dBm', snr: '9.2' },
			},
		];

		for (const sample of samples) add(sample);
	}

	/**
	 * Open the live console stream and append each frame the engine emits. Idempotent:
	 * a second call while connected is a no-op. The engine assigns frame ids, so the
	 * id is dropped and re-issued locally by `add`.
	 */
	function connect(): void {
		if (stream) return;
		stream = createConsoleStream({
			onMessage: ({ id: _id, ...frame }) => add(frame),
			onOpen: () => { connected.value = true; },
			onClose: () => { connected.value = false; },
		});
	}

	/** Close the live console stream. */
	function disconnect(): void {
		stream?.close();
		stream = null;
		connected.value = false;
	}

	const getters = createMessagesGetters({ items, selectedId, deviceFilter });

	return { items, selectedId, deviceFilter, connected, add, select, setDeviceFilter, clear, seedSamples, connect, disconnect, ...getters };
});
