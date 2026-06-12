<template>
	<div class="bs-config">
		<div class="bs-config__pair">
			<q-input
				class="bs-config__field"
				:model-value="modelValue.lnsUri"
				:label="`${t('connections.basicstation.lnsUri')} *`"
				:placeholder="t('connections.basicstation.lnsUriPlaceholder')"
				outlined
				dense
				stack-label
				hide-bottom-space
				lazy-rules
				:rules="[requiredRule]"
				@update:model-value="(val) => patch({ lnsUri: String(val ?? '') })"
			/>
			<q-input
				class="bs-config__field"
				:model-value="modelValue.gatewayEui"
				:label="`${t('connections.basicstation.gatewayEui')} *`"
				:placeholder="t('connections.basicstation.gatewayEuiPlaceholder')"
				outlined
				dense
				stack-label
				hide-bottom-space
				lazy-rules
				:rules="[requiredRule]"
				@update:model-value="(val) => patch({ gatewayEui: String(val ?? '') })"
			/>
		</div>

		<div class="bs-config__pair">
			<q-select
				class="bs-config__field"
				:model-value="modelValue.region"
				:options="regionOptions"
				:label="t('connections.lorawan.region')"
				outlined
				dense
				stack-label
				hide-bottom-space
				emit-value
				map-options
				@update:model-value="(val) => patch({ region: val })"
			/>
			<q-select
				class="bs-config__field"
				:model-value="modelValue.macVersion"
				:options="macVersionOptions"
				:label="t('connections.lorawan.macVersion')"
				outlined
				dense
				stack-label
				hide-bottom-space
				emit-value
				map-options
				@update:model-value="(val) => patch({ macVersion: val })"
			/>
		</div>

		<div class="bs-config__pair">
			<q-select
				class="bs-config__field"
				:model-value="modelValue.class"
				:options="classOptions"
				:label="t('connections.lorawan.deviceClass')"
				outlined
				dense
				stack-label
				hide-bottom-space
				emit-value
				map-options
				@update:model-value="(val) => patch({ class: val })"
			/>
			<q-select
				class="bs-config__field"
				:model-value="modelValue.activation"
				:options="activationOptions"
				:label="t('connections.lorawan.activation')"
				outlined
				dense
				stack-label
				hide-bottom-space
				emit-value
				map-options
				@update:model-value="(val) => patch({ activation: val })"
			/>
		</div>

		<!-- OTAA -->
		<template v-if="modelValue.activation === 'otaa'">
			<div class="bs-config__pair">
				<q-input
					class="bs-config__field"
					:model-value="modelValue.devEui"
					:label="`${t('connections.lorawan.devEui')} *`"
					outlined
					dense
					stack-label
					hide-bottom-space
					lazy-rules
					:rules="[requiredRule]"
					@update:model-value="(val) => patch({ devEui: String(val ?? '') })"
				/>
				<q-input
					class="bs-config__field"
					:model-value="modelValue.joinEui"
					:label="t('connections.lorawan.joinEui')"
					outlined
					dense
					stack-label
					hide-bottom-space
					@update:model-value="(val) => patch({ joinEui: String(val ?? '') })"
				/>
			</div>
			<q-input
				:model-value="modelValue.appKey"
				:label="`${t('connections.lorawan.appKey')} *`"
				outlined
				dense
				stack-label
				hide-bottom-space
				lazy-rules
				:rules="[requiredRule]"
				@update:model-value="(val) => patch({ appKey: String(val ?? '') })"
			/>
			<q-input
				v-if="is11"
				:model-value="modelValue.nwkKey"
				:label="`${t('connections.lorawan.nwkKey')} *`"
				outlined
				dense
				stack-label
				hide-bottom-space
				lazy-rules
				:rules="[requiredRule]"
				@update:model-value="(val) => patch({ nwkKey: String(val ?? '') })"
			/>
		</template>

		<!-- ABP -->
		<template v-else>
			<q-input
				:model-value="modelValue.devAddr"
				:label="`${t('connections.lorawan.devAddr')} *`"
				outlined
				dense
				stack-label
				hide-bottom-space
				lazy-rules
				:rules="[requiredRule]"
				@update:model-value="(val) => patch({ devAddr: String(val ?? '') })"
			/>
			<q-input
				:model-value="modelValue.nwkSKey"
				:label="`${t('connections.lorawan.nwkSKey')} *`"
				outlined
				dense
				stack-label
				hide-bottom-space
				lazy-rules
				:rules="[requiredRule]"
				@update:model-value="(val) => patch({ nwkSKey: String(val ?? '') })"
			/>
			<q-input
				:model-value="modelValue.appSKey"
				:label="`${t('connections.lorawan.appSKey')} *`"
				outlined
				dense
				stack-label
				hide-bottom-space
				lazy-rules
				:rules="[requiredRule]"
				@update:model-value="(val) => patch({ appSKey: String(val ?? '') })"
			/>
		</template>
	</div>
</template>

<script setup lang="ts">
/** TYPE IMPORTS */
import type { GatewayRegion, LoraWanActivation, LoraWanMacVersion } from '@services/sim';
import type { BasicsStationConnectionConfigEmits, BasicsStationConnectionConfigProps } from './interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPOSABLES */
import { useTranslations } from '@composables/i18n';

const REGIONS: GatewayRegion[] = ['EU868', 'US915', 'AU915', 'AS923', 'CN470', 'IN865', 'KR920', 'RU864'];
const MAC_VERSIONS: LoraWanMacVersion[] = ['1.0.2', '1.0.3', '1.0.4', '1.1.0'];

/** PROPS & EMITS */
const props = defineProps<BasicsStationConnectionConfigProps>();
const emit = defineEmits<BasicsStationConnectionConfigEmits>();

/** COMPOSABLES & STORES */
const { t } = useTranslations();

/** COMPUTED */
const is11 = computed(() => props.modelValue.macVersion.startsWith('1.1'));

const regionOptions = computed(() => REGIONS.map((region) => ({ label: region, value: region })));
const macVersionOptions = computed(() => MAC_VERSIONS.map((version) => ({ label: version, value: version })));
const activationOptions = computed<{ label: string; value: LoraWanActivation }[]>(() => [
	{ label: t('connections.lorawan.otaa'), value: 'otaa' },
	{ label: t('connections.lorawan.abp'), value: 'abp' },
]);

const classOptions = computed(() => [
	{ label: t('connections.lorawan.classA'), value: 'A' as const },
	{ label: t('connections.lorawan.classC'), value: 'C' as const },
]);

/** FUNCTIONS */

/**
 * Validation rule for a required field.
 * @param {string | null | undefined} value - the field value
 * @returns {true | string} true when valid, otherwise the error message
 */
function requiredRule(value: string | null | undefined): true | string {
	return (!!value && value.trim().length > 0) || t('validation.required');
}

/**
 * Emit a new config object with the changed fields merged in.
 * @param {Partial<BasicsStationConnectionConfigProps['modelValue']>} partial - changed fields
 */
function patch(partial: Partial<BasicsStationConnectionConfigProps['modelValue']>): void {
	emit('update:modelValue', { ...props.modelValue, ...partial });
}
</script>

<style scoped lang="scss">
.bs-config {
	display: flex;
	flex-direction: column;
	gap: var(--mapex-spacing-md);

	&__pair {
		display: flex;
		gap: var(--mapex-spacing-md);
	}

	&__field {
		flex: 1 1 0;
		min-width: 0;
	}
}
</style>
