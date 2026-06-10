/** TYPE IMPORTS */
import type { ConsoleMessage } from '@sim/schema';

/** Callbacks for the live console WebSocket. */
export interface ConsoleStreamHandlers {
	onMessage: (message: ConsoleMessage) => void;
	onOpen?: () => void;
	onClose?: () => void;
}

/** A handle to the live console stream; call `close` to stop it for good. */
export interface ConsoleStream {
	close: () => void;
}
