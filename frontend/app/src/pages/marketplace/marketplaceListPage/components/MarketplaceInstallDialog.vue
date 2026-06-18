<template>
	<q-dialog :model-value="modelValue" @update:model-value="emit('update:modelValue', $event)">
		<q-card class="install-dialog">
			<!-- Header -->
			<q-card-section class="row items-center q-pb-sm">
				<q-icon name="add_box" color="primary" size="sm" class="q-mr-sm" />
				<div>
					<div class="text-h6">{{ t('marketplace.install.title') }}</div>
					<div class="text-caption text-grey-7">{{ item?.name }}</div>
				</div>
				<q-space />
				<q-btn v-close-popup flat round dense icon="close" />
			</q-card-section>

			<q-separator />

			<!-- Loading the template -->
			<q-card-section v-if="loading || !form" class="column items-center q-py-xl">
				<q-spinner color="primary" size="2.5em" class="q-mb-sm" />
				<div class="text-grey-7">{{ t('marketplace.install.loading') }}</div>
			</q-card-section>

			<!-- Editable draft -->
			<q-card-section v-else class="install-body q-gutter-md">
				<!-- Name -->
				<q-input
					v-model="form.name"
					dense
					outlined
					:label="t('devices.fields.name')"
				/>

				<!-- Device ID with random button -->
				<q-input
					v-model="form.deviceId"
					dense
					outlined
					:label="t('devices.fields.deviceId')"
					:placeholder="t('devices.fields.deviceIdPlaceholder')"
				>
					<template #append>
						<q-btn flat dense round icon="mdi-refresh" size="sm" @click="regenerateDeviceId">
							<q-tooltip>{{ t('devices.fields.deviceIdRegenerate') }}</q-tooltip>
						</q-btn>
					</template>
				</q-input>

				<q-separator />

				<!-- Protocol-specific connection config (reused from the device wizard) -->
				<div class="text-subtitle2 text-primary">{{ t('marketplace.install.connection') }}</div>
				<component :is="configComponent" v-if="configComponent" v-model="form.config" />
			</q-card-section>

			<q-separator />

			<!-- Footer -->
			<q-card-actions class="q-pa-md">
				<q-btn v-close-popup flat no-caps :label="t('marketplace.install.cancel')" />
				<q-space />
				<q-btn
					unelevated
					no-caps
					color="primary"
					icon="add"
					:label="t('marketplace.install.confirm')"
					:disable="!canInstall"
					:loading="installing"
					@click="confirm"
				/>
			</q-card-actions>
		</q-card>
	</q-dialog>
</template>

<script setup lang="ts">
defineOptions({ name: 'MarketplaceInstallDialog' });

/** TYPE IMPORTS */
import type { DeviceInput } from '@services/sim';
import type { MarketplaceCardItem } from '../interfaces/marketplaceListPage.interface';

/** VUE IMPORTS */
import { computed, ref, watch } from 'vue';

/** COMPONENTS */
import { PROTOCOL_REGISTRY } from '@components/protocols/ProtocolRegistry';

/** COMPOSABLES */
import { useTranslations } from '@composables/i18n';

/** STORES */
import { useGatewaysStore } from '@stores/gateways';

/** PROPS & EMITS */
const props = defineProps<{
	modelValue: boolean;
	item: MarketplaceCardItem | null;
	/** The prepared device draft (template + minted deviceId), or null while loading. */
	draft: DeviceInput | null;
	loading: boolean;
	installing: boolean;
}>();
const emit = defineEmits<{
	'update:modelValue': [value: boolean];
	confirm: [input: DeviceInput];
}>();

/** COMPOSABLES & STORES */
const { t } = useTranslations();
const gatewaysStore = useGatewaysStore();

/** STATE — a local, editable copy of the draft so edits don't mutate the prop. */
const form = ref<DeviceInput | null>(null);

/** COMPUTED */

/** The protocol's connection-config component (reused from the device wizard). */
const configComponent = computed(() =>
	form.value ? PROTOCOL_REGISTRY[form.value.protocolId]?.configComponent : undefined,
);

/** Whether the draft is complete enough to create (identity + a valid config). */
const canInstall = computed(() => {
	if (!form.value) return false;
	if (!form.value.name.trim() || !form.value.deviceId.trim()) return false;
	return PROTOCOL_REGISTRY[form.value.protocolId]?.validate(form.value.config).valid ?? false;
});

/** WATCHERS */

/** Copy the incoming draft into the local editable form whenever it arrives. */
watch(
	() => props.draft,
	(draft) => {
		form.value = draft ? (JSON.parse(JSON.stringify(draft)) as DeviceInput) : null;
	},
	{ immediate: true },
);

/** Load gateways when the dialog opens so the LoRaWAN gateway picker has options. */
watch(
	() => props.modelValue,
	(open) => {
		if (open && !gatewaysStore.items.length) void gatewaysStore.fetch();
	},
);

/** FUNCTIONS */

/**
 * Mint a fresh, readable deviceId from the model name, e.g. `em300-th-a1b2c3d4`.
 * Mirrors the store's default so the random button stays consistent.
 */
function regenerateDeviceId(): void {
	if (!form.value) return;
	const base = (props.item?.model ?? 'device').toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-|-$/g, '');
	const suffix =
		typeof crypto !== 'undefined' && 'randomUUID' in crypto
			? crypto.randomUUID().slice(0, 8)
			: Math.floor(Math.random() * 0xffffffff).toString(16).padStart(8, '0');
	form.value.deviceId = `${base}-${suffix}`;
}

/** Emit the finalized draft for the page to create. */
function confirm(): void {
	if (!form.value || !canInstall.value) return;
	emit('confirm', JSON.parse(JSON.stringify(form.value)) as DeviceInput);
}
</script>

<style scoped lang="scss">
.install-dialog {
	width: 560px;
	max-width: 95vw;
}

.install-body {
	max-height: 65vh;
	overflow-y: auto;
}
</style>
