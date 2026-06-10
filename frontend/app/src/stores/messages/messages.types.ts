/** Direction of a console message. */
export type MessageDirection = 'up' | 'down' | 'system';

/** Nature of a console message. */
export type MessageKind = 'data' | 'auth' | 'join' | 'downlink' | 'status';

/**
 * A single entry in the console message stream: every data uplink/downlink and
 * auth/join handshake the engine emits. The store uses the live frame contract
 * directly so the websocket payload and the in-memory item never drift.
 */
export type { ConsoleMessage } from '@services/sim';
