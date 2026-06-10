/** TYPE IMPORTS */
import type { Log, LogQuery } from '@services/sim';

/** VUE IMPORTS */
import { computed, ref } from 'vue';

/** SERVICES */
import { defineStore } from 'pinia';
import { sim } from '@services/sim';

let seq = 0;
function nextId(): string {
	seq += 1;
	return `log-${seq}`;
}

/** Sample message templates used to populate the table while offline. */
const SAMPLE_TEMPLATES: Omit<Log, 'id' | 'created'>[] = [
	{
		protocol: 'http', deviceId: 'dev-http-1', deviceName: 'HTTP sensor 01', direction: 'up', kind: 'data',
		status: '200', summary: 'POST /v1/ingest', payload: '{\n  "temperature": 21.5,\n  "humidity": 60\n}',
	},
	{
		protocol: 'mqtt', deviceId: 'dev-mqtt-1', deviceName: 'MQTT sensor 01', direction: 'system', kind: 'auth',
		status: 'OK', summary: 'CONNECT accepted (user/pass)', payload: '{\n  "username": "device",\n  "result": "authorized"\n}',
	},
	{
		protocol: 'mqtt', deviceId: 'dev-mqtt-1', deviceName: 'MQTT sensor 01', direction: 'up', kind: 'data',
		status: 'qos1', summary: 'PUBLISH telemetry/dev-mqtt-1', payload: '{\n  "level": 0.72\n}',
	},
	{
		protocol: 'mqtt', deviceId: 'dev-mqtt-1', deviceName: 'MQTT sensor 01', direction: 'down', kind: 'downlink',
		status: 'sent', summary: 'command/dev-mqtt-1', payload: '{\n  "cmd": "set-interval",\n  "seconds": 30\n}',
	},
	{
		protocol: 'lorawan', deviceId: 'sensor-1', deviceName: 'LoRa sensor 01', direction: 'system', kind: 'join',
		status: 'accepted', summary: 'OTAA join-accept (1.0.3)', payload: 'JoinAccept\n  DevAddr: 26 0B AC 12\n  MIC ok',
	},
	{
		protocol: 'lorawan', deviceId: 'sensor-1', deviceName: 'LoRa sensor 01', direction: 'up', kind: 'data',
		status: 'FCnt 14', summary: 'Uplink FPort 2', payload: 'FRMPayload (decoded)\n{\n  "temp": 19.8\n}',
	},
];

function buildSamples(): Log[] {
	const out: Log[] = [];
	const now = Date.now();
	for (let i = 0; i < 42; i += 1) {
		const tpl = SAMPLE_TEMPLATES[i % SAMPLE_TEMPLATES.length];
		if (!tpl) continue;
		// Newest first: each sample is one minute older than the previous.
		out.push({ ...tpl, id: nextId(), created: new Date(now - i * 60_000).toISOString() });
	}
	return out;
}

/**
 * Persisted logs: the SQLite-backed history of device messages. Reads call the
 * engine's paginated/filterable endpoint; while the engine is offline it falls
 * back to seeded samples paginated and filtered client-side so the table is
 * explorable.
 */
export const useLogsStore = defineStore('logs', () => {
	/** STATE */
	const all = ref<Log[]>([]);
	const items = ref<Log[]>([]);
	const total = ref(0);
	const loading = ref(false);
	const page = ref(1);
	const itemsPerPage = ref(15);
	const filters = ref<Record<string, string>>({});
	const lastUpdatedAt = ref<number | undefined>(undefined);

	/** COMPUTED */
	const totalPages = computed(() => Math.max(1, Math.ceil(total.value / itemsPerPage.value)));

	function applyLocal(): void {
		if (!all.value.length) all.value = buildSamples();

		const f = filters.value;
		let list = all.value;
		if (f.protocol) list = list.filter((log) => log.protocol === f.protocol);
		if (f.kind) list = list.filter((log) => log.kind === f.kind);
		if (f.direction) list = list.filter((log) => log.direction === f.direction);
		if (f.device) list = list.filter((log) => log.deviceId === f.device);
		if (f.q) {
			const q = f.q.toLowerCase();
			list = list.filter((log) => `${log.summary} ${log.payload} ${log.deviceName}`.toLowerCase().includes(q));
		}

		total.value = list.length;
		const start = (page.value - 1) * itemsPerPage.value;
		items.value = list.slice(start, start + itemsPerPage.value);
	}

	async function fetch(): Promise<void> {
		loading.value = true;
		try {
			const query: LogQuery = {
				limit: itemsPerPage.value,
				offset: (page.value - 1) * itemsPerPage.value,
				...filters.value,
			};
			const res = await sim.logs.list(query);
			items.value = res.items;
			total.value = res.total;
		} catch {
			applyLocal();
		} finally {
			loading.value = false;
			lastUpdatedAt.value = Date.now();
		}
	}

	function setPage(value: number): void {
		page.value = value;
		void fetch();
	}

	function setItemsPerPage(value: number): void {
		itemsPerPage.value = value;
		page.value = 1;
		void fetch();
	}

	function setFilters(value: Record<string, string>): void {
		filters.value = value;
		page.value = 1;
		void fetch();
	}

	return {
		items, total, loading, page, itemsPerPage, filters, totalPages, lastUpdatedAt,
		fetch, setPage, setItemsPerPage, setFilters,
	};
});
