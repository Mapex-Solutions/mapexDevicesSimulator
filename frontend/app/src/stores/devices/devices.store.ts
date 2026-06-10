/** TYPE IMPORTS */
import type { Device, DeviceInput } from '@services/sim';

/** VUE IMPORTS */
import { ref } from 'vue';

/** SERVICES */
import { defineStore } from 'pinia';
import { sim } from '@services/sim';

/** LOCAL IMPORTS */
import { createDevicesGetters } from './devices.getters';
import { seedDevices } from './devices.seed';

/**
 * Generate a local id for offline-created devices.
 * @returns {string} a unique id
 */
function localId(): string {
	const cryptoApi = globalThis.crypto;
	const suffix = cryptoApi && typeof cryptoApi.randomUUID === 'function' ? cryptoApi.randomUUID() : Math.random().toString(36).slice(2);
	return `dev-local-${suffix}`;
}

/**
 * Simulated devices. While the engine is offline, reads fall back to seed data
 * and writes are applied locally so the UI is fully usable in development.
 */
export const useDevicesStore = defineStore('devices', () => {
	/** STATE */
	const items = ref<Device[]>([]);
	const loading = ref(false);
	const saving = ref(false);

	async function fetch(): Promise<void> {
		loading.value = true;
		try {
			const data = await sim.devices.list();
			items.value = Array.isArray(data) && data.length > 0 ? data : seedDevices();
		} catch {
			if (!items.value.length) items.value = seedDevices();
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
		} catch {
			const created: Device = { id: localId(), created: new Date().toISOString(), ...input };
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
		} catch {
			const existing = items.value.find((item) => item.id === id);
			const updated: Device = { id, created: existing?.created ?? new Date().toISOString(), ...input };
			items.value = items.value.map((item) => (item.id === id ? updated : item));
			return updated;
		} finally {
			saving.value = false;
		}
	}

	async function remove(id: string): Promise<void> {
		try {
			await sim.devices.remove({ id });
		} catch {
			// Engine offline: fall through to a local removal.
		}
		items.value = items.value.filter((item) => item.id !== id);
	}

	const getters = createDevicesGetters({ items });

	return { items, loading, saving, fetch, create, update, remove, ...getters };
});
