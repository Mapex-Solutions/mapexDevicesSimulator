<template>
	<q-page class="q-pt-lg">
		<!-- Header -->
		<PageHeader
			icon="receipt_long"
			icon-color="primary"
			:title="t('logs.title')"
			:description="t('logs.subtitle')"
		/>

		<!-- Filters -->
		<div class="text-caption text-grey-7 q-mb-xs">{{ t('logs.filters') }}</div>
		<div class="row items-center q-col-gutter-sm q-mb-md">
			<div class="col">
				<q-input
					v-model="quickSearch"
					outlined
					dense
					clearable
					hide-bottom-space
					:placeholder="t('logs.searchPlaceholder')"
					class="filter-input"
					@keyup.enter="applyQuickFilters"
					@clear="applyQuickFilters"
				>
					<template #prepend><q-icon name="search" color="grey-6" /></template>
				</q-input>
			</div>

			<div class="col-auto" style="min-width: 160px">
				<q-select
					v-model="quickProtocol"
					outlined
					dense
					emit-value
					map-options
					:options="protocolOptions"
					:label="t('console.filter.protocol')"
					class="filter-input"
					@update:model-value="applyQuickFilters"
					hide-bottom-space
				/>
			</div>

			<div class="col-auto">
				<q-btn round flat icon="tune" color="grey-7" @click="showFiltersDrawer = true">
					<q-badge
						v-if="advancedFiltersCount > 0 || hasPendingAdvancedFilters"
						:color="hasPendingAdvancedFilters ? 'warning' : 'primary'"
						floating
						rounded
						:label="advancedFiltersCount || '!'"
					/>
					<AppTooltip :content="hasPendingAdvancedFilters ? t('logs.pendingFilters') : t('logs.advancedFilters')" />
				</q-btn>
			</div>
		</div>

		<!-- Active filter chips -->
		<div v-if="hasActiveFilters" class="row items-center q-mb-md q-gutter-xs">
			<q-chip
				v-for="chip in visibleFilterChips"
				:key="chip.key"
				removable
				dense
				outline
				color="primary"
				size="sm"
				@remove="removeFilter(chip.key)"
			>
				<span class="text-weight-medium">{{ chip.label }}:</span>&nbsp;{{ chip.value }}
			</q-chip>

			<q-badge v-if="hiddenFiltersCount > 0" color="primary" outline class="q-pa-xs cursor-pointer">
				+{{ hiddenFiltersCount }}
				<AppTooltip>
					<div v-for="chip in hiddenFilterChips" :key="chip.key" class="q-mb-xs">
						<span class="text-weight-medium">{{ chip.label }}:</span> {{ chip.value }}
					</div>
				</AppTooltip>
			</q-badge>

			<q-btn
				flat
				dense
				size="sm"
				color="grey-7"
				icon="filter_alt_off"
				:label="t('logs.clearAll')"
				class="q-ml-sm"
				@click="clearAllFilters"
			/>
		</div>

		<!-- Results header -->
		<div class="row items-center q-pt-md q-mb-md">
			<div class="col">
				<div class="row items-center">
					<q-icon name="receipt_long" size="sm" color="primary" class="q-mr-sm" />
					<div class="text-subtitle1 text-weight-medium text-primary">{{ t('logs.listTitle') }}</div>
				</div>
			</div>
			<div class="col-auto">
				<ListHeaderMenu
					icon="receipt_long"
					:item-label="t('logs.itemLabel')"
					:item-label-plural="t('logs.itemLabelPlural')"
					:items-count="logsStore.total"
					:items-per-page="logsStore.itemsPerPage"
					:columns="menuColumns"
					:filtered="hasActiveFilters"
					:refreshing="logsStore.loading"
					:last-updated-at="logsStore.lastUpdatedAt"
					@update:items-per-page="handleItemsPerPageChange"
					@update:columns="handleColumnsUpdate"
					@refresh="logsStore.fetch"
				/>
			</div>
		</div>

		<!-- Loading -->
		<div v-if="logsStore.loading" class="row justify-center q-my-lg">
			<q-spinner color="primary" size="3em" />
		</div>

		<!-- Rows -->
		<div v-else class="row">
			<div v-for="log in logsStore.items" :key="log.id" class="col-12 q-mb-xs">
				<DataRow :data="log" :columns="visibleColumns" :show-actions="false" />
			</div>

			<div v-if="!logsStore.items.length" class="col-12">
				<ListCardEmpty
					icon="receipt_long"
					:title="t('logs.emptyTitle')"
					:description="t('logs.emptyDescription')"
					@button-click="clearAllFilters"
				/>
			</div>
		</div>

		<!-- Pagination -->
		<div class="row justify-center q-mt-lg q-mb-lg">
			<q-pagination
				v-if="logsStore.items.length"
				:model-value="logsStore.page"
				direction-links
				boundary-links
				class="rounded-borders"
				color="primary"
				active-color="primary"
				:max="logsStore.totalPages"
				@update:model-value="handlePageChange"
			/>
		</div>

		<!-- Advanced filters -->
		<AdvancedFiltersDrawer
			v-model="showFiltersDrawer"
			:title="t('logs.filters')"
			:fields="advancedFilterFields"
			:values="advancedFilterValues"
			@apply="handleAdvancedFiltersApply"
			@reset="handleAdvancedFiltersReset"
			@pending-change="handlePendingChange"
		/>
	</q-page>
</template>

<script setup lang="ts">
/** TYPE IMPORTS */
import type { ListHeaderMenuColumn } from '@components/ListHeaderMenu';
import type { DataRowColumn } from '@components/DataRow';
import type { FilterField, FilterValues } from '@components/AdvancedFiltersDrawer';
import type { LogDirection, ProtocolId } from '@services/sim';

/** VUE IMPORTS */
import { computed, onMounted, ref } from 'vue';

/** COMPONENTS */
import { PageHeader } from '@components/PageHeader';
import { ListHeaderMenu } from '@components/ListHeaderMenu';
import { DataRow } from '@components/DataRow';
import { ListCardEmpty } from '@components/ListCardEmpty';
import { AdvancedFiltersDrawer } from '@components/AdvancedFiltersDrawer';
import { AppTooltip } from '@components/AppTooltip';

/** COMPOSABLES */
import { useTranslations } from '@composables/i18n';

/** STORES */
import { useLogsStore } from '@stores/logs';

const MAX_VISIBLE_CHIPS = 2;
const PROTOCOL_ICON: Record<string, string> = {
	http: 'mdi-web',
	mqtt: 'mdi-transit-connection-variant',
	lorawan: 'mdi-access-point',
	basicstation: 'mdi-radio-tower',
};

/** COMPOSABLES & STORES */
const { t } = useTranslations();
const logsStore = useLogsStore();

/** STATE */
const quickSearch = ref('');
const quickProtocol = ref<string | null>(null);
const showFiltersDrawer = ref(false);
const advancedFilterValues = ref<FilterValues>({ kind: null, direction: null });
const hasPendingAdvancedFilters = ref(false);
const filters = ref<{ search?: string; protocol?: string; kind?: string; direction?: string }>({});
const columnVisibility = ref({ direction: true, kind: true, status: true });

/** COMPUTED */
const protocolOptions = computed(() => [
	{ label: t('console.allDevices'), value: null },
	...(['http', 'mqtt', 'lorawan', 'basicstation'] as ProtocolId[]).map((p) => ({ label: t(`protocol.${p}`), value: p })),
]);

const advancedFilterFields = computed<FilterField[]>(() => [
	{
		key: 'kind', type: 'select', label: t('console.filter.kind'), icon: 'category',
		options: [
			{ label: t('console.allDevices'), value: null },
			...(['data', 'auth', 'join', 'downlink', 'status'] as const).map((k) => ({ label: t(`console.kind.${k}`), value: k })),
		],
	},
	{
		key: 'direction', type: 'select', label: t('console.filter.direction'), icon: 'swap_vert',
		options: [
			{ label: t('console.allDevices'), value: null },
			{ label: t('console.dirUp'), value: 'up' },
			{ label: t('console.dirDown'), value: 'down' },
			{ label: t('console.dirSystem'), value: 'system' },
		],
	},
]);

const hasActiveFilters = computed(() =>
	Boolean(filters.value.search || filters.value.protocol || filters.value.kind || filters.value.direction),
);

const advancedFiltersCount = computed(() => {
	let count = 0;
	if (filters.value.kind) count += 1;
	if (filters.value.direction) count += 1;
	return count;
});

const activeFilterChips = computed(() => {
	const chips: { key: string; label: string; value: string }[] = [];
	if (filters.value.search) chips.push({ key: 'search', label: t('console.filter.search'), value: filters.value.search });
	if (filters.value.protocol) chips.push({ key: 'protocol', label: t('console.filter.protocol'), value: t(`protocol.${filters.value.protocol}`) });
	if (filters.value.kind) chips.push({ key: 'kind', label: t('console.filter.kind'), value: t(`console.kind.${filters.value.kind}`) });
	if (filters.value.direction) chips.push({ key: 'direction', label: t('console.filter.direction'), value: dirLabel(filters.value.direction as LogDirection) });
	return chips;
});

const visibleFilterChips = computed(() => activeFilterChips.value.slice(0, MAX_VISIBLE_CHIPS));
const hiddenFilterChips = computed(() => activeFilterChips.value.slice(MAX_VISIBLE_CHIPS));
const hiddenFiltersCount = computed(() => hiddenFilterChips.value.length);

const menuColumns = computed<ListHeaderMenuColumn[]>(() => [
	{ key: 'direction', label: t('logs.col.direction'), visible: columnVisibility.value.direction },
	{ key: 'kind', label: t('logs.col.kind'), visible: columnVisibility.value.kind },
	{ key: 'status', label: t('logs.col.status'), visible: columnVisibility.value.status },
]);

const logColumns = computed<DataRowColumn[]>(() => [
	{ key: 'icon', label: '', type: 'avatar', visible: 'always', width: 72, icon: (_v, row) => PROTOCOL_ICON[row.protocol] ?? 'mdi-chip', color: () => 'primary', tooltip: (_v, row) => t(`protocol.${row.protocol}`) },
	{ key: 'deviceName', label: t('logs.col.device'), type: 'text', visible: 'always', width: 200, ellipsis: true, secondaryKey: 'deviceId' },
	{ key: 'protocol', label: t('logs.col.protocol'), type: 'chip', visible: 'laptop', width: 130, format: (v) => t(`protocol.${v}`) },
	{ key: 'direction', label: t('logs.col.direction'), type: 'chip', visible: 'laptop', width: 120, format: (v) => dirLabel(v as LogDirection), color: (v) => dirColor(v as LogDirection) },
	{ key: 'kind', label: t('logs.col.kind'), type: 'badge', visible: 'laptop', width: 110, format: (v) => t(`console.kind.${v}`) },
	{ key: 'summary', label: t('logs.col.summary'), type: 'text', visible: 'always', ellipsis: true },
	{ key: 'status', label: t('logs.col.status'), type: 'text', visible: 'desktop', width: 110 },
	{ key: 'created', label: t('logs.col.time'), type: 'text', visible: 'laptop', width: 190, format: (v) => formatDateTime(v as string) },
]);

const visibleColumns = computed(() =>
	logColumns.value.filter((col) => {
		if (col.key === 'direction') return columnVisibility.value.direction;
		if (col.key === 'kind') return columnVisibility.value.kind;
		if (col.key === 'status') return columnVisibility.value.status;
		return true;
	}),
);

/** FUNCTIONS */

/**
 * Format an ISO timestamp as a short local date and time.
 * @param {string} iso - the ISO timestamp
 * @returns {string} the formatted date and time
 */
function formatDateTime(iso: string): string {
	if (!iso) return '—';
	const date = new Date(iso);
	if (Number.isNaN(date.getTime())) return iso;
	return `${date.toLocaleDateString()} ${date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}`;
}

/**
 * Human label for a message direction.
 * @param {LogDirection} direction - the direction
 */
function dirLabel(direction: LogDirection): string {
	if (direction === 'up') return t('console.dirUp');
	if (direction === 'down') return t('console.dirDown');
	return t('console.dirSystem');
}

/**
 * Chip color for a message direction.
 * @param {LogDirection} direction - the direction
 */
function dirColor(direction: LogDirection): string {
	if (direction === 'up') return 'teal';
	if (direction === 'down') return 'primary';
	return 'grey';
}

/**
 * Push the current filters to the store (which resets the page and reloads).
 */
function syncStore(): void {
	const map: Record<string, string> = {};
	if (filters.value.search) map.q = filters.value.search;
	if (filters.value.protocol) map.protocol = filters.value.protocol;
	if (filters.value.kind) map.kind = filters.value.kind;
	if (filters.value.direction) map.direction = filters.value.direction;
	logsStore.setFilters(map);
}

/**
 * Apply the quick filters (search + protocol).
 */
function applyQuickFilters(): void {
	if (quickSearch.value) filters.value.search = quickSearch.value;
	else delete filters.value.search;

	if (quickProtocol.value) filters.value.protocol = quickProtocol.value;
	else delete filters.value.protocol;

	syncStore();
}

/**
 * Apply advanced filters from the drawer.
 * @param {FilterValues} values - the drawer values
 */
function handleAdvancedFiltersApply(values: FilterValues): void {
	advancedFilterValues.value = values;
	if (values.kind) filters.value.kind = String(values.kind);
	else delete filters.value.kind;
	if (values.direction) filters.value.direction = String(values.direction);
	else delete filters.value.direction;
	hasPendingAdvancedFilters.value = false;
	syncStore();
}

/**
 * Track pending (unapplied) changes from the drawer.
 * @param {boolean} pending - whether there are pending changes
 */
function handlePendingChange(pending: boolean): void {
	hasPendingAdvancedFilters.value = pending;
}

function handleAdvancedFiltersReset(): void {
	advancedFilterValues.value = { kind: null, direction: null };
}

/**
 * Remove a single active filter.
 * @param {string} key - the filter key
 */
function removeFilter(key: string): void {
	if (key === 'search') { delete filters.value.search; quickSearch.value = ''; }
	else if (key === 'protocol') { delete filters.value.protocol; quickProtocol.value = null; }
	else if (key === 'kind') { delete filters.value.kind; advancedFilterValues.value.kind = null; }
	else if (key === 'direction') { delete filters.value.direction; advancedFilterValues.value.direction = null; }
	syncStore();
}

function clearAllFilters(): void {
	filters.value = {};
	quickSearch.value = '';
	quickProtocol.value = null;
	advancedFilterValues.value = { kind: null, direction: null };
	hasPendingAdvancedFilters.value = false;
	syncStore();
}

/**
 * Update column visibility from the list header menu.
 * @param {ListHeaderMenuColumn[]} columns - the updated columns
 */
function handleColumnsUpdate(columns: ListHeaderMenuColumn[]): void {
	for (const col of columns) {
		if (col.key === 'direction') columnVisibility.value.direction = col.visible;
		if (col.key === 'kind') columnVisibility.value.kind = col.visible;
		if (col.key === 'status') columnVisibility.value.status = col.visible;
	}
}

/**
 * Change page.
 * @param {number} page - the new page
 */
function handlePageChange(page: number): void {
	logsStore.setPage(page);
}

/**
 * Change items per page.
 * @param {number} value - the new value
 */
function handleItemsPerPageChange(value: number): void {
	logsStore.setItemsPerPage(value);
}

/** LIFECYCLE HOOKS */
onMounted(() => {
	void logsStore.fetch();
});
</script>

<style scoped lang="scss">
.filter-input {
	:deep(.q-field__control) {
		border-radius: var(--mapex-radius-md);
	}
}

// The summary column carries no fixed width. With the table's default auto layout
// it grows to fit its (nowrap) content and forces a horizontal scrollbar, cutting
// off the trailing columns. A fixed layout makes the summary take the remaining
// space and truncate with the shared cell ellipsis instead of widening the row.
:deep(.data-row-table .q-table) {
	table-layout: fixed;
	width: 100%;
}

// The avatar is an inline-flex box; without this the cell's text-overflow
// rule renders an ellipsis ("...") under the icon when it overflows.
:deep(.data-row-cell--avatar) {
	overflow: visible;
	text-overflow: clip;
}

// These rows have no actions menu, so the last column (the timestamp) is the
// final cell. Give it a right gap that mirrors the icon's left inset. The
// extra specificity overrides DataRow's own td:last-child padding rule.
:deep(.data-row-table .data-row-card td:last-child) {
	padding-right: var(--mapex-spacing-lg) !important;
}
</style>
