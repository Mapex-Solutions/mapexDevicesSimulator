/** TYPE IMPORTS */
import type { ConsoleMessage } from '@sim/schema';
import type { ConsoleStreamHandlers, ConsoleStream } from './interfaces';

/** SCHEMAS */
import { ZodConsoleMessageSchema } from '@sim/schema';

/** CLIENT */
import { resolveWsBase } from '../../../client';

/**
 * Opens the live console stream at <ws>/ws and forwards each validated frame to
 * `onMessage`. Reconnects with a fixed backoff while open, since the sidecar may
 * restart; malformed frames are dropped rather than throwing into the socket.
 *
 * @param {ConsoleStreamHandlers} handlers - message/open/close callbacks
 * @param {string} path - the websocket path on the sidecar (defaults to /ws)
 * @returns {ConsoleStream} a handle whose `close` stops the stream for good
 */
export function createConsoleStream(handlers: ConsoleStreamHandlers, path = '/ws'): ConsoleStream {
	let socket: WebSocket | null = null;
	let timer: ReturnType<typeof setTimeout> | null = null;
	let closed = false;

	const url = `${resolveWsBase()}${path}`;

	function connect(): void {
		if (closed) return;
		socket = new WebSocket(url);

		socket.onopen = () => handlers.onOpen?.();

		socket.onmessage = (event: MessageEvent) => {
			let raw: unknown;
			try {
				raw = JSON.parse(typeof event.data === 'string' ? event.data : '');
			} catch {
				return;
			}
			const parsed = ZodConsoleMessageSchema.safeParse(raw);
			if (parsed.success) handlers.onMessage(parsed.data as ConsoleMessage);
		};

		socket.onclose = () => {
			handlers.onClose?.();
			if (!closed) timer = setTimeout(connect, 2000);
		};

		socket.onerror = () => socket?.close();
	}

	connect();

	return {
		close: () => {
			closed = true;
			if (timer) clearTimeout(timer);
			socket?.close();
		},
	};
}
