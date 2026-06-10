/** VUE IMPORTS */
import { computed } from 'vue';

/** STORES */
import { useDevicesStore } from '@stores/devices';
import { useMessagesStore } from '@stores/messages';

/** A gateway's live link state, derived from its devices' console status frames. */
export type GatewayConnection = 'online' | 'connecting' | 'offline' | 'unknown';

const ONLINE = new Set(['connected', 'joined', 'subscribed', 'activated']);
const CONNECTING = new Set(['connecting', 'subscribing', 'reconnecting', 'join-request', 'join-accept']);

/**
 * Derive each gateway's live connection state from the console stream. LoRaWAN
 * devices report their link lifecycle as system/status frames; a gateway reflects
 * the latest frame from any device that transmits through it.
 *
 * @returns {{ connectionOf: (gatewayId: string) => GatewayConnection }} the lookup
 */
export function useGatewayConnections(): { connectionOf: (gatewayId: string) => GatewayConnection } {
	const devicesStore = useDevicesStore();
	const messagesStore = useMessagesStore();

	const statusByGateway = computed<Map<string, GatewayConnection>>(() => {
		const deviceToGateway = new Map<string, string>();
		for (const device of devicesStore.items) {
			if (device.protocolId === 'lorawan' && device.config.kind === 'lorawan') {
				deviceToGateway.set(device.deviceId, device.config.gatewayId);
			}
		}

		const latest = new Map<string, string>();
		for (const message of messagesStore.items) {
			if (message.direction !== 'system') continue;
			const gatewayId = deviceToGateway.get(message.deviceId);
			if (!gatewayId) continue;
			latest.set(gatewayId, message.status ?? message.summary);
		}

		const out = new Map<string, GatewayConnection>();
		for (const [gatewayId, status] of latest) {
			out.set(gatewayId, ONLINE.has(status) ? 'online' : CONNECTING.has(status) ? 'connecting' : 'offline');
		}
		return out;
	});

	/**
	 * Resolve one gateway's live connection state.
	 * @param {string} gatewayId - the gateway id
	 * @returns {GatewayConnection} the derived state
	 */
	function connectionOf(gatewayId: string): GatewayConnection {
		return statusByGateway.value.get(gatewayId) ?? 'unknown';
	}

	return { connectionOf };
}
