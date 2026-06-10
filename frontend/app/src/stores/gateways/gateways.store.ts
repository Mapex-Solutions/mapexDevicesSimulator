/** TYPE IMPORTS */
import type { Gateway, GatewayInput } from '@services/sim';

/** VUE IMPORTS */
import { ref } from 'vue';

/** SERVICES */
import { defineStore } from 'pinia';
import { sim } from '@services/sim';

/** LOCAL IMPORTS */
import { createGatewaysGetters } from './gateways.getters';
import { seedGateways } from './gateways.seed';

/**
 * Generate a local id for offline-created gateways.
 * @returns {string} a unique id
 */
function localId(): string {
	const cryptoApi = globalThis.crypto;
	const suffix = cryptoApi && typeof cryptoApi.randomUUID === 'function' ? cryptoApi.randomUUID() : Math.random().toString(36).slice(2);
	return `gw-local-${suffix}`;
}

/**
 * Simulated gateways. While the engine is offline, reads fall back to seed data
 * and writes are applied locally so the UI is fully usable in development.
 */
export const useGatewaysStore = defineStore('gateways', () => {
	/** STATE */
	const items = ref<Gateway[]>([]);
	const loading = ref(false);
	const saving = ref(false);

	async function fetch(): Promise<void> {
		loading.value = true;
		try {
			const data = await sim.gateways.list();
			items.value = Array.isArray(data) && data.length > 0 ? data : seedGateways();
		} catch {
			if (!items.value.length) items.value = seedGateways();
		} finally {
			loading.value = false;
		}
	}

	async function create(input: GatewayInput): Promise<Gateway> {
		saving.value = true;
		try {
			const created = await sim.gateways.create(input);
			items.value = [...items.value, created];
			return created;
		} catch {
			const created: Gateway = { id: localId(), created: new Date().toISOString(), ...input };
			items.value = [...items.value, created];
			return created;
		} finally {
			saving.value = false;
		}
	}

	async function update(id: string, input: GatewayInput): Promise<Gateway> {
		saving.value = true;
		try {
			const updated = await sim.gateways.update({ id }, input);
			items.value = items.value.map((item) => (item.id === id ? updated : item));
			return updated;
		} catch {
			const existing = items.value.find((item) => item.id === id);
			const updated: Gateway = { id, created: existing?.created ?? new Date().toISOString(), ...input };
			items.value = items.value.map((item) => (item.id === id ? updated : item));
			return updated;
		} finally {
			saving.value = false;
		}
	}

	async function remove(id: string): Promise<void> {
		try {
			await sim.gateways.remove({ id });
		} catch {
			// Engine offline: fall through to a local removal.
		}
		items.value = items.value.filter((item) => item.id !== id);
	}

	const getters = createGatewaysGetters({ items });

	return { items, loading, saving, fetch, create, update, remove, ...getters };
});
