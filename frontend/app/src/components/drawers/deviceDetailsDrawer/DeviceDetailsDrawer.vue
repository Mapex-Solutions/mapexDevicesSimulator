<template>
	<!-- Invisible backdrop for click-outside detection (matches MapexOS). -->
	<Teleport to="body">
		<div v-if="modelValue" class="drawer-backdrop" @click="close" />
	</Teleport>

	<q-drawer
		overlay
		bordered
		side="right"
		:model-value="modelValue"
		:width="450"
		@update:model-value="emit('update:modelValue', $event)"
	>
		<!-- Header -->
		<q-toolbar class="drawer-header">
			<q-icon name="memory" size="sm" class="q-mr-sm" color="primary" />
			<q-toolbar-title class="text-weight-medium">{{ t('devices.drawer.title') }}</q-toolbar-title>
			<q-btn flat round dense icon="close" class="drawer-close-btn" @click="close">
				<q-tooltip>{{ t('devices.drawer.close') }}</q-tooltip>
			</q-btn>
		</q-toolbar>

		<q-separator />

		<!-- Content -->
		<div class="drawer-content">
			<q-scroll-area class="fit">
				<div v-if="device" class="q-px-md q-py-lg">
					<!-- Basic Information -->
					<div class="section q-mb-md">
						<div class="section-header">
							<q-icon name="info" color="primary" size="sm" class="q-mr-sm" />
							<span class="text-subtitle1 text-weight-medium">{{ t('devices.drawer.sections.basic') }}</span>
						</div>
						<q-separator class="q-my-sm" />

						<div class="field-row">
							<div class="field-label">{{ t('devices.col.name') }}</div>
							<div class="field-value text-weight-medium">{{ device.name || '-' }}</div>
						</div>

						<div class="field-row">
							<div class="field-label">{{ t('devices.fields.deviceId') }}</div>
							<div class="field-value">
								<q-chip dense square size="sm" icon="fingerprint" class="value-chip">{{ device.deviceId || '-' }}</q-chip>
							</div>
						</div>

						<div class="row q-col-gutter-sm">
							<div class="col-7">
								<div class="field-row">
									<div class="field-label">{{ t('devices.col.protocol') }}</div>
									<div class="field-value">
										<q-chip dense square size="sm" color="primary" text-color="white" :icon="protocolIcon(device.protocolId)">
											{{ t(`protocol.${device.protocolId}`) }}
										</q-chip>
									</div>
								</div>
							</div>
							<div class="col-5">
								<div class="field-row">
									<div class="field-label">{{ t('devices.col.status') }}</div>
									<div class="field-value">
										<q-chip dense square size="sm" :color="device.enabled ? 'positive' : 'grey'" text-color="white">
											{{ device.enabled ? t('devices.on') : t('devices.off') }}
										</q-chip>
									</div>
								</div>
							</div>
						</div>
					</div>

					<!-- Connection -->
					<div class="section q-mb-md">
						<div class="section-header">
							<q-icon name="cable" color="primary" size="sm" class="q-mr-sm" />
							<span class="text-subtitle1 text-weight-medium">{{ t('devices.drawer.sections.connection') }}</span>
						</div>
						<q-separator class="q-my-sm" />

						<div v-for="field in connectionFields" :key="field.label" class="field-row">
							<div class="field-label">{{ field.label }}</div>
							<div class="field-value" :class="{ 'text-grey-6': !field.value }">{{ field.value || '-' }}</div>
						</div>
					</div>

					<!-- Attributes -->
					<div class="section q-mb-md">
						<div class="section-header">
							<q-icon name="sell" color="primary" size="sm" class="q-mr-sm" />
							<span class="text-subtitle1 text-weight-medium">{{ t('devices.attributes') }}</span>
						</div>
						<q-separator class="q-my-sm" />

						<div v-if="attributeRows.length" class="row q-gutter-xs">
							<q-chip v-for="attr in attributeRows" :key="attr.key" dense square size="sm" class="value-chip">
								{{ attr.key }}: {{ attr.value }}
							</q-chip>
						</div>
						<div v-else class="text-grey-6 text-caption">{{ t('devices.drawer.noAttributes') }}</div>
					</div>

					<!-- Meta -->
					<div class="section">
						<div class="section-header">
							<q-icon name="schedule" color="primary" size="sm" class="q-mr-sm" />
							<span class="text-subtitle1 text-weight-medium">{{ t('devices.drawer.sections.meta') }}</span>
						</div>
						<q-separator class="q-my-sm" />

						<div class="row q-col-gutter-sm">
							<div class="col-6">
								<div class="field-row">
									<div class="field-label">{{ t('devices.col.events') }}</div>
									<div class="field-value">{{ device.events?.length ?? 0 }}</div>
								</div>
							</div>
							<div class="col-6">
								<div class="field-row">
									<div class="field-label">{{ t('devices.storeLogs') }}</div>
									<div class="field-value">{{ device.storeLogs ? t('devices.on') : t('devices.off') }}</div>
								</div>
							</div>
						</div>

						<div class="field-row">
							<div class="field-label">{{ t('devices.col.created') }}</div>
							<div class="field-value">{{ device.created ? formatDate(device.created) : '-' }}</div>
						</div>
					</div>
				</div>
			</q-scroll-area>
		</div>

		<!-- Footer Actions -->
		<q-separator />
		<div class="drawer-footer">
			<q-btn flat no-caps :label="t('devices.drawer.cancel')" @click="close" />
			<q-space />
			<q-btn unelevated no-caps icon="edit" color="primary" :label="t('devices.drawer.edit')" :disable="!device" @click="handleEdit" />
		</div>
	</q-drawer>
</template>

<script setup lang="ts">
defineOptions({ name: 'DeviceDetailsDrawer' });

/** TYPE IMPORTS */
import type { Device } from '@services/sim';

/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPONENTS */
import { protocolIcon } from '@components/protocols/ProtocolRegistry';

/** COMPOSABLES */
import { useTranslations } from '@composables/i18n';
import { useDateFormat } from '@composables/datetime';

/** STORES */
import { useGatewaysStore } from '@stores/gateways';

/** PROPS & EMITS */
const props = defineProps<{ modelValue: boolean; device: Device | null }>();
const emit = defineEmits<{
	'update:modelValue': [value: boolean];
	edit: [device: Device];
}>();

/** COMPOSABLES & STORES */
const { t } = useTranslations();
const { formatDate } = useDateFormat();
const gatewaysStore = useGatewaysStore();

/**
 * Resolve a LoRaWAN gateway's display name from its id (via the gateways store),
 * falling back to the raw id when the gateway is not loaded.
 * @param {string} id - the gateway id stored on the device config
 * @returns {string} the gateway name, or the id, or a dash
 */
function gatewayName(id: string): string {
	if (!id) return '';
	return gatewaysStore.items.find((gateway) => gateway.id === id)?.name ?? id;
}

/** COMPUTED */

/** The protocol-specific connection fields shown in the Connection section. */
const connectionFields = computed<Array<{ label: string; value: string }>>(() => {
	const device = props.device;
	if (!device) return [];
	const config = device.config;
	if (config.kind === 'http') {
		return [
			{ label: t('connections.http.url'), value: config.url },
			{ label: t('connections.http.method'), value: config.method },
			{ label: t('connections.http.authMode'), value: config.authMode },
		];
	}
	if (config.kind === 'mqtt') {
		return [
			{ label: t('connections.mqtt.brokerUrl'), value: config.brokerUrl },
			{ label: t('connections.mqtt.clientId'), value: config.clientId },
			{ label: t('connections.mqtt.baseTopic'), value: config.baseTopic },
			{ label: t('connections.mqtt.authMode'), value: config.authMode },
		];
	}
	if (config.kind === 'lorawan') {
		return [
			{ label: t('connections.lorawan.gateway'), value: gatewayName(config.gatewayId) },
			{ label: t('connections.lorawan.region'), value: config.region },
			{ label: t('connections.lorawan.macVersion'), value: config.macVersion },
			{ label: t('connections.lorawan.deviceClass'), value: config.class },
			{ label: t('connections.lorawan.activation'), value: config.activation.toUpperCase() },
			{ label: t('connections.lorawan.devEui'), value: config.devEui },
		];
	}
	// basicstation
	return [
		{ label: t('connections.basicstation.lnsUri'), value: config.lnsUri },
		{ label: t('connections.basicstation.gatewayEui'), value: config.gatewayEui },
		{ label: t('connections.lorawan.region'), value: config.region },
		{ label: t('connections.lorawan.activation'), value: config.activation.toUpperCase() },
		{ label: t('connections.lorawan.devEui'), value: config.devEui },
	];
});

/** Device attributes as a sorted key/value list. */
const attributeRows = computed(() =>
	Object.entries(props.device?.attributes ?? {}).map(([key, value]) => ({ key, value })),
);

/** FUNCTIONS */

/** Close the drawer. */
function close(): void {
	emit('update:modelValue', false);
}

/** Emit edit for the current device and close. */
function handleEdit(): void {
	if (!props.device) return;
	emit('edit', props.device);
	close();
}
</script>

<style lang="scss" scoped>
:deep(.q-drawer__content) {
	display: flex;
	flex-direction: column;
	height: 100%;
}

.drawer-header {
	flex-shrink: 0;
	background: var(--mapex-header-bg);
	border-bottom: 1px solid var(--mapex-header-border);

	.q-toolbar__title {
		font-size: 1.1rem;
		color: var(--q-primary);
	}
}

.drawer-close-btn {
	color: var(--mapex-text-secondary);
}

:global(.drawer-backdrop) {
	position: fixed;
	top: 0;
	left: 0;
	right: 450px;
	bottom: 0;
	background: transparent;
	z-index: 5999;
	cursor: default;
}

.drawer-content {
	flex: 1;
	min-height: 0;
	overflow: hidden;

	:deep(.q-scrollarea__content) {
		width: 100%;
		max-width: 100%;
		overflow-x: hidden;
	}
}

.drawer-footer {
	flex-shrink: 0;
	display: flex;
	align-items: center;
	padding: 12px 16px;
	background: var(--mapex-header-bg);
	border-top: 1px solid var(--mapex-header-border);
	box-shadow: 0 -2px 8px var(--mapex-elevation-shadow);
}

.section {
	.section-header {
		display: flex;
		align-items: center;
		color: var(--q-primary);
		margin-bottom: 8px;
	}
}

.field-row {
	display: flex;
	flex-direction: column;
	padding: 10px 0;
	border-bottom: 1px solid var(--mapex-divider);

	&:last-child {
		border-bottom: none;
	}

	.field-label {
		font-size: 0.7rem;
		font-weight: 600;
		text-transform: uppercase;
		color: var(--mapex-text-secondary);
		margin-bottom: 4px;
		letter-spacing: 0.8px;
	}

	.field-value {
		font-size: 0.9rem;
		color: var(--mapex-text-primary);
		word-break: break-word;
		line-height: 1.4;
	}
}

.value-chip {
	max-width: 100%;
}
</style>
