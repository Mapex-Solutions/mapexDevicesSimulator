/** Lifecycle statuses that mean a healthy, established link or a delivered message. */
const POSITIVE = new Set(['connected', 'subscribed', 'joined', 'activated', 'received', 'online', 'sent']);

/** Statuses that mean the send failed or the link is down. */
const NEGATIVE = new Set(['error', 'timeout', 'disconnected', 'gateway-offline', 'join-failed', 'failed', 'unreachable', 'refused']);

/** Transitional statuses that are neither success nor failure yet. */
const PENDING = new Set(['connecting', 'subscribing', 'reconnecting', 'join-request', 'join-accept', 'pending']);

/**
 * Map a console status to a Quasar color token, so an HTTP code or a link
 * lifecycle status reads at a glance instead of being a flat grey badge.
 *
 * HTTP-range codes drive the color band (2xx green, 3xx blue, 4xx amber,
 * 5xx red); known lifecycle statuses get a semantic color; anything else
 * (e.g. qos1, FCnt 14) stays neutral. Shades are chosen to keep the badge's
 * default white text readable.
 * @param {string | undefined} status - the message status
 * @returns {string} the Quasar color token for the badge
 */
export function statusColor(status?: string): string {
	if (!status) return 'grey-7';

	const code = Number.parseInt(status, 10);
	if (Number.isFinite(code) && code >= 100 && code < 600) {
		if (code < 300) return 'green-7';
		if (code < 400) return 'blue-7';
		if (code < 500) return 'orange-8';
		return 'red-7';
	}

	const key = status.toLowerCase();
	if (POSITIVE.has(key)) return 'green-7';
	if (NEGATIVE.has(key)) return 'red-7';
	if (PENDING.has(key)) return 'orange-8';
	return 'grey-7';
}
