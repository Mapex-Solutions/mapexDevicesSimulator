<template>
	<q-page class="q-pa-lg">
		<div v-if="loading" class="flex flex-center q-pa-xl">
			<q-spinner color="primary" size="42px" />
		</div>

		<div v-else>
			<PageHeader
				icon="router"
				icon-color="primary"
				:title="pageTitle"
				:description="t('createGateway.description')"
				:button="{ label: t('createGateway.back'), icon: 'arrow_back', flat: true, onClick: goBack }"
			/>

			<div class="row q-col-gutter-lg">
				<!-- Stepper (left) -->
				<div class="col-12 col-md-4">
					<StepperVertical
						:title="t('createGateway.stepperTitle')"
						:subtitle="t('createGateway.stepperSubtitle')"
						:info-text="t('createGateway.stepperInfo')"
						:current-step-label="t('createGateway.currentStep')"
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
								class="col-12 col-sm-6"
								v-model="draft.name"
								:label="`${t('gateways.fields.name')} *`"
								outlined
								dense
								lazy-rules
								:rules="[requiredRule]"
								hide-bottom-space
							/>
							<q-select
								class="col-12 col-sm-6"
								v-model="draft.enabled"
								:options="boolOptions"
								:label="t('gateways.fields.status')"
								outlined
								dense
								emit-value
								map-options
								hide-bottom-space
							/>
							<q-input
								class="col-12 col-sm-6"
								v-model="draft.eui"
								:label="`${t('gateways.fields.eui')} *`"
								:placeholder="t('gateways.fields.euiPlaceholder')"
								outlined
								dense
								lazy-rules
								:rules="[requiredRule]"
								hide-bottom-space
							/>
							<q-select
								class="col-12 col-sm-6"
								v-model="draft.region"
								:options="regionOptions"
								:label="t('gateways.fields.region')"
								outlined
								dense
								emit-value
								map-options
								hide-bottom-space
							/>
							<q-input
								class="col-12"
								v-model="draft.description"
								:label="t('gateways.fields.description')"
								type="textarea"
								autogrow
								outlined
								dense
								hide-bottom-space
							/>
						</div>

						<!-- Step 2: LNS link -->
						<div v-else-if="currentStep === 2" class="column q-gutter-md">
							<div class="step-hint">{{ t('createGateway.linkHint') }}</div>
							<q-select
								v-model="draft.link.protocol"
								:options="protocolOptions"
								:label="t('gatewayLink.protocol')"
								outlined
								dense
								emit-value
								map-options
								hide-bottom-space
							/>

							<q-input
								v-if="draft.link.protocol === 'basicstation'"
								v-model="draft.link.lnsUri"
								:label="`${t('gatewayLink.lnsUri')} *`"
								:placeholder="t('gatewayLink.lnsUriPlaceholder')"
								outlined
								dense
								lazy-rules
								:rules="[requiredRule]"
								hide-bottom-space
							/>

							<div v-else class="gw-link__pair">
								<q-input
									class="gw-link__host"
									v-model="draft.link.host"
									:label="`${t('gatewayLink.host')} *`"
									:placeholder="t('gatewayLink.hostPlaceholder')"
									outlined
									dense
									lazy-rules
									:rules="[requiredRule]"
									hide-bottom-space
								/>
								<q-input
									class="gw-link__port"
									:model-value="draft.link.port"
									:label="t('gatewayLink.port')"
									type="number"
									min="1"
									outlined
									dense
									@update:model-value="(val) => (draft.link.port = clampPort(val))"
									hide-bottom-space
								/>
							</div>
						</div>
					</template>
				</FormCard>
			</div>
		</div>
	</q-page>
</template>

<script setup lang="ts">
/** TYPE IMPORTS */
import type { Gateway, GatewayInput, GatewayLinkProtocol, GatewayRegion } from '@services/sim';
import type { FormCardHeader, FormCardNavigation } from '@components/FormCard';
import type { StepperVerticalItem } from '@components/StepperVertical';

/** VUE IMPORTS */
import { computed, onMounted, reactive, ref } from 'vue';

/** COMPONENTS */
import { PageHeader } from '@components/PageHeader';
import { FormCard } from '@components/FormCard';
import { StepperVertical } from '@components/StepperVertical';

/** COMPOSABLES */
import { useTranslations } from '@composables/i18n';
import { useStepperNavigation } from '@composables/shared/form';

/** UTILS */
import { notifyFail } from '@utils/alert';

/** SERVICES */
import { useRoute, useRouter } from 'vue-router';

/** STORES */
import { useGatewaysStore } from '@stores/gateways';

const TOTAL_STEPS = 2;
const REGIONS: GatewayRegion[] = ['EU868', 'US915', 'AU915', 'AS923', 'CN470', 'IN865', 'KR920', 'RU864'];
const LINK_PROTOCOLS: GatewayLinkProtocol[] = ['basicstation', 'udp'];

/** COMPOSABLES & STORES */
const { t } = useTranslations();
const route = useRoute();
const router = useRouter();
const gatewaysStore = useGatewaysStore();

/** STATE */
const editingId = ref<string | null>((route.params.id as string | undefined) ?? null);
const isEditMode = computed(() => editingId.value !== null);
const loading = ref(false);
const saving = ref(false);
const currentStep = ref(1);
const draft = reactive<GatewayInput>(blankDraft());

/** COMPUTED */
const boolOptions = computed(() => [
	{ label: t('devices.on'), value: true },
	{ label: t('devices.off'), value: false },
]);
const regionOptions = computed(() => REGIONS.map((region) => ({ label: region, value: region })));
const protocolOptions = computed(() =>
	LINK_PROTOCOLS.map((protocol) => ({ label: t(`gatewayLink.${protocol}`), value: protocol })),
);

const steps = computed<StepperVerticalItem[]>(() => [
	{ title: t('createGateway.steps.identity'), description: t('createGateway.steps.identityDesc'), icon: 'badge' },
	{ title: t('createGateway.steps.link'), description: t('createGateway.steps.linkDesc'), icon: 'cloud_sync' },
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

const pageTitle = computed(() => (isEditMode.value ? t('createGateway.editTitle') : t('createGateway.title')));

/** FUNCTIONS */

/**
 * Build a blank gateway draft.
 */
function blankDraft(): GatewayInput {
	return {
		name: '',
		eui: '',
		enabled: true,
		region: 'EU868',
		description: '',
		link: { protocol: 'basicstation', lnsUri: 'wss://127.0.0.1:1887', host: '127.0.0.1', port: 1700 },
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
 * Clamp a raw port input to a positive integer.
 * @param {string | number | null} value - the raw field value
 * @returns {number} the clamped port
 */
function clampPort(value: string | number | null): number {
	return Math.max(1, Math.floor(Number(value) || 1));
}

/**
 * Validation message for a single step, or null when the step is valid.
 * @param {number} step - the step to validate
 * @returns {string | null} the error message or null
 */
function stepError(step: number): string | null {
	if (step === 1) {
		if (!draft.name.trim()) return t('validation.nameRequired');
		if (!draft.eui.trim()) return t('validation.euiRequired');
	}

	if (step === 2) {
		if (draft.link.protocol === 'basicstation' && !draft.link.lnsUri.trim()) return t('validation.linkRequired');
		if (draft.link.protocol === 'udp' && !draft.link.host.trim()) return t('validation.linkRequired');
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
			notifyFail({ message });
			return;
		}
	}
	changeStep(step);
}

/**
 * Move to a step, clamped to the valid range.
 * @param {number} step - the target step
 */
function changeStep(step: number): void {
	currentStep.value = Math.min(TOTAL_STEPS, Math.max(1, step));
}

function goBack(): void {
	void router.push({ name: 'gateways' });
}

/**
 * Persist the gateway and return to the list.
 */
async function onSave(): Promise<void> {
	const error = firstError();
	if (error) {
		changeStep(error.step);
		notifyFail({ message: error.message });
		return;
	}

	const input: GatewayInput = {
		name: draft.name,
		eui: draft.eui,
		enabled: draft.enabled,
		region: draft.region,
		description: draft.description,
		link: { ...draft.link },
	};

	saving.value = true;
	try {
		if (editingId.value) await gatewaysStore.update(editingId.value, input);
		else await gatewaysStore.create(input);
		goBack();
	} catch (err) {
		notifyFail({ message: err instanceof Error ? err.message : t('common.saveFailed') });
	} finally {
		saving.value = false;
	}
}

/**
 * Populate the draft from an existing gateway in edit mode.
 */
async function loadGateway(): Promise<void> {
	if (!editingId.value) return;
	loading.value = true;
	try {
		await gatewaysStore.fetch();
		const gateway: Gateway | undefined = gatewaysStore.getById(editingId.value);
		if (!gateway) return;
		Object.assign(draft, {
			name: gateway.name,
			eui: gateway.eui,
			enabled: gateway.enabled,
			region: gateway.region,
			description: gateway.description,
			link: { ...gateway.link },
		});
	} finally {
		loading.value = false;
	}
}

useStepperNavigation({ currentStep, totalSteps: TOTAL_STEPS, changeStep });

/** LIFECYCLE HOOKS */
onMounted(() => {
	void loadGateway();
});
</script>

<style scoped lang="scss">
.step-hint {
	font-size: var(--mapex-font-sm);
	color: var(--mapex-text-secondary);
	margin-bottom: var(--mapex-spacing-xs);
}

.gw-link {
	&__pair {
		display: flex;
		gap: var(--mapex-spacing-md);
	}

	&__host {
		flex: 1 1 0;
		min-width: 0;
	}

	&__port {
		flex: 0 0 140px;
	}
}
</style>
