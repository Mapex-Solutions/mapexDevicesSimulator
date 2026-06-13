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
					:items-count="logsStore.items.length"
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
				<DataRow :data="log" :columns="visibleColumns" :show-actions="false" @click="openDetail(log)" />
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

		<!-- Pagination: cursor-based, so previous/next only (no total page count) -->
		<div v-if="logsStore.items.length" class="row justify-center items-center q-gutter-sm q-mt-lg q-mb-lg">
			<q-btn
				flat
				dense
				no-caps
				icon="chevron_left"
				:label="t('logs.previous')"
				:disable="!logsStore.hasPrev || logsStore.loading"
				@click="logsStore.prev"
			/>
			<q-btn
				flat
				dense
				no-caps
				icon-right="chevron_right"
				:label="t('logs.next')"
				:disable="!logsStore.hasNext || logsStore.loading"
				@click="logsStore.next"
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

		<!-- Log detail: the same drawer the console uses, showing payload + response -->
		<GenericDrawer
			v-model="detailOpen"
			:title="t('console.details')"
			icon="mdi-text-box-search-outline"
			:close-tooltip="t('common.close')"
		>
			<MessageDetail :message="selectedMessage" />
		</GenericDrawer>
	</q-page>
</template>

<script setup lang="ts">
/** TYPE IMPORTS */
import type { ListHeaderMenuColumn } from '@components/ListHeaderMenu';
import type { DataRowColumn } from '@components/DataRow';
import type { FilterField, FilterValues, FilterAutocompleteOption } from '@components/AdvancedFiltersDrawer';
import type { ConsoleMessage } from '@stores/messages';
import type { Log, LogDirection, ProtocolId } from '@services/sim';

/** VUE IMPORTS */
import { computed, onMounted, ref } from 'vue';

/** COMPONENTS */
import { PageHeader } from '@components/PageHeader';
import { ListHeaderMenu } from '@components/ListHeaderMenu';
import { DataRow } from '@components/DataRow';
import { ListCardEmpty } from '@components/ListCardEmpty';
import { AdvancedFiltersDrawer } from '@components/AdvancedFiltersDrawer';
import { AppTooltip } from '@components/AppTooltip';
import { GenericDrawer } from '@components/GenericDrawer';
import { MessageDetail } from '@components/MessageDetail';
import { protocolIcon } from '@components/protocols/ProtocolRegistry';

/** COMPOSABLES */
import { useTranslations } from '@composables/i18n';

/** STORES */
import { useLogsStore } from '@stores/logs';
import { useDevicesStore } from '@stores/devices';

const MAX_VISIBLE_CHIPS = 2;

/** COMPOSABLES & STORES */
const { t } = useTranslations();
const logsStore = useLogsStore();
const devicesStore = useDevicesStore();

/** STATE */
const quickSearch = ref('');
const quickProtocol = ref<string | null>(null);
const showFiltersDrawer = ref(false);
const advancedFilterValues = ref<FilterValues>({ kind: null, direction: null, device: null, event: null, dateFrom: null, dateTo: null });
const hasPendingAdvancedFilters = ref(false);
const filters = ref<{ search?: string; protocol?: string; kind?: string; direction?: string; device?: string; event?: string; dateFrom?: string; dateTo?: string }>({});
const columnVisibility = ref({ direction: true, kind: true, status: true });
const detailOpen = ref(false);
const selectedLog = ref<Log | null>(null);

/** COMPUTED */
const protocolOptions = computed(() => [
	{ label: t('console.allDevices'), value: null },
	...(['http', 'mqtt', 'lorawan', 'basicstation'] as ProtocolId[]).map((p) => ({ label: t(`protocol.${p}`), value: p })),
]);

// Ordered by what someone scanning the history reaches for first: when it
// happened, then which device, then which event, and only then the technical
// type/direction filters.
const advancedFilterFields = computed<FilterField[]>(() => [
	{ key: 'dateFrom', type: 'input', label: t('logs.filter.dateFrom'), icon: 'event', inputType: 'date' },
	{ key: 'dateTo', type: 'input', label: t('logs.filter.dateTo'), icon: 'event', inputType: 'date' },
	{
		key: 'device', type: 'autocomplete', label: t('logs.filter.device'), icon: 'memory',
		placeholder: t('logs.filter.devicePlaceholder'), fetchOptions: fetchDeviceOptions,
	},
	{ key: 'event', type: 'input', label: t('logs.filter.event'), icon: 'bolt', placeholder: t('logs.filter.eventPlaceholder') },
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

const hasActiveFilters = computed(() => Object.keys(filters.value).length > 0);

const advancedFiltersCount = computed(() => {
	const keys = ['kind', 'direction', 'device', 'event', 'dateFrom', 'dateTo'] as const;
	return keys.filter((k) => filters.value[k]).length;
});

const activeFilterChips = computed(() => {
	const chips: { key: string; label: string; value: string }[] = [];
	if (filters.value.search) chips.push({ key: 'search', label: t('console.filter.search'), value: filters.value.search });
	if (filters.value.protocol) chips.push({ key: 'protocol', label: t('console.filter.protocol'), value: t(`protocol.${filters.value.protocol}`) });
	if (filters.value.dateFrom) chips.push({ key: 'dateFrom', label: t('logs.filter.dateFrom'), value: filters.value.dateFrom });
	if (filters.value.dateTo) chips.push({ key: 'dateTo', label: t('logs.filter.dateTo'), value: filters.value.dateTo });
	if (filters.value.device) chips.push({ key: 'device', label: t('logs.filter.device'), value: filters.value.device });
	if (filters.value.event) chips.push({ key: 'event', label: t('logs.filter.event'), value: filters.value.event });
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
	{ key: 'icon', label: '', type: 'avatar', visible: 'always', width: 72, icon: (_v, row) => protocolIcon(row.protocol), color: () => 'primary', tooltip: (_v, row) => t(`protocol.${row.protocol}`) },
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

// Map the selected log into the console message shape so the shared detail drawer
// renders it (payload + response + status), with the event name surfaced as meta.
const selectedMessage = computed<ConsoleMessage | null>(() => {
	const log = selectedLog.value;
	if (!log) return null;
	const message: ConsoleMessage = {
		id: log.id,
		ts: log.created ?? '',
		protocol: log.protocol,
		deviceId: log.deviceId,
		deviceName: log.deviceName,
		direction: log.direction,
		kind: log.kind,
		summary: log.summary,
		payload: log.payload,
	};
	if (log.status) message.status = log.status;
	if (log.response) message.response = log.response;
	if (log.eventName) message.meta = { event: log.eventName };
	return message;
});

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
	if (filters.value.device) map.device = filters.value.device;
	if (filters.value.event) map.event = filters.value.event;
	// The date inputs are day-precision; widen them to the full day so the bounds
	// are inclusive against the timestamp the engine stores.
	if (filters.value.dateFrom) map.dateFrom = `${filters.value.dateFrom}T00:00:00Z`;
	if (filters.value.dateTo) map.dateTo = `${filters.value.dateTo}T23:59:59Z`;
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
	applyDrawerValue('kind', values.kind);
	applyDrawerValue('direction', values.direction);
	applyDrawerValue('device', values.device);
	applyDrawerValue('event', values.event);
	applyDrawerValue('dateFrom', values.dateFrom);
	applyDrawerValue('dateTo', values.dateTo);
	hasPendingAdvancedFilters.value = false;
	syncStore();
}

/**
 * Copy one drawer value into the active filters, or drop it when empty.
 * @param {keyof typeof filters.value} key - the filter key
 * @param {unknown} value - the drawer value
 */
function applyDrawerValue(key: 'kind' | 'direction' | 'device' | 'event' | 'dateFrom' | 'dateTo', value: unknown): void {
	if (value) filters.value[key] = String(value);
	else delete filters.value[key];
}

/**
 * Track pending (unapplied) changes from the drawer.
 * @param {boolean} pending - whether there are pending changes
 */
function handlePendingChange(pending: boolean): void {
	hasPendingAdvancedFilters.value = pending;
}

function handleAdvancedFiltersReset(): void {
	advancedFilterValues.value = { kind: null, direction: null, device: null, event: null, dateFrom: null, dateTo: null };
}

/**
 * Remove a single active filter.
 * @param {string} key - the filter key
 */
function removeFilter(key: string): void {
	if (key === 'search') { delete filters.value.search; quickSearch.value = ''; }
	else if (key === 'protocol') { delete filters.value.protocol; quickProtocol.value = null; }
	else {
		delete filters.value[key as keyof typeof filters.value];
		advancedFilterValues.value[key] = null;
	}
	syncStore();
}

function clearAllFilters(): void {
	filters.value = {};
	quickSearch.value = '';
	quickProtocol.value = null;
	handleAdvancedFiltersReset();
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
 * Change items per page.
 * @param {number} value - the new value
 */
function handleItemsPerPageChange(value: number): void {
	logsStore.setItemsPerPage(value);
}

/**
 * Search the configured devices for the deviceId filter autocomplete.
 * @param {string} search - the typed term
 * @returns {Promise<FilterAutocompleteOption[]>} the matching devices
 */
async function fetchDeviceOptions(search: string): Promise<FilterAutocompleteOption[]> {
	const term = search.trim().toLowerCase();
	return devicesStore.items
		.filter((device) => !term || device.name.toLowerCase().includes(term) || device.deviceId.toLowerCase().includes(term))
		.slice(0, 20)
		.map((device) => ({ id: device.deviceId, label: device.name, caption: device.deviceId }));
}

/**
 * Open the detail drawer on a log row.
 * @param {Log} log - the clicked log
 */
function openDetail(log: Log): void {
	selectedLog.value = log;
	detailOpen.value = true;
}

/** LIFECYCLE HOOKS */
onMounted(() => {
	logsStore.fetch();
	void devicesStore.fetch();
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
