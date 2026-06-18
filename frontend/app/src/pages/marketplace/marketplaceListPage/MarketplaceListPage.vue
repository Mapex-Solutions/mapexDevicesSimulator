<template>
	<q-page class="q-pt-lg">
		<PageHeader
			icon="storefront"
			icon-color="primary"
			:title="t('marketplace.title')"
			:description="t('marketplace.subtitle')"
		/>

		<!-- Top filter bar -->
		<q-card flat class="filter-bar q-pa-md q-mb-lg">
			<div class="row q-col-gutter-md items-center">
				<div class="col-12 col-md-4">
					<q-input
						v-model="filters.search"
						dense
						outlined
						clearable
						debounce="250"
						:placeholder="t('marketplace.searchPlaceholder')"
					>
						<template #prepend><q-icon name="search" /></template>
					</q-input>
				</div>
				<div class="col-12 col-sm-4 col-md">
					<q-select
						v-model="filters.protocol"
						dense
						outlined
						clearable
						emit-value
						map-options
						:options="protocolOptions"
						:label="t('marketplace.protocol')"
					/>
				</div>
				<div class="col-12 col-sm-4 col-md">
					<q-select
						v-model="filters.readingType"
						dense
						outlined
						clearable
						emit-value
						map-options
						:options="readingOptions"
						:label="t('marketplace.readingType')"
					/>
				</div>
				<div class="col-12 col-sm-4 col-md">
					<q-select
						v-model="filters.manufacturer"
						dense
						outlined
						clearable
						emit-value
						map-options
						:options="manufacturerOptions"
						:label="t('marketplace.manufacturer')"
					/>
				</div>
				<div class="col-auto">
					<q-btn
						flat
						dense
						no-caps
						color="grey-7"
						icon="filter_alt_off"
						:label="t('marketplace.reset')"
						:disable="!hasActiveFilters"
						@click="resetFilters"
					/>
				</div>
			</div>
		</q-card>

		<!-- Offline state -->
		<ListCardEmpty
			v-if="marketplaceStore.status === 'offline'"
			icon="wifi_off"
			:title="t('marketplace.offlineTitle')"
			:description="t('marketplace.offline')"
			:button-label="t('marketplace.retry')"
			button-icon="refresh"
			@button-click="reload"
		/>

		<!-- Loading -->
		<div v-else-if="marketplaceStore.status === 'loading'" class="row justify-center q-my-xl">
			<q-spinner color="primary" size="3em" />
		</div>

		<!-- Catalog grid -->
		<template v-else>
			<div class="row q-col-gutter-md">
				<div v-for="card in cards" :key="card.id" class="col-12 col-sm-6 col-md-4">
					<MarketplaceCard
						:item="card"
						:installing="marketplaceStore.installingId === card.id"
						:add-label="t('marketplace.add')"
						:details-label="t('marketplace.details')"
						:manual-label="t('marketplace.manual')"
						:codec-label="t('marketplace.codec')"
						@open="openDetail"
						@install="onInstall"
					/>
				</div>
			</div>

			<ListCardEmpty
				v-if="!cards.length"
				icon="storefront"
				:title="t('marketplace.emptyTitle')"
				:description="t('marketplace.empty')"
				:button-label="hasActiveFilters ? t('marketplace.reset') : t('marketplace.retry')"
				:button-icon="hasActiveFilters ? 'filter_alt_off' : 'refresh'"
				@button-click="onEmptyAction"
			/>
		</template>

		<!-- Details modal -->
		<MarketplaceDetailDialog
			v-model="detailOpen"
			:item="selectedItem"
			:info="detailInfo"
			:codecs="detailCodecs"
			:loading="detailLoading"
			:installing="!!selectedItem && marketplaceStore.installingId === selectedItem.id"
			@install="onInstall"
		/>

		<!-- Install / customize modal -->
		<MarketplaceInstallDialog
			:model-value="installOpen"
			:item="installItem"
			:draft="installDraft"
			:loading="!installDraft"
			:installing="installing"
			@update:model-value="onInstallOpenChange"
			@confirm="onConfirmInstall"
		/>
	</q-page>
</template>

<script setup lang="ts">
/** TYPE IMPORTS */
import type { MarketplaceInformation, MarketplaceCodec, DeviceInput } from '@services/sim';
import type { MarketplaceQuery } from '@stores/marketplace';
import type {
	MarketplaceCardItem,
	MarketplaceFilters,
	MarketplaceReadingMeta,
} from './interfaces/marketplaceListPage.interface';

/** VUE IMPORTS */
import { computed, onMounted, onBeforeUnmount, ref, watch } from 'vue';

/** COMPONENTS */
import { PageHeader } from '@components/PageHeader';
import { ListCardEmpty } from '@components/ListCardEmpty';
import MarketplaceCard from './components/MarketplaceCard.vue';
import MarketplaceDetailDialog from './components/MarketplaceDetailDialog.vue';
import MarketplaceInstallDialog from './components/MarketplaceInstallDialog.vue';

/** COMPOSABLES */
import { useTranslations } from '@composables/i18n';

/** UTILS */
import { notifySuccess, notifyFail } from '@utils/alert';

/** SERVICES & STORES */
import { useRouter } from 'vue-router';
import { sim } from '@services/sim';
import { useMarketplaceStore } from '@stores/marketplace';

/** COMPOSABLES & STORES */
const { t, te, locale } = useTranslations();

/**
 * Localized label for a reading-type facet: the app i18n label keyed by value,
 * falling back to the catalog's server-provided label when no key exists.
 */
function readingLabel(value: string, fallback: string): string {
	const key = `marketplace.facets.readingType.${value}`;
	return te(key) ? t(key) : fallback;
}
const router = useRouter();
const marketplaceStore = useMarketplaceStore();

/** STATE */
const filters = ref<MarketplaceFilters>({ search: '', protocol: null, readingType: null, manufacturer: null });
const detailOpen = ref(false);
const selectedItem = ref<MarketplaceCardItem | null>(null);
const detailInfo = ref<MarketplaceInformation | null>(null);
const detailCodecs = ref<MarketplaceCodec[]>([]);
const detailLoading = ref(false);
const installOpen = ref(false);
const installItem = ref<MarketplaceCardItem | null>(null);
const installDraft = ref<DeviceInput | null>(null);
const installing = ref(false);
/** Whether the install dialog was opened from the detail dialog, so cancelling
 * the install returns to it instead of leaving the user on the bare grid. */
const reopenDetailAfterInstall = ref(false);

/** COMPUTED */

/** Map of reading-type value → display meta, from the server facets. */
const readingMetaByValue = computed((): Record<string, MarketplaceReadingMeta> => {
	const map: Record<string, MarketplaceReadingMeta> = {};
	for (const reading of marketplaceStore.facets?.readingTypes ?? []) {
		map[reading.value] = { value: reading.value, label: readingLabel(reading.value, reading.label), icon: reading.icon || 'mdi-gauge' };
	}
	return map;
});

/** Map of protocol value → display label, from the server facets. */
const protocolLabelByValue = computed((): Record<string, string> => {
	const map: Record<string, string> = {};
	for (const protocol of marketplaceStore.facets?.protocols ?? []) map[protocol.value] = protocol.label;
	return map;
});

const protocolOptions = computed(() =>
	(marketplaceStore.facets?.protocols ?? []).map((p) => ({ label: p.label, value: p.value })),
);

const readingOptions = computed(() =>
	(marketplaceStore.facets?.readingTypes ?? []).map((r) => ({ label: readingLabel(r.value, r.label), value: r.value })),
);

const manufacturerOptions = computed(() =>
	(marketplaceStore.facets?.manufacturers ?? []).map((m) => ({ label: m.label, value: m.value })),
);

const hasActiveFilters = computed(
	() => !!filters.value.search || !!filters.value.protocol || !!filters.value.readingType || !!filters.value.manufacturer,
);

/** The active filters mapped to the catalog query (empty values dropped). */
const query = computed((): MarketplaceQuery => {
	const built: MarketplaceQuery = { lang: locale.value };
	if (filters.value.protocol) built.protocol = filters.value.protocol;
	if (filters.value.readingType) built.readingType = filters.value.readingType;
	if (filters.value.manufacturer) built.manufacturer = filters.value.manufacturer;
	if (filters.value.search) built.search = filters.value.search;
	return built;
});

/** The catalog cards enriched with display labels for the protocol and readings. */
const cards = computed((): MarketplaceCardItem[] =>
	marketplaceStore.items.map((item) => ({
		...item,
		protocolLabel: protocolLabelByValue.value[item.protocol] ?? item.protocol.toUpperCase(),
		readingMetas: item.readingTypes.map(
			(value) => readingMetaByValue.value[value] ?? { value, label: value, icon: 'mdi-gauge' },
		),
	})),
);

/** WATCHERS */

/** Re-query the catalog whenever a filter changes. */
watch(query, (next) => void marketplaceStore.fetch(next));

/** FUNCTIONS */

function resetFilters(): void {
	filters.value = { search: '', protocol: null, readingType: null, manufacturer: null };
}

/** Reload the current query (offline retry / refresh). */
function reload(): void {
	void marketplaceStore.fetch(query.value);
}

/** Empty-grid action: clear filters when any are active, otherwise reload. */
function onEmptyAction(): void {
	if (hasActiveFilters.value) resetFilters();
	else reload();
}

/**
 * Open the details modal for a card and fetch its information sheet and codecs.
 * @param {MarketplaceCardItem} item - the selected catalog item
 */
async function openDetail(item: MarketplaceCardItem): Promise<void> {
	selectedItem.value = item;
	detailInfo.value = null;
	detailCodecs.value = [];
	detailOpen.value = true;
	detailLoading.value = true;
	try {
		const [info, codecs] = await Promise.all([
			sim.marketplace.information({ vendor: item.vendor, slug: item.slug }),
			sim.marketplace.codecs({ vendor: item.vendor, slug: item.slug }),
		]);
		detailInfo.value = info;
		detailCodecs.value = codecs;
	} catch {
		detailInfo.value = null;
		detailCodecs.value = [];
	} finally {
		detailLoading.value = false;
	}
}

/**
 * Open the install dialog for a catalog model, pre-filled with its template so the
 * user can adjust identity/credentials (and pick a gateway for LoRaWAN) before the
 * device is created.
 * @param {MarketplaceCardItem} item - the catalog item to install
 */
async function onInstall(item: MarketplaceCardItem): Promise<void> {
	// Never stack on the detail dialog: remember it was open, close it, then open
	// the install dialog. Cancelling the install reopens the detail (see below).
	reopenDetailAfterInstall.value = detailOpen.value;
	detailOpen.value = false;
	installItem.value = item;
	installDraft.value = null;
	installOpen.value = true;
	try {
		installDraft.value = await marketplaceStore.prepareInstall(item);
	} catch {
		installOpen.value = false;
		reopenDetailAfterInstall.value = false;
		notifyFail({ message: t('marketplace.addFailed') });
	}
}

/**
 * React to the install dialog opening/closing. When it closes WITHOUT a confirmed
 * install (cancel / backdrop / Esc) and it was opened from the detail dialog,
 * reopen the detail so the user lands back where they were.
 * @param {boolean} open - the dialog's next open state
 */
function onInstallOpenChange(open: boolean): void {
	installOpen.value = open;
	if (!open && reopenDetailAfterInstall.value) {
		reopenDetailAfterInstall.value = false;
		detailOpen.value = true;
	}
}

/**
 * Create the device from the (user-edited) draft, then offer to open the devices
 * list.
 * @param {DeviceInput} input - the finalized device draft
 */
async function onConfirmInstall(input: DeviceInput): Promise<void> {
	installing.value = true;
	try {
		const device = await marketplaceStore.confirmInstall(input);
		// A successful install closes both modals — clear the flag first so the
		// install dialog's close handler does not reopen the detail.
		reopenDetailAfterInstall.value = false;
		installOpen.value = false;
		detailOpen.value = false;
		notifySuccess({
			message: t('marketplace.added', { name: device.name }),
			actions: [{ label: t('marketplace.view'), color: 'white', handler: () => void router.push({ name: 'devices' }) }],
		});
	} catch {
		notifyFail({ message: t('marketplace.addFailed') });
	} finally {
		installing.value = false;
	}
}

/** Re-probe the catalog when connectivity changes. */
function onConnectivityChange(): void {
	reload();
}

/** LIFECYCLE HOOKS */
onMounted(() => {
	void marketplaceStore.fetch(query.value);
	window.addEventListener('online', onConnectivityChange);
	window.addEventListener('offline', onConnectivityChange);
});

onBeforeUnmount(() => {
	window.removeEventListener('online', onConnectivityChange);
	window.removeEventListener('offline', onConnectivityChange);
});
</script>

<style scoped lang="scss">
.filter-bar {
	background: var(--mapex-surface-elevated);
	border: 1px solid var(--mapex-card-border);
	border-radius: var(--mapex-radius-lg);
}
</style>
