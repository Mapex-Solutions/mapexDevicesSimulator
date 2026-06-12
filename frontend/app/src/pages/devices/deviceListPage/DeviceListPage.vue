<template>
	<q-page class="q-pt-lg">
		<PageHeader
			icon="memory"
			icon-color="primary"
			:title="t('devices.title')"
			:description="t('devices.subtitle')"
			:button="{ label: t('devices.newDevice'), icon: 'add', color: 'primary', onClick: openCreate }"
		/>

		<!-- Results header -->
		<div class="row items-center q-pt-md q-mb-md">
			<div class="col">
				<div class="row items-center">
					<q-icon name="memory" size="sm" color="primary" class="q-mr-sm" />
					<div class="text-subtitle1 text-weight-medium text-primary">{{ t('devices.title') }}</div>
				</div>
			</div>
			<div class="col-auto">
				<ListHeaderMenu
					icon="memory"
					:item-label="t('devices.itemLabel')"
					:item-label-plural="t('devices.itemLabelPlural')"
					:items-count="visibleDevices.length"
					:items-per-page="visibleDevices.length || 1"
					:columns="menuColumns"
					:show-items-per-page="false"
					:refreshing="devicesStore.loading"
					@update:columns="handleColumnsUpdate"
					@refresh="devicesStore.fetch"
				/>
			</div>
		</div>

		<!-- Rows -->
		<div v-if="devicesStore.loading" class="row justify-center q-my-lg">
			<q-spinner color="primary" size="3em" />
		</div>

		<div v-else class="row">
			<div v-for="device in visibleDevices" :key="device.id" class="col-12 q-mb-xs">
				<DataRow
					:data="device"
					:columns="visibleColumns"
					:actions="{ showView: false }"
					@edit="openEdit"
					@delete="onDelete"
				/>
			</div>

			<div v-if="!visibleDevices.length" class="col-12">
				<ListCardEmpty
					icon="mdi-chip"
					:title="t('devices.emptyTitle')"
					:description="t('devices.empty')"
					:button-label="t('devices.newDevice')"
					button-icon="add"
					@button-click="openCreate"
				/>
			</div>
		</div>
	</q-page>
</template>

<script setup lang="ts">
/** TYPE IMPORTS */
import type { Device } from '@services/sim';
import type { DataRowColumn } from '@components/DataRow';
import type { ListHeaderMenuColumn } from '@components/ListHeaderMenu';

/** VUE IMPORTS */
import { computed, onMounted, ref } from 'vue';

/** COMPONENTS */
import { PageHeader } from '@components/PageHeader';
import { ListHeaderMenu } from '@components/ListHeaderMenu';
import { DataRow } from '@components/DataRow';
import { ListCardEmpty } from '@components/ListCardEmpty';

/** COMPOSABLES */
import { useTranslations } from '@composables/i18n';

/** UTILS */
import { useQuasar } from 'quasar';

/** SERVICES */
import { useRoute, useRouter } from 'vue-router';

/** STORES */
import { useDevicesStore } from '@stores/devices';

/** COMPOSABLES & STORES */
const { t } = useTranslations();
const $q = useQuasar();
const route = useRoute();
const router = useRouter();
const devicesStore = useDevicesStore();

const PROTOCOL_ICON: Record<string, string> = {
	http: 'mdi-web',
	mqtt: 'mdi-transit-connection-variant',
	lorawan: 'mdi-access-point',
	basicstation: 'mdi-radio-tower',
};

/** STATE */
const columnVisibility = ref({ protocol: true, target: true, events: true, created: true, status: true });

/** COMPUTED */
const visibleDevices = computed(() => {
	const filter = route.query.protocol;
	if (typeof filter !== 'string') return devicesStore.items;
	return devicesStore.items.filter((device) => device.protocolId === filter);
});

const menuColumns = computed<ListHeaderMenuColumn[]>(() => [
	{ key: 'protocol', label: t('devices.col.protocol'), visible: columnVisibility.value.protocol },
	{ key: 'target', label: t('devices.col.target'), visible: columnVisibility.value.target },
	{ key: 'events', label: t('devices.col.events'), visible: columnVisibility.value.events },
	{ key: 'created', label: t('devices.col.created'), visible: columnVisibility.value.created },
	{ key: 'status', label: t('devices.col.status'), visible: columnVisibility.value.status },
]);

const deviceColumns = computed<DataRowColumn[]>(() => [
	{ key: 'icon', label: '', type: 'avatar', visible: 'always', width: 72, icon: (_v, row) => PROTOCOL_ICON[row.protocolId] ?? 'mdi-chip', color: () => 'primary', tooltip: (_v, row) => t(`protocol.${row.protocolId}`) },
	{ key: 'name', label: t('devices.col.name'), type: 'text', visible: 'always', width: 240, ellipsis: true, secondaryKey: 'deviceId' },
	{ key: 'protocolId', label: t('devices.col.protocol'), type: 'chip', visible: 'laptop', width: 120, format: (v) => t(`protocol.${v}`), color: () => 'primary' },
	{ key: 'config', label: t('devices.col.target'), type: 'text', visible: 'laptop', ellipsis: true, format: (_v, row) => targetOf(row) },
	{ key: 'events', label: t('devices.col.events'), type: 'badge', visible: 'laptop', width: 90, format: (v) => String((v as unknown[])?.length ?? 0), color: () => 'primary' },
	{ key: 'enabled', label: t('devices.col.status'), type: 'chip', visible: 'always', width: 110, format: (v) => (v ? t('devices.on') : t('devices.off')), color: (v) => (v ? 'positive' : 'grey') },
	{ key: 'created', label: t('devices.col.created'), type: 'text', visible: 'laptop', width: 130, format: (v) => formatDate(v as string) },
]);

const visibleColumns = computed(() =>
	deviceColumns.value.filter((col) => {
		if (col.key === 'protocolId') return columnVisibility.value.protocol;
		if (col.key === 'config') return columnVisibility.value.target;
		if (col.key === 'events') return columnVisibility.value.events;
		if (col.key === 'created') return columnVisibility.value.created;
		if (col.key === 'enabled') return columnVisibility.value.status;
		return true;
	}),
);

/** FUNCTIONS */

/**
 * The device's target address per protocol.
 * @param {Device} device - the device to describe
 * @returns {string} the endpoint URL, broker URL or LoRaWAN summary
 */
function targetOf(device: Device): string {
	const config = device.config;
	if (config.kind === 'mqtt') return config.brokerUrl;
	if (config.kind === 'lorawan') return `${config.region} · ${config.macVersion}`;
	if (config.kind === 'basicstation') return config.lnsUri;
	return config.url;
}

/**
 * Format an ISO timestamp as a short local date.
 * @param {string} iso - the ISO timestamp
 * @returns {string} the formatted date
 */
function formatDate(iso: string): string {
	if (!iso) return '—';
	const date = new Date(iso);
	return Number.isNaN(date.getTime()) ? '—' : date.toLocaleDateString();
}

/**
 * Update column visibility from the list header menu.
 * @param {ListHeaderMenuColumn[]} columns - the updated columns
 */
function handleColumnsUpdate(columns: ListHeaderMenuColumn[]): void {
	for (const col of columns) {
		if (col.key === 'protocol') columnVisibility.value.protocol = col.visible;
		if (col.key === 'target') columnVisibility.value.target = col.visible;
		if (col.key === 'events') columnVisibility.value.events = col.visible;
		if (col.key === 'created') columnVisibility.value.created = col.visible;
		if (col.key === 'status') columnVisibility.value.status = col.visible;
	}
}

function openCreate(): void {
	void router.push({ name: 'device-new' });
}

/**
 * Open the device wizard in edit mode.
 * @param {Device} device - the device to edit
 */
function openEdit(device: Device): void {
	void router.push({ name: 'device-edit', params: { id: device.id } });
}

/**
 * Confirm and delete a device.
 * @param {Device} device - the device to delete
 */
function onDelete(device: Device): void {
	$q.dialog({
		title: t('common.delete'),
		message: t('common.deleteConfirm', { name: device.name }),
		cancel: true,
	}).onOk(async () => {
		try {
			await devicesStore.remove(device.id);
		} catch (err) {
			$q.notify({ type: 'negative', message: err instanceof Error ? err.message : t('common.deleteFailed') });
		}
	});
}

/** LIFECYCLE HOOKS */
onMounted(() => {
	void devicesStore.fetch();
});
</script>

<style scoped lang="scss">
// The avatar is an inline-flex box; without this the cell's text-overflow
// rule renders an ellipsis ("...") under the icon when it overflows.
:deep(.data-row-cell--avatar) {
	overflow: visible;
	text-overflow: clip;
}
</style>
