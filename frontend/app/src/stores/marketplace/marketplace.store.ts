/** TYPE IMPORTS */
import type {
	Device,
	DeviceInput,
	MarketplaceCatalogItem,
	MarketplaceFacets,
} from '@services/sim';
import type { MarketplaceQuery } from './marketplace.types';

/** VUE IMPORTS */
import { ref } from 'vue';

/** SERVICES */
import { defineStore } from 'pinia';
import { sim } from '@services/sim';

/**
 * Connectivity of the online marketplace catalog. `offline` means the catalog
 * server could not be reached — distinct from the sidecar engine status, which
 * governs installed devices.
 */
export type MarketplaceStatus = 'idle' | 'loading' | 'online' | 'offline';

/**
 * The device marketplace: an online, read-only catalog the user browses to add a
 * pre-configured device. Listing and filtering are served by the marketplace
 * service; installing a chosen model goes through the local engine (the devices
 * API), which owns the SQLite write.
 */
export const useMarketplaceStore = defineStore('marketplace', () => {
	/** STATE */
	const items = ref<MarketplaceCatalogItem[]>([]);
	const total = ref(0);
	const facets = ref<MarketplaceFacets | null>(null);
	const status = ref<MarketplaceStatus>('idle');
	const installingId = ref<string | null>(null);

	/**
	 * Fetch a filtered page of catalog cards. The first call also loads the facet
	 * options. A network failure flips the store to `offline` and clears the grid.
	 *
	 * @param {MarketplaceQuery} [query] - the active filters
	 */
	async function fetch(query: MarketplaceQuery = {}): Promise<void> {
		if (typeof navigator !== 'undefined' && navigator.onLine === false) {
			status.value = 'offline';
			items.value = [];
			return;
		}

		status.value = 'loading';
		try {
			const [page, facetSet] = await Promise.all([
				sim.marketplace.list(query),
				facets.value ? Promise.resolve(facets.value) : sim.marketplace.facets(),
			]);
			items.value = page.items;
			total.value = page.total;
			facets.value = facetSet;
			status.value = 'online';
		} catch {
			items.value = [];
			facets.value = null;
			status.value = 'offline';
		}
	}

	/**
	 * Install a catalog model as a new device. Fetches the model's simulator
	 * template, mints a fresh deviceId (the engine assigns the row id and the keys
	 * are the template's defaults), and creates the device through the engine.
	 *
	 * @param {MarketplaceCatalogItem} item - the catalog model to install
	 * @returns {Promise<Device>} the created device
	 */
	async function install(item: MarketplaceCatalogItem): Promise<Device> {
		installingId.value = item.id;
		try {
			const template = await sim.marketplace.simulator({ vendor: item.vendor, slug: item.slug });
			const input: DeviceInput = { ...template, deviceId: mintDeviceId(item) };
			return await sim.devices.create(input);
		} finally {
			installingId.value = null;
		}
	}

	return { items, total, facets, status, installingId, fetch, install };
});

/**
 * Build a readable, unique-enough deviceId for an installed model, e.g.
 * `em300-th-a1b2c3d4`. The engine identifies a device by this id across the
 * console and logs, so it must not collide with an existing one.
 *
 * @param {MarketplaceCatalogItem} item - the model being installed
 * @returns {string} a fresh deviceId
 */
function mintDeviceId(item: MarketplaceCatalogItem): string {
	const slug = item.model.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-|-$/g, '');
	const suffix =
		typeof crypto !== 'undefined' && 'randomUUID' in crypto
			? crypto.randomUUID().slice(0, 8)
			: Math.floor(Math.random() * 0xffffffff).toString(16).padStart(8, '0');
	return `${slug}-${suffix}`;
}
