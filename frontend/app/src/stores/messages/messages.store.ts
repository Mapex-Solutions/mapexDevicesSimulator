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
 * across all protocols, fed by the engine over the websocket. Empty until the
 * engine emits frames; connection state is exposed via `connected`.
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

	return { items, selectedId, deviceFilter, connected, add, select, setDeviceFilter, clear, connect, disconnect, ...getters };
});
