/** TYPE IMPORTS */
import type { Gateway, GatewayInput } from '@services/sim';

/** VUE IMPORTS */
import { ref } from 'vue';

/** SERVICES */
import { defineStore } from 'pinia';
import { sim } from '@services/sim';

/** LOCAL IMPORTS */
import { createGatewaysGetters } from './gateways.getters';

/**
 * Gateways backed by the engine over REST. Reads reflect the backend exactly
 * (empty when it has none); writes go straight through and let failures surface
 * to the caller. Engine-offline is shown globally via the sidecar health status.
 */
export const useGatewaysStore = defineStore('gateways', () => {
	/** STATE */
	const items = ref<Gateway[]>([]);
	const loading = ref(false);
	const saving = ref(false);

	async function fetch(): Promise<void> {
		loading.value = true;
		try {
			items.value = await sim.gateways.list();
		} catch {
			items.value = [];
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
		} finally {
			saving.value = false;
		}
	}

	async function remove(id: string): Promise<void> {
		await sim.gateways.remove({ id });
		items.value = items.value.filter((item) => item.id !== id);
	}

	const getters = createGatewaysGetters({ items });

	return { items, loading, saving, fetch, create, update, remove, ...getters };
});
