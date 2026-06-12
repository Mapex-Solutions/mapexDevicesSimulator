<template>
	<q-page class="q-pa-lg">
		<div v-if="loading" class="flex flex-center q-pa-xl">
			<q-spinner color="primary" size="42px" />
		</div>

		<div v-else>
			<PageHeader
				icon="memory"
				icon-color="primary"
				:title="pageTitle"
				:description="t('createDevice.description')"
				:button="{ label: t('createDevice.back'), icon: 'arrow_back', flat: true, onClick: goBack }"
			/>

			<div class="row q-col-gutter-lg">
				<!-- Stepper (left) -->
				<div class="col-12 col-md-4">
					<StepperVertical
						:title="t('createDevice.stepperTitle')"
						:subtitle="t('createDevice.stepperSubtitle')"
						:info-text="t('createDevice.stepperInfo')"
						:current-step-label="t('createDevice.currentStep')"
						:current-step="currentStep"
						:steps="steps"
						:allow-step-navigation="isEditMode"
						@step-click="changeStep"
					/>
				</div>

				<!-- Body (right) -->
				<FormCard
					:header="currentHeader"
					:navigation="navigation"
					:button-labels="buttonLabels"
					@previous="changeStep"
					@next="onNext"
					@save="onSave"
				>
					<template #form>
						<!-- Step 1: Identity -->
						<div v-if="currentStep === 1" class="row q-col-gutter-md">
							<q-input
								class="col-12"
								v-model="draft.name"
								:label="`${t('devices.fields.name')} *`"
								outlined
								dense
								lazy-rules
								:rules="[requiredRule]"
								@update:model-value="(val) => (draft.name = String(val ?? ''))"
								hide-bottom-space
							/>
							<q-select
								class="col-12 col-sm-6"
								v-model="draft.protocolId"
								:options="protocolOptions"
								:label="t('devices.protocol')"
								outlined
								dense
								emit-value
								map-options
								@update:model-value="onProtocolChange"
								hide-bottom-space
							/>
							<q-input
								class="col-12 col-sm-6"
								v-model="draft.deviceId"
								:label="t('devices.fields.deviceId')"
								:placeholder="t('devices.fields.deviceIdPlaceholder')"
								outlined
								dense
								hide-bottom-space
							>
								<template #append>
									<q-btn flat dense round icon="mdi-refresh" size="sm" @click="regenerateDeviceId">
										<AppTooltip>{{ t('devices.fields.deviceIdRegenerate') }}</AppTooltip>
									</q-btn>
								</template>
							</q-input>
							<q-select
								class="col-12 col-sm-6"
								v-model="draft.enabled"
								:options="boolOptions"
								:label="t('devices.status')"
								outlined
								dense
								emit-value
								map-options
								hide-bottom-space
							/>
							<q-select
								class="col-12 col-sm-6"
								v-model="draft.storeLogs"
								:options="boolOptions"
								:label="t('devices.storeLogs')"
								outlined
								dense
								emit-value
								map-options
								hide-bottom-space
							/>
						</div>

						<!-- Step 2: Target -->
						<div v-else-if="currentStep === 2">
							<component :is="activeDef.configComponent" v-if="activeDef" v-model="draft.config" />
							<q-banner v-else dense class="protocol-soon">
								<template #avatar><q-icon name="mdi-information-outline" color="primary" /></template>
								{{ t('createDevice.protocolSoon') }}
							</q-banner>
						</div>

						<!-- Step 3: Events -->
						<div v-else-if="currentStep === 3" class="column q-gutter-md">
							<div class="step-hint">{{ t('deviceEvents.intro', { protocol: t(`protocol.${draft.protocolId}`) }) }}</div>

							<!-- Events unsupported for this protocol -->
							<q-banner v-if="!eventsSupported" dense class="protocol-soon">
								<template #avatar><q-icon name="mdi-information-outline" color="primary" /></template>
								{{ t('deviceEvents.unsupported', { protocol: t(`protocol.${draft.protocolId}`) }) }}
							</q-banner>

							<DeviceEventsEditor
								v-else
								:model-value="draft.events"
								:protocol-id="draft.protocolId"
								@update:model-value="(v) => (draft.events = v)"
							/>
						</div>
					</template>
				</FormCard>
			</div>
		</div>
	</q-page>
</template>

<script setup lang="ts">
/** TYPE IMPORTS */
import type { Device, DeviceInput, ProtocolConfig, ProtocolId } from '@services/sim';
import type { FormCardHeader, FormCardNavigation } from '@components/FormCard';
import type { StepperVerticalItem } from '@components/StepperVertical';

/** VUE IMPORTS */
import { computed, onMounted, reactive, ref } from 'vue';

/** COMPONENTS */
import { PageHeader } from '@components/PageHeader';
import { FormCard } from '@components/FormCard';
import { StepperVertical } from '@components/StepperVertical';
import { ENABLED_PROTOCOLS, PROTOCOL_REGISTRY } from '@components/protocols/ProtocolRegistry';
import { AppTooltip } from '@components/AppTooltip';
import { DeviceEventsEditor } from './components/DeviceEventsEditor';

/** COMPOSABLES */
import { useTranslations } from '@composables/i18n';
import { useStepperNavigation } from '@composables/shared/form';

/** UTILS */
import { useQuasar } from 'quasar';

/** SERVICES */
import { useRoute, useRouter } from 'vue-router';

/** STORES */
import { useDevicesStore } from '@stores/devices';

const TOTAL_STEPS = 3;

/** COMPOSABLES & STORES */
const { t } = useTranslations();
const $q = useQuasar();
const route = useRoute();
const router = useRouter();
const devicesStore = useDevicesStore();

/** STATE */
const editingId = ref<string | null>((route.params.id as string | undefined) ?? null);
const isEditMode = computed(() => editingId.value !== null);
const loading = ref(false);
const saving = ref(false);
const currentStep = ref(1);
const draft = reactive<DeviceInput>(blankDraft());

/** COMPUTED */
const protocolOptions = computed(() =>
	(['http', 'mqtt', 'lorawan', 'basicstation'] as ProtocolId[]).map((p) => ({ label: t(`protocol.${p}`), value: p })),
);
const activeDef = computed(() => PROTOCOL_REGISTRY[draft.protocolId]);
const eventsSupported = computed(
	() =>
		draft.protocolId === 'http' ||
		draft.protocolId === 'mqtt' ||
		draft.protocolId === 'lorawan' ||
		draft.protocolId === 'basicstation',
);

const boolOptions = computed(() => [
	{ label: t('devices.on'), value: true },
	{ label: t('devices.off'), value: false },
]);

const steps = computed<StepperVerticalItem[]>(() => [
	{ title: t('createDevice.steps.identity'), description: t('createDevice.steps.identityDesc'), icon: 'badge' },
	{ title: t('createDevice.steps.target'), description: t('createDevice.steps.targetDesc'), icon: 'cloud_upload' },
	{ title: t('createDevice.steps.events'), description: t('createDevice.steps.eventsDesc'), icon: 'bolt' },
]);

const currentHeader = computed<FormCardHeader>(() => {
	const step = steps.value[currentStep.value - 1];
	return { icon: step?.icon ?? 'badge', title: step?.title ?? '', description: step?.description ?? '' };
});

const navigation = computed<FormCardNavigation>(() => ({
	currentStep: currentStep.value,
	totalSteps: TOTAL_STEPS,
	showPreviousButton: true,
	showNextButton: true,
	showSaveButton: true,
	loadingSaveButton: saving.value,
}));

const buttonLabels = computed(() => ({
	previous: t('common.back'),
	next: t('common.next'),
	save: t('common.save'),
}));

const pageTitle = computed(() => (isEditMode.value ? t('createDevice.editTitle') : t('createDevice.title')));

/** FUNCTIONS */

/**
 * Build a blank device draft seeded with the first enabled protocol.
 */
function blankDraft(): DeviceInput {
	const def = ENABLED_PROTOCOLS[0];
	return {
		name: '',
		deviceId: generateUuid(),
		protocolId: def?.id ?? 'http',
		enabled: true,
		storeLogs: true,
		config: def?.defaultConfig() ?? buildEmptyConfig(),
		attributes: {},
		events: [],
	};
}

/**
 * Generate a UUID for a new device, falling back when crypto is unavailable.
 * @returns {string} the generated identifier
 */
function generateUuid(): string {
	const cryptoApi = globalThis.crypto;
	if (cryptoApi && typeof cryptoApi.randomUUID === 'function') return cryptoApi.randomUUID();
	return `dev-${Math.floor(Math.random() * 1e9)}`;
}

/**
 * Replace the device identifier with a fresh UUID.
 */
function regenerateDeviceId(): void {
	draft.deviceId = generateUuid();
}

function buildEmptyConfig(): ProtocolConfig {
	return {
		kind: 'http', url: '', method: 'POST', headers: [{ key: 'Content-Type', value: 'application/json' }],
		authMode: 'none', apiKeyHeader: 'X-API-Key', apiKey: '', basicUser: '', basicPass: '',
	};
}

/**
 * Validation rule for a non-empty text field.
 * @param {string | null | undefined} value - the field value
 * @returns {true | string} true when valid, otherwise the error message
 */
function requiredRule(value: string | null | undefined): true | string {
	return (!!value && value.trim().length > 0) || t('validation.required');
}

/**
 * Validation message for a single step, or null when the step is valid.
 * @param {number} step - the step to validate
 * @returns {string | null} the error message or null
 */
function stepError(step: number): string | null {
	if (step === 1 && !draft.name.trim()) return t('validation.nameRequired');

	if (step === 2) {
		const config = draft.config;
		if (config.kind === 'http') {
			if (!config.url.trim()) return t('validation.urlRequired');
			if (config.authMode === 'apiKey' && !config.apiKey.trim()) return t('validation.credentialsRequired');
			if (config.authMode === 'basic' && (!config.basicUser.trim() || !config.basicPass.trim())) {
				return t('validation.credentialsRequired');
			}
		}
		if (config.kind === 'mqtt') {
			if (!config.brokerUrl.trim()) return t('validation.brokerRequired');
			if (config.authMode === 'userpass' && !config.username.trim()) return t('validation.credentialsRequired');
			if (config.authMode === 'tls' && !config.tlsCertPem.trim()) return t('validation.credentialsRequired');
		}
		if (config.kind === 'lorawan') {
			if (!config.gatewayId) return t('validation.gatewayRequired');
			if (config.activation === 'otaa' && (!config.devEui.trim() || !config.appKey.trim())) {
				return t('validation.keysRequired');
			}
			if (config.activation === 'abp' && (!config.devAddr.trim() || !config.nwkSKey.trim() || !config.appSKey.trim())) {
				return t('validation.keysRequired');
			}
		}
		if (config.kind === 'basicstation') {
			if (!config.lnsUri.trim() || !config.gatewayEui.trim()) return t('validation.linkRequired');
			if (config.activation === 'otaa' && (!config.devEui.trim() || !config.appKey.trim())) {
				return t('validation.keysRequired');
			}
			if (config.activation === 'abp' && (!config.devAddr.trim() || !config.nwkSKey.trim() || !config.appSKey.trim())) {
				return t('validation.keysRequired');
			}
		}
	}

	return null;
}

/**
 * First validation problem across all steps, or null when ready to save.
 * @returns {{ step: number; message: string } | null} the offending step and message
 */
function firstError(): { step: number; message: string } | null {
	for (let step = 1; step <= TOTAL_STEPS; step += 1) {
		const message = stepError(step);
		if (message) return { step, message };
	}
	return null;
}

/**
 * Advance to a step, validating the current step first when moving forward.
 * @param {number} step - the target step
 */
function onNext(step: number): void {
	if (step > currentStep.value) {
		const message = stepError(currentStep.value);
		if (message) {
			$q.notify({ type: 'negative', message });
			return;
		}
	}
	changeStep(step);
}

/**
 * Reset the config to the chosen protocol's defaults.
 * @param {ProtocolId} protocolId - the chosen protocol
 */
function onProtocolChange(protocolId: ProtocolId): void {
	const def = PROTOCOL_REGISTRY[protocolId];
	if (def) draft.config = def.defaultConfig();
}

/**
 * Move to a step, clamped to the valid range.
 * @param {number} step - the target step
 */
function changeStep(step: number): void {
	currentStep.value = Math.min(TOTAL_STEPS, Math.max(1, step));
}

function goBack(): void {
	void router.push({ name: 'devices' });
}

/**
 * Persist the device and return to the list.
 */
async function onSave(): Promise<void> {
	const error = firstError();
	if (error) {
		changeStep(error.step);
		$q.notify({ type: 'negative', message: error.message });
		return;
	}

	const input: DeviceInput = {
		name: draft.name,
		deviceId: draft.deviceId.trim() || generateUuid(),
		protocolId: draft.protocolId,
		enabled: draft.enabled,
		storeLogs: draft.storeLogs,
		config: { ...draft.config },
		attributes: { ...draft.attributes },
		events: draft.events.map((event) => ({ ...event })),
	};

	saving.value = true;
	try {
		if (editingId.value) await devicesStore.update(editingId.value, input);
		else await devicesStore.create(input);
		goBack();
	} catch (err) {
		$q.notify({ type: 'negative', message: err instanceof Error ? err.message : t('common.saveFailed') });
	} finally {
		saving.value = false;
	}
}

/**
 * Populate the draft from an existing device in edit mode.
 */
async function loadDevice(): Promise<void> {
	if (!editingId.value) return;
	loading.value = true;
	try {
		await devicesStore.fetch();
		const device: Device | undefined = devicesStore.getById(editingId.value);
		if (!device) return;
		Object.assign(draft, {
			name: device.name,
			deviceId: device.deviceId || generateUuid(),
			protocolId: device.protocolId,
			enabled: device.enabled,
			storeLogs: device.storeLogs,
			config: { ...device.config },
			attributes: { ...device.attributes },
			events: device.events.map((event) => ({ ...event })),
		});
	} finally {
		loading.value = false;
	}
}

useStepperNavigation({ currentStep, totalSteps: TOTAL_STEPS, changeStep });

/** LIFECYCLE HOOKS */
onMounted(() => {
	void loadDevice();
});
</script>

<style scoped lang="scss">
.step-hint {
	font-size: var(--mapex-font-sm);
	color: var(--mapex-text-secondary);
	margin-bottom: var(--mapex-spacing-xs);
}

.protocol-soon {
	border-radius: var(--mapex-radius-sm);
	background: var(--mapex-surface-sunken);
	color: var(--mapex-text-secondary);
}
</style>
