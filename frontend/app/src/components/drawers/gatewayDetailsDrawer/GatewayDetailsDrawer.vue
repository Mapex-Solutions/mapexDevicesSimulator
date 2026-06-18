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
			<q-icon name="router" size="sm" class="q-mr-sm" color="primary" />
			<q-toolbar-title class="text-weight-medium">{{ t('gateways.drawer.title') }}</q-toolbar-title>
			<q-btn flat round dense icon="close" class="drawer-close-btn" @click="close">
				<q-tooltip>{{ t('gateways.drawer.close') }}</q-tooltip>
			</q-btn>
		</q-toolbar>

		<q-separator />

		<!-- Content -->
		<div class="drawer-content">
			<q-scroll-area class="fit">
				<div v-if="gateway" class="q-px-md q-py-lg">
					<!-- Basic Information -->
					<div class="section q-mb-md">
						<div class="section-header">
							<q-icon name="info" color="primary" size="sm" class="q-mr-sm" />
							<span class="text-subtitle1 text-weight-medium">{{ t('gateways.drawer.sections.basic') }}</span>
						</div>
						<q-separator class="q-my-sm" />

						<div class="field-row">
							<div class="field-label">{{ t('gateways.fields.name') }}</div>
							<div class="field-value text-weight-medium">{{ gateway.name || '-' }}</div>
						</div>

						<div class="field-row">
							<div class="field-label">{{ t('gateways.fields.eui') }}</div>
							<div class="field-value">
								<q-chip v-if="gateway.eui" dense square size="sm" icon="fingerprint" class="value-chip">{{ gateway.eui }}</q-chip>
								<span v-else class="text-grey-6">{{ t('gateways.noEui') }}</span>
							</div>
						</div>

						<div class="field-row">
							<div class="field-label">{{ t('gateways.fields.status') }}</div>
							<div class="field-value">
								<q-chip dense square size="sm" :color="gateway.enabled ? 'positive' : 'grey'" text-color="white">
									{{ gateway.enabled ? t('devices.on') : t('devices.off') }}
								</q-chip>
							</div>
						</div>
					</div>

					<!-- Link -->
					<div class="section q-mb-md">
						<div class="section-header">
							<q-icon name="cable" color="primary" size="sm" class="q-mr-sm" />
							<span class="text-subtitle1 text-weight-medium">{{ t('gateways.drawer.sections.link') }}</span>
						</div>
						<q-separator class="q-my-sm" />

						<div class="row q-col-gutter-sm">
							<div class="col-7">
								<div class="field-row">
									<div class="field-label">{{ t('gateways.col.connection') }}</div>
									<div class="field-value">
										<q-chip dense square size="sm" color="primary" text-color="white">{{ linkProtocolLabel }}</q-chip>
									</div>
								</div>
							</div>
							<div class="col-5">
								<div class="field-row">
									<div class="field-label">{{ t('gateways.fields.region') }}</div>
									<div class="field-value">{{ gateway.region }}</div>
								</div>
							</div>
						</div>

						<div class="field-row">
							<div class="field-label">{{ t('gateways.drawer.endpoint') }}</div>
							<div class="field-value">{{ endpoint || '-' }}</div>
						</div>
					</div>

					<!-- Description -->
					<div class="section q-mb-md">
						<div class="section-header">
							<q-icon name="notes" color="primary" size="sm" class="q-mr-sm" />
							<span class="text-subtitle1 text-weight-medium">{{ t('gateways.fields.description') }}</span>
						</div>
						<q-separator class="q-my-sm" />
						<div class="field-value" :class="{ 'text-grey-6': !gateway.description }">
							{{ gateway.description || t('gateways.drawer.noDescription') }}
						</div>
					</div>

					<!-- Meta -->
					<div class="section">
						<div class="section-header">
							<q-icon name="schedule" color="primary" size="sm" class="q-mr-sm" />
							<span class="text-subtitle1 text-weight-medium">{{ t('gateways.drawer.sections.meta') }}</span>
						</div>
						<q-separator class="q-my-sm" />
						<div class="field-row">
							<div class="field-label">{{ t('gateways.col.created') }}</div>
							<div class="field-value">{{ gateway.created ? formatDate(gateway.created) : '-' }}</div>
						</div>
					</div>
				</div>
			</q-scroll-area>
		</div>

		<!-- Footer Actions -->
		<q-separator />
		<div class="drawer-footer">
			<q-btn flat no-caps :label="t('gateways.drawer.cancel')" @click="close" />
			<q-space />
			<q-btn unelevated no-caps icon="edit" color="primary" :label="t('gateways.drawer.edit')" :disable="!gateway" @click="handleEdit" />
		</div>
	</q-drawer>
</template>

<script setup lang="ts">
defineOptions({ name: 'GatewayDetailsDrawer' });

/** TYPE IMPORTS */
import type { Gateway } from '@services/sim';

/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPOSABLES */
import { useTranslations } from '@composables/i18n';
import { useDateFormat } from '@composables/datetime';

/** PROPS & EMITS */
const props = defineProps<{ modelValue: boolean; gateway: Gateway | null }>();
const emit = defineEmits<{
	'update:modelValue': [value: boolean];
	edit: [gateway: Gateway];
}>();

/** COMPOSABLES */
const { t } = useTranslations();
const { formatDate } = useDateFormat();

/** COMPUTED */

/** Human label for the gateway's LNS link protocol. */
const linkProtocolLabel = computed(() =>
	props.gateway?.link.protocol === 'basicstation' ? t('protocol.basicstation') : 'Semtech UDP',
);

/** The LNS endpoint: the WebSocket URI for Basics Station, or host:port for UDP. */
const endpoint = computed(() => {
	const link = props.gateway?.link;
	if (!link) return '';
	return link.protocol === 'basicstation' ? link.lnsUri : `${link.host}:${link.port}`;
});

/** FUNCTIONS */

/** Close the drawer. */
function close(): void {
	emit('update:modelValue', false);
}

/** Emit edit for the current gateway and close. */
function handleEdit(): void {
	if (!props.gateway) return;
	emit('edit', props.gateway);
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
