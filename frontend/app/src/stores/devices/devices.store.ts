/** TYPE IMPORTS */
import type { Device, DeviceInput } from '@services/sim';

/** VUE IMPORTS */
import { ref } from 'vue';

/** SERVICES */
import { defineStore } from 'pinia';
import { sim } from '@services/sim';

/** LOCAL IMPORTS */
import { createDevicesGetters } from './devices.getters';

/**
 * Devices backed by the engine over REST. Reads reflect the backend exactly
 * (empty when it has none); writes go straight through and let failures surface
 * to the caller. Engine-offline is shown globally via the sidecar health status.
 */
export const useDevicesStore = defineStore('devices', () => {
	/** STATE */
	const items = ref<Device[]>([]);
	const loading = ref(false);
	const saving = ref(false);

	async function fetch(): Promise<void> {
		loading.value = true;
		try {
			items.value = await sim.devices.list();
		} catch {
			items.value = [];
		} finally {
			loading.value = false;
		}
	}

	async function create(input: DeviceInput): Promise<Device> {
		saving.value = true;
		try {
			const created = await sim.devices.create(input);
			items.value = [...items.value, created];
			return created;
		} finally {
			saving.value = false;
		}
	}

	async function update(id: string, input: DeviceInput): Promise<Device> {
		saving.value = true;
		try {
			const updated = await sim.devices.update({ id }, input);
			items.value = items.value.map((item) => (item.id === id ? updated : item));
			return updated;
		} finally {
			saving.value = false;
		}
	}

	async function remove(id: string): Promise<void> {
		await sim.devices.remove({ id });
		items.value = items.value.filter((item) => item.id !== id);
	}

	const getters = createDevicesGetters({ items });

	return { items, loading, saving, fetch, create, update, remove, ...getters };
});
