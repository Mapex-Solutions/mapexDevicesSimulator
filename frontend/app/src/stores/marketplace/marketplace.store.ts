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
	 * Fetch a page of catalog cards. With `append`, the page is concatenated onto
	 * the current grid (infinite scroll); otherwise it replaces it. The first call
	 * also loads the facet options. A network failure on the first page flips the
	 * store to `offline` and clears the grid; a failure while appending keeps what
	 * is already loaded so a transient load-more error never blanks the list.
	 *
	 * @param {MarketplaceQuery} [query] - the active filters plus page/perPage
	 * @param {boolean} [append] - concatenate onto the grid instead of replacing it
	 */
	async function fetch(query: MarketplaceQuery = {}, append = false): Promise<void> {
		if (typeof navigator !== 'undefined' && navigator.onLine === false) {
			if (!append) {
				status.value = 'offline';
				items.value = [];
				total.value = 0;
			}
			return;
		}

		status.value = 'loading';
		try {
			const [page, facetSet] = await Promise.all([
				sim.marketplace.list(query),
				facets.value ? Promise.resolve(facets.value) : sim.marketplace.facets(),
			]);
			items.value = append ? [...items.value, ...page.items] : page.items;
			total.value = page.total;
			facets.value = facetSet;
			status.value = 'online';
		} catch {
			if (append) {
				// Keep the already-loaded cards; the next scroll can retry.
				status.value = 'online';
			} else {
				items.value = [];
				total.value = 0;
				facets.value = null;
				status.value = 'offline';
			}
		}
	}

	/** Drop the current grid so the next fetch reloads from the first page. */
	function reset(): void {
		items.value = [];
		total.value = 0;
	}

	/**
	 * Fetch a catalog model's simulator template and build an editable device
	 * draft: the template's defaults plus a freshly minted deviceId. The user
	 * adjusts identity/credentials in the install dialog before it is created, so
	 * this does NOT touch the engine yet.
	 *
	 * @param {MarketplaceCatalogItem} item - the catalog model being installed
	 * @returns {Promise<DeviceInput>} an editable device draft
	 */
	async function prepareInstall(item: MarketplaceCatalogItem): Promise<DeviceInput> {
		installingId.value = item.id;
		try {
			const template = await sim.marketplace.simulator({ vendor: item.vendor, slug: item.slug });
			return { ...template, deviceId: mintDeviceId(item) };
		} finally {
			installingId.value = null;
		}
	}

	/**
	 * Create the device from the (possibly user-edited) draft through the engine,
	 * which owns the SQLite write.
	 *
	 * @param {DeviceInput} input - the finalized device draft
	 * @returns {Promise<Device>} the created device
	 */
	async function confirmInstall(input: DeviceInput): Promise<Device> {
		return await sim.devices.create(input);
	}

	return { items, total, facets, status, installingId, fetch, reset, prepareInstall, confirmInstall };
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
