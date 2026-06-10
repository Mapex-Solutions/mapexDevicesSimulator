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
					:actions="{ showView: false }"
					@edit="openEdit"
					@delete="onDelete"
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
import type { Gateway } from '@services/sim';
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
import { useRouter } from 'vue-router';

/** STORES */
import { useGatewaysStore } from '@stores/gateways';

/** COMPOSABLES & STORES */
const { t } = useTranslations();
const $q = useQuasar();
const router = useRouter();
const gatewaysStore = useGatewaysStore();

const LINK_ICON: Record<string, string> = {
	basicstation: 'mdi-access-point',
	udp: 'mdi-lan-connect',
};

/** STATE */
const columnVisibility = ref({ link: true, region: true, status: true, created: true });

/** COMPUTED */
const gateways = computed(() => gatewaysStore.items);

const menuColumns = computed<ListHeaderMenuColumn[]>(() => [
	{ key: 'link', label: t('gateways.col.link'), visible: columnVisibility.value.link },
	{ key: 'region', label: t('gateways.col.region'), visible: columnVisibility.value.region },
	{ key: 'status', label: t('gateways.col.status'), visible: columnVisibility.value.status },
	{ key: 'created', label: t('gateways.col.created'), visible: columnVisibility.value.created },
]);

const gatewayColumns = computed<DataRowColumn[]>(() => [
	{ key: 'link.protocol', label: '', type: 'avatar', visible: 'always', width: 72, icon: (v) => LINK_ICON[v as string] ?? 'mdi-access-point', color: () => 'primary', tooltip: (v) => t(`gatewayLink.${v}`) },
	{ key: 'name', label: t('gateways.col.name'), type: 'text', visible: 'always', width: 240, ellipsis: true, secondaryKey: 'eui' },
	{ key: 'link', label: t('gateways.col.link'), type: 'text', visible: 'laptop', ellipsis: true, format: (_v, row) => linkTargetOf(row) },
	{ key: 'region', label: t('gateways.col.region'), type: 'chip', visible: 'always', width: 110, color: () => 'primary' },
	{ key: 'enabled', label: t('gateways.col.status'), type: 'chip', visible: 'always', width: 110, format: (v) => (v ? t('devices.on') : t('devices.off')), color: (v) => (v ? 'positive' : 'grey') },
	{ key: 'created', label: t('gateways.col.created'), type: 'text', visible: 'laptop', width: 130, format: (v) => formatDate(v as string) },
]);

const visibleColumns = computed(() =>
	gatewayColumns.value.filter((col) => {
		if (col.key === 'link') return columnVisibility.value.link;
		if (col.key === 'region') return columnVisibility.value.region;
		if (col.key === 'enabled') return columnVisibility.value.status;
		if (col.key === 'created') return columnVisibility.value.created;
		return true;
	}),
);

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
		if (col.key === 'link') columnVisibility.value.link = col.visible;
		if (col.key === 'region') columnVisibility.value.region = col.visible;
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
 * Confirm and delete a gateway.
 * @param {Gateway} gateway - the gateway to delete
 */
function onDelete(gateway: Gateway): void {
	$q.dialog({
		title: t('common.delete'),
		message: t('common.deleteConfirm', { name: gateway.name }),
		cancel: true,
	}).onOk(() => {
		void gatewaysStore.remove(gateway.id);
	});
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
