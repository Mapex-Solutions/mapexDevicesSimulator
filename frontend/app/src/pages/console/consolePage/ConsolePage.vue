<template>
	<q-page class="console-page">
		<PageHeader
			icon="terminal"
			icon-color="primary"
			:title="t('console.title')"
			:description="t('console.subtitle')"
		/>

		<div class="console">
			<!-- Devices -->
			<aside class="console__devices">
				<div class="console__pane-head row items-center justify-between">
					<span>{{ t('console.devices') }}</span>
					<div class="row items-center q-gutter-xs">
						<q-badge :color="onlineGateways ? 'positive' : 'grey-6'" text-color="white" class="console__stat">
							<q-icon name="router" size="13px" class="q-mr-xs" />{{ pad(onlineGateways) }}
							<AppTooltip :content="t('console.onlineGateways')" />
						</q-badge>
						<q-badge :color="onlineDevices ? 'positive' : 'grey-6'" text-color="white" class="console__stat">
							<q-icon name="memory" size="13px" class="q-mr-xs" />{{ pad(onlineDevices) }}
							<AppTooltip :content="t('console.onlineDevices')" />
						</q-badge>
					</div>
				</div>

				<div class="console__search row items-center no-wrap q-gutter-xs">
					<q-input class="col" v-model="search" dense outlined hide-bottom-space :placeholder="t('console.searchDevices')" clearable>
						<template #prepend><q-icon name="search" /></template>
					</q-input>
					<q-btn flat dense round icon="mdi-filter-variant" :color="deviceProtocolFilter ? 'primary' : undefined">
						<q-badge v-if="deviceProtocolFilter" floating rounded color="primary" />
						<AppTooltip>{{ t('console.filter.protocol') }}</AppTooltip>
						<q-menu anchor="bottom right" self="top right">
							<q-list dense style="min-width: 180px">
								<q-item
									clickable
									v-close-popup
									:active="!deviceProtocolFilter"
									active-class="console__device--active"
									@click="deviceProtocolFilter = null"
								>
									<q-item-section avatar><q-icon name="mdi-all-inclusive" /></q-item-section>
									<q-item-section>{{ t('console.allDevices') }}</q-item-section>
								</q-item>
								<q-item
									v-for="opt in protocolFilterOptions"
									:key="opt.value"
									clickable
									v-close-popup
									:active="deviceProtocolFilter === opt.value"
									active-class="console__device--active"
									@click="deviceProtocolFilter = opt.value"
								>
									<q-item-section avatar><q-icon :name="protocolIcon(opt.value)" /></q-item-section>
									<q-item-section>{{ opt.label }}</q-item-section>
								</q-item>
							</q-list>
						</q-menu>
					</q-btn>
				</div>

				<q-list class="console__device-list">
					<q-item
						clickable
						:active="!messagesStore.deviceFilter"
						active-class="console__device--active"
						@click="messagesStore.setDeviceFilter(null)"
					>
						<q-item-section avatar><q-icon name="mdi-all-inclusive" /></q-item-section>
						<q-item-section>{{ t('console.allDevices') }}</q-item-section>
					</q-item>

					<q-item
						v-for="device in deviceList"
						:key="device.deviceId"
						clickable
						:active="messagesStore.deviceFilter === device.deviceId"
						active-class="console__device--active"
						@click="messagesStore.setDeviceFilter(device.deviceId)"
					>
						<q-item-section avatar><q-icon :name="protocolIcon(device.protocol)" /></q-item-section>
						<q-item-section>
							<q-item-label>{{ device.name }}</q-item-label>
							<q-item-label caption>{{ t(`protocol.${device.protocol}`) }}</q-item-label>
						</q-item-section>
						<q-item-section v-if="device.uuid" side>
							<q-btn flat dense round size="sm" color="primary" icon="mdi-flash" @click.stop="openFire(device.uuid)">
								<AppTooltip :content="t('console.fireEvent')" />
							</q-btn>
						</q-item-section>
					</q-item>

					<div v-if="!deviceList.length" class="console__devices-empty">{{ t('console.emptyDevices') }}</div>
				</q-list>

				<div class="console__devices-foot">
					<q-btn flat dense no-caps icon="mdi-plus" :label="t('dashboard.newDevice')" :to="{ name: 'device-new' }" />
				</div>
			</aside>

			<!-- Message log -->
			<main class="console__log">
				<div class="console__pane-head row items-center justify-between">
					<div class="row items-center q-gutter-sm">
						<span>{{ t('console.title') }}</span>
						<q-badge color="grey-7" :label="String(ordered.length)" />
					</div>
					<div class="row q-gutter-xs">
						<q-btn
							flat
							dense
							no-caps
							icon="mdi-filter-variant"
							:label="t('console.filters')"
							:color="activeFilterCount ? 'primary' : undefined"
							@click="filterOpen = true"
						>
							<q-badge v-if="activeFilterCount" floating rounded color="primary" :label="String(activeFilterCount)" />
						</q-btn>
						<q-btn flat dense no-caps icon="mdi-broom" :label="t('console.clear')" @click="messagesStore.clear()" />
					</div>
				</div>

				<div v-if="!ordered.length" class="console__empty">{{ t('console.emptyLog') }}</div>

				<div v-else class="console__rows">
					<button
						v-for="message in ordered"
						:key="message.id"
						type="button"
						class="msg"
						:class="{ 'msg--active': message.id === messagesStore.selectedId }"
						@click="selectMessage(message.id)"
					>
						<q-icon :name="dirIcon(message.direction)" :color="dirColor(message.direction)" size="18px" class="msg__dir" />
						<span class="msg__ts">{{ message.ts }}</span>
						<q-badge outline color="grey" :label="t(`protocol.${message.protocol}`)" class="msg__proto" />
						<span class="msg__device">{{ message.deviceName }}</span>
						<span class="msg__summary">{{ message.summary }}</span>
						<q-badge v-if="message.status" :color="statusColor(message.status)" text-color="white" :label="message.status" class="msg__status" />
					</button>
				</div>
			</main>
		</div>

		<!-- Message detail (opens on click, closable) -->
		<GenericDrawer
			v-model="detailOpen"
			:title="t('console.details')"
			icon="mdi-text-box-search-outline"
			:close-tooltip="t('common.close')"
			@close="onDetailClose"
		>
			<MessageDetail :message="messagesStore.selected" />
		</GenericDrawer>

		<GenericModal v-model="filterOpen" :title="t('console.filters')" icon="mdi-filter-variant">
			<MessageFilterBar :protocol="activeProtocol" v-model="filterValues" />
			<template #footer>
				<q-btn flat :label="t('console.filter.clear')" @click="filterValues = {}" />
				<q-btn v-close-popup color="primary" :label="t('common.close')" />
			</template>
		</GenericModal>

		<FireEventDialog v-model="fireOpen" :device-id="fireDeviceId" />
	</q-page>
</template>

<script setup lang="ts">
/** TYPE IMPORTS */
import type { MessageDirection } from '@stores/messages';
import type { ProtocolId } from '@services/sim';
import type { FilterValues } from '@utils/message-filters';

/** VUE IMPORTS */
import { computed, onMounted, onUnmounted, ref } from 'vue';

/** COMPONENTS */
import { AppTooltip } from '@components/AppTooltip';
import { FireEventDialog } from '@components/FireEventDialog';
import { GenericDrawer } from '@components/GenericDrawer';
import { GenericModal } from '@components/GenericModal';
import { PageHeader } from '@components/PageHeader';
import { protocolIcon } from '@components/protocols/ProtocolRegistry';
import { MessageDetail } from '@components/MessageDetail';
import { MessageFilterBar } from './components/MessageFilterBar';

/** COMPOSABLES */
import { useTranslations } from '@composables/i18n';

/** UTILS */
import { applyMessageFilters, getMessageFilterFields } from '@utils/message-filters';
import { statusColor } from '@utils/status-color';

/** STORES */
import { useAppStore } from '@stores/app';
import { useDevicesStore } from '@stores/devices';
import { useGatewaysStore } from '@stores/gateways';
import { useMessagesStore } from '@stores/messages';

/** COMPOSABLES & STORES */
const { t } = useTranslations();
const appStore = useAppStore();
const devicesStore = useDevicesStore();
const gatewaysStore = useGatewaysStore();
const messagesStore = useMessagesStore();

/** STATE */
const search = ref('');
const fireOpen = ref(false);
const fireDeviceId = ref<string | null>(null);
const detailOpen = ref(false);
const filterOpen = ref(false);
const filterValues = ref<FilterValues>({});
const deviceProtocolFilter = ref<ProtocolId | null>(null);

/** COMPUTED */

/**
 * Left-pane device list: configured devices unioned with any device seen in the
 * message stream, filtered by the search box.
 *
 * Keyed by the user-facing `deviceId` — that is the only identifier the stream
 * carries and the value the console filters on. Keying by the catalog UUID
 * instead duplicated every device the moment it produced its first frame (the
 * stream entry never matched the UUID-keyed catalog entry). `uuid` keeps the
 * catalog id for firing; it is null for a device seen only in the stream.
 */
const deviceList = computed(() => {
	const map = new Map<string, { deviceId: string; uuid: string | null; name: string; protocol: ProtocolId }>();
	for (const device of devicesStore.items) {
		map.set(device.deviceId, { deviceId: device.deviceId, uuid: device.id, name: device.name, protocol: device.protocolId });
	}
	for (const message of messagesStore.items) {
		if (!map.has(message.deviceId)) {
			map.set(message.deviceId, { deviceId: message.deviceId, uuid: null, name: message.deviceName, protocol: message.protocol });
		}
	}

	const term = search.value.trim().toLowerCase();
	let list = [...map.values()];
	if (deviceProtocolFilter.value) list = list.filter((device) => device.protocol === deviceProtocolFilter.value);
	if (term) list = list.filter((device) => device.name.toLowerCase().includes(term));
	return list;
});

const protocolFilterOptions = computed(() =>
	(['http', 'mqtt', 'lorawan', 'basicstation'] as ProtocolId[]).map((protocol) => ({
		label: t(`protocol.${protocol}`),
		value: protocol,
	})),
);

const activeProtocol = computed<ProtocolId | null>(() => {
	const id = messagesStore.deviceFilter;
	if (id) return deviceList.value.find((device) => device.deviceId === id)?.protocol ?? null;
	const selected = filterValues.value.protocol;
	return selected ? (selected as ProtocolId) : null;
});

const filterFields = computed(() => getMessageFilterFields(activeProtocol.value));

const ordered = computed(() => {
	const filtered = applyMessageFilters(messagesStore.filtered, filterFields.value, filterValues.value);
	return [...filtered].reverse();
});

const activeFilterCount = computed(() => Object.values(filterValues.value).filter((value) => value.trim() !== '').length);

// Online = enabled, since an enabled device/gateway holds a live session in the engine.
const onlineDevices = computed(() => devicesStore.items.filter((device) => device.enabled).length);
const onlineGateways = computed(() => gatewaysStore.items.filter((gateway) => gateway.enabled).length);

/** FUNCTIONS */

/**
 * Zero-pad a count to two digits for the header stat badges.
 * @param {number} value - the count
 * @returns {string} the padded count
 */
function pad(value: number): string {
	return String(value).padStart(2, '0');
}

/**
 * Icon for a message direction.
 * @param {MessageDirection} direction - the message direction
 */
function dirIcon(direction: MessageDirection): string {
	if (direction === 'up') return 'mdi-arrow-up';
	if (direction === 'down') return 'mdi-arrow-down';
	return 'mdi-cog-outline';
}

/**
 * Color for a message direction.
 * @param {MessageDirection} direction - the message direction
 */
function dirColor(direction: MessageDirection): string {
	if (direction === 'up') return 'teal';
	if (direction === 'down') return 'primary';
	return 'grey';
}

/**
 * Select a message and open the detail drawer on it.
 * @param {string} id - the message id to inspect
 */
function selectMessage(id: string): void {
	messagesStore.select(id);
	detailOpen.value = true;
}

/**
 * Clear the selection when the detail drawer closes, so the row highlight drops
 * and reopening the same row works cleanly.
 */
function onDetailClose(): void {
	messagesStore.select(null);
}

/**
 * Open the fire-event dialog for a specific device.
 * @param {string} deviceId - the device to fire from
 */
function openFire(deviceId: string): void {
	fireDeviceId.value = deviceId;
	fireOpen.value = true;
}

/** LIFECYCLE HOOKS */
onMounted(() => {
	void appStore.startHealthPolling();
	void devicesStore.fetch();
	void gatewaysStore.fetch();
	messagesStore.connect();
});

onUnmounted(() => {
	messagesStore.disconnect();
});
</script>

<style scoped lang="scss">
.console-page {
	padding: var(--mapex-spacing-md);
	display: flex;
	flex-direction: column;
	height: calc(100vh - 116px);
	min-height: 540px;
}

.console {
	display: flex;
	flex: 1;
	min-height: 0;
	border: 1px solid var(--mapex-card-border);
	border-radius: var(--mapex-radius-md);
	overflow: hidden;
	background: var(--mapex-surface-bg);
}

.console__pane-head {
	// Fixed bar height so all three pane headers line up, even though the middle
	// one carries the filter/clear buttons and the others are plain text. Keep
	// this in sync with the message-detail header height.
	display: flex;
	align-items: center;
	min-height: 52px;
	padding: 0 var(--mapex-spacing-lg);
	border-bottom: 1px solid var(--mapex-divider);
	font-size: var(--mapex-font-sm);
	font-weight: var(--mapex-font-weight-semibold);
	color: var(--mapex-text-primary);
}

.console__stat {
	font-family: var(--mapex-mono-font);
	font-weight: var(--mapex-font-weight-semibold);
	padding: var(--mapex-spacing-2xs) var(--mapex-spacing-xs);
}

.console__devices {
	width: 264px;
	flex-shrink: 0;
	display: flex;
	flex-direction: column;
	border-right: 1px solid var(--mapex-divider);
	background: var(--mapex-surface-elevated);

	.console__search {
		padding: var(--mapex-spacing-sm) var(--mapex-spacing-md);
	}

	.console__device-list {
		flex: 1;
		overflow-y: auto;
	}

	.console__device--active {
		background: var(--mapex-active-bg);
		color: var(--mapex-primary);
	}

	.console__devices-empty {
		padding: var(--mapex-spacing-lg);
		text-align: center;
		color: var(--mapex-text-muted);
		font-size: var(--mapex-font-sm);
	}

	.console__devices-foot {
		border-top: 1px solid var(--mapex-divider);
		padding: var(--mapex-spacing-xs);
	}
}

.console__log {
	flex: 1;
	min-width: 0;
	display: flex;
	flex-direction: column;

	.console__empty {
		flex: 1;
		display: flex;
		align-items: center;
		justify-content: center;
		color: var(--mapex-text-muted);
		padding: var(--mapex-spacing-2xl);
		text-align: center;
	}

	.console__rows {
		flex: 1;
		overflow-y: auto;
	}
}

.msg {
	display: flex;
	align-items: center;
	gap: var(--mapex-spacing-sm);
	width: 100%;
	text-align: left;
	border: none;
	border-bottom: 1px solid var(--mapex-divider);
	background: transparent;
	padding: var(--mapex-spacing-sm) var(--mapex-spacing-md);
	cursor: pointer;
	font-family: var(--mapex-mono-font);
	font-size: var(--mapex-font-xs);
	color: var(--mapex-text-primary);
	transition: var(--mapex-transition-fast);

	&:hover {
		background: var(--mapex-surface-highlight);
	}

	&--active {
		background: var(--mapex-active-bg);
	}

	&__dir {
		flex-shrink: 0;
	}

	&__ts {
		color: var(--mapex-text-muted);
		flex-shrink: 0;
	}

	&__proto {
		flex-shrink: 0;
	}

	&__device {
		color: var(--mapex-primary);
		flex-shrink: 0;
		max-width: 140px;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	&__summary {
		flex: 1;
		min-width: 0;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
		color: var(--mapex-text-secondary);
	}

	&__status {
		flex-shrink: 0;
	}
}
</style>
