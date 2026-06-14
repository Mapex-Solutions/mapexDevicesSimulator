<template>
	<q-page class="q-pt-lg">
		<PageHeader
			icon="router"
			icon-color="primary"
			:title="t('gateways.title')"
			:description="t('gateways.subtitle')"
			:button="{ label: t('gateways.newGateway'), icon: 'add', color: 'primary', onClick: openCreate }"
		/>

		<!-- Results header -->
		<div class="row items-center q-pt-md q-mb-md">
			<div class="col">
				<div class="row items-center">
					<q-icon name="router" size="sm" color="primary" class="q-mr-sm" />
					<div class="text-subtitle1 text-weight-medium text-primary">{{ t('gateways.title') }}</div>
				</div>
			</div>
			<div class="col-auto">
				<ListHeaderMenu
					icon="router"
					:item-label="t('gateways.itemLabel')"
					:item-label-plural="t('gateways.itemLabelPlural')"
					:items-count="gateways.length"
					:items-per-page="gateways.length || 1"
					:columns="menuColumns"
					:show-items-per-page="false"
					:refreshing="gatewaysStore.loading"
					@update:columns="handleColumnsUpdate"
					@refresh="gatewaysStore.fetch"
				/>
			</div>
		</div>

		<!-- Rows -->
		<div v-if="gatewaysStore.loading" class="row justify-center q-my-lg">
			<q-spinner color="primary" size="3em" />
		</div>

		<div v-else class="row">
			<div v-for="gateway in gateways" :key="gateway.id" class="col-12 q-mb-xs">
				<DataRow
					:data="gateway"
					:columns="visibleColumns"
					:actions="rowActions"
					@edit="openEdit"
					@delete="onDelete"
					@action="onAction"
				/>
			</div>

			<div v-if="!gateways.length" class="col-12">
				<ListCardEmpty
					icon="mdi-access-point"
					:title="t('gateways.emptyTitle')"
					:description="t('gateways.empty')"
					:button-label="t('gateways.newGateway')"
					button-icon="add"
					@button-click="openCreate"
				/>
			</div>
		</div>
	</q-page>
</template>

<script setup lang="ts">
/** TYPE IMPORTS */
import type { Gateway, GatewayInput } from '@services/sim';
import type { DataRowColumn, DataRowActionConfig } from '@components/DataRow';
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
import { useDateFormat } from '@composables/datetime';
import { useGatewayConnections } from '@composables/gateways';

/** UTILS */
import { notifySuccess, notifyFail, dialogDelete } from '@utils/alert';

/** SERVICES */
import { useRouter } from 'vue-router';

/** STORES */
import { useGatewaysStore } from '@stores/gateways';

/** COMPOSABLES & STORES */
const { t } = useTranslations();
const { formatDate } = useDateFormat();
const router = useRouter();
const gatewaysStore = useGatewaysStore();
const { connectionOf } = useGatewayConnections();

const LINK_ICON: Record<string, string> = {
	basicstation: 'mdi-access-point',
	udp: 'mdi-lan-connect',
};

const CONNECTION_COLOR: Record<string, string> = {
	online: 'positive',
	connecting: 'warning',
	offline: 'negative',
	unknown: 'grey',
};

/** STATE */
const columnVisibility = ref({ link: true, region: true, connection: true, status: true, created: true });

/** COMPUTED */
const gateways = computed(() =>
	gatewaysStore.items.map((gateway) => ({ ...gateway, connection: connectionOf(gateway.id) })),
);

const menuColumns = computed<ListHeaderMenuColumn[]>(() => [
	{ key: 'link', label: t('gateways.col.link'), visible: columnVisibility.value.link },
	{ key: 'region', label: t('gateways.col.region'), visible: columnVisibility.value.region },
	{ key: 'connection', label: t('gateways.col.connection'), visible: columnVisibility.value.connection },
	{ key: 'status', label: t('gateways.col.status'), visible: columnVisibility.value.status },
	{ key: 'created', label: t('gateways.col.created'), visible: columnVisibility.value.created },
]);

const gatewayColumns = computed<DataRowColumn[]>(() => [
	{ key: 'link.protocol', label: '', type: 'avatar', visible: 'always', width: 72, icon: (v) => LINK_ICON[v as string] ?? 'mdi-access-point', color: () => 'primary', tooltip: (v) => t(`gatewayLink.${v}`) },
	{ key: 'name', label: t('gateways.col.name'), type: 'text', visible: 'always', width: 240, ellipsis: true, secondaryKey: 'eui' },
	{ key: 'link', label: t('gateways.col.link'), type: 'text', visible: 'laptop', ellipsis: true, format: (_v, row) => linkTargetOf(row) },
	{ key: 'region', label: t('gateways.col.region'), type: 'chip', visible: 'always', width: 110, color: () => 'primary' },
	{ key: 'connection', label: t('gateways.col.connection'), type: 'chip', visible: 'always', width: 120, format: (v) => t(`gatewayConn.${v as string}`), color: (v) => CONNECTION_COLOR[v as string] ?? 'grey' },
	{ key: 'enabled', label: t('gateways.col.status'), type: 'chip', visible: 'always', width: 110, format: (v) => (v ? t('devices.on') : t('devices.off')), color: (v) => (v ? 'positive' : 'grey') },
	{ key: 'created', label: t('gateways.col.created'), type: 'text', visible: 'laptop', width: 130, format: (v) => formatDate(v as string) },
]);

const visibleColumns = computed(() =>
	gatewayColumns.value.filter((col) => {
		if (col.key === 'link') return columnVisibility.value.link;
		if (col.key === 'region') return columnVisibility.value.region;
		if (col.key === 'connection') return columnVisibility.value.connection;
		if (col.key === 'enabled') return columnVisibility.value.status;
		if (col.key === 'created') return columnVisibility.value.created;
		return true;
	}),
);

const rowActions = computed<DataRowActionConfig>(() => ({
	showView: false,
	customActions: [{ key: 'duplicate', label: t('common.duplicate'), icon: 'content_copy' }],
}));

/** FUNCTIONS */

/**
 * The gateway's LNS connection address per link protocol.
 * @param {Gateway} gateway - the gateway to describe
 * @returns {string} the LNS URI or host:port
 */
function linkTargetOf(gateway: Gateway): string {
	return gateway.link.protocol === 'udp' ? `${gateway.link.host}:${gateway.link.port}` : gateway.link.lnsUri;
}

/**
 * Update column visibility from the list header menu.
 * @param {ListHeaderMenuColumn[]} columns - the updated columns
 */
function handleColumnsUpdate(columns: ListHeaderMenuColumn[]): void {
	for (const col of columns) {
		if (col.key === 'link') columnVisibility.value.link = col.visible;
		if (col.key === 'region') columnVisibility.value.region = col.visible;
		if (col.key === 'connection') columnVisibility.value.connection = col.visible;
		if (col.key === 'status') columnVisibility.value.status = col.visible;
		if (col.key === 'created') columnVisibility.value.created = col.visible;
	}
}

function openCreate(): void {
	void router.push({ name: 'gateway-new' });
}

/**
 * Open the gateway wizard in edit mode.
 * @param {Gateway} gateway - the gateway to edit
 */
function openEdit(gateway: Gateway): void {
	void router.push({ name: 'gateway-edit', params: { id: gateway.id } });
}

/**
 * Handle a custom row action.
 * @param {string} key - the action key
 * @param {Gateway} gateway - the row's gateway
 */
function onAction(key: string, gateway: Gateway): void {
	if (key === 'duplicate') void duplicateGateway(gateway);
}

/**
 * Mint a 16-hex EUI not already in use, so a copied gateway has its own identity
 * on the LNS instead of clashing with the original.
 * @returns {string} a unique gateway EUI
 */
function uniqueEui(): string {
	const taken = new Set(gatewaysStore.items.map((gateway) => gateway.eui.toUpperCase()));
	let candidate = randomEui();
	while (taken.has(candidate)) candidate = randomEui();
	return candidate;
}

/**
 * Generate a random 16-hex-character gateway EUI.
 * @returns {string} the EUI
 */
function randomEui(): string {
	const bytes = new Uint8Array(8);
	crypto.getRandomValues(bytes);
	return Array.from(bytes, (b) => b.toString(16).padStart(2, '0')).join('').toUpperCase();
}

/**
 * Duplicate a gateway through the create endpoint, copying every field, renaming
 * the copy, and minting a fresh EUI so it is a distinct gateway.
 * @param {Gateway} gateway - the gateway to duplicate
 */
async function duplicateGateway(gateway: Gateway): Promise<void> {
	// Deep-clone through JSON: the row comes from the reactive store, so its nested
	// link is a Vue proxy that structuredClone rejects; a DTO is plain JSON.
	const source = JSON.parse(JSON.stringify(gateway)) as Gateway;
	const input: GatewayInput = {
		name: t('common.duplicateName', { name: gateway.name }),
		eui: uniqueEui(),
		enabled: source.enabled,
		region: source.region,
		description: source.description,
		link: source.link,
	};
	try {
		await gatewaysStore.create(input);
		notifySuccess({ message: t('common.duplicated') });
	} catch {
		notifyFail({ message: t('common.saveFailed') });
	}
}

/**
 * Confirm and delete a gateway.
 * @param {Gateway} gateway - the gateway to delete
 */
async function onDelete(gateway: Gateway): Promise<void> {
	const confirmed = await dialogDelete({
		title: t('common.delete'),
		message: t('common.deleteConfirm', { name: gateway.name }),
		ok: { label: t('common.delete'), color: 'negative' },
		cancel: { label: t('common.cancel'), flat: true },
	});
	if (!confirmed) return;
	try {
		await gatewaysStore.remove(gateway.id);
	} catch (err) {
		notifyFail({ message: err instanceof Error ? err.message : t('common.deleteFailed') });
	}
}

/** LIFECYCLE HOOKS */
onMounted(() => {
	void gatewaysStore.fetch();
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
