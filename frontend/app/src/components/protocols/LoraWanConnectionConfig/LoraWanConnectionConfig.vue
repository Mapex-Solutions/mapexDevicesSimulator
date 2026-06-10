<template>
	<div class="lora-config">
		<q-select
			:model-value="modelValue.gatewayId"
			:options="gatewayOptions"
			:label="`${t('connections.lorawan.gateway')} *`"
			:hint="gatewayOptions.length ? undefined : t('connections.lorawan.noGateways')"
			outlined
			dense
			stack-label
			hide-bottom-space
			emit-value
			map-options
			@update:model-value="(val) => patch({ gatewayId: String(val ?? '') })"
		/>

		<div class="lora-config__pair">
			<q-select
				class="lora-config__field"
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
				class="lora-config__field"
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

		<q-select
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

		<!-- OTAA -->
		<template v-if="modelValue.activation === 'otaa'">
			<div class="lora-config__pair">
				<q-input
					class="lora-config__field"
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
					class="lora-config__field"
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
import type { LoraWanConnectionConfigEmits, LoraWanConnectionConfigProps } from './interfaces';

/** VUE IMPORTS */
import { computed, onMounted } from 'vue';

/** COMPOSABLES */
import { useTranslations } from '@composables/i18n';

/** STORES */
import { useGatewaysStore } from '@stores/gateways';

const REGIONS: GatewayRegion[] = ['EU868', 'US915', 'AU915', 'AS923', 'CN470', 'IN865', 'KR920', 'RU864'];
const MAC_VERSIONS: LoraWanMacVersion[] = ['1.0.2', '1.0.3', '1.0.4', '1.1.0'];

/** PROPS & EMITS */
const props = defineProps<LoraWanConnectionConfigProps>();
const emit = defineEmits<LoraWanConnectionConfigEmits>();

/** COMPOSABLES & STORES */
const { t } = useTranslations();
const gatewaysStore = useGatewaysStore();

/** COMPUTED */
const is11 = computed(() => props.modelValue.macVersion.startsWith('1.1'));

const gatewayOptions = computed(() =>
	gatewaysStore.items.map((gateway) => ({
		label: gateway.eui ? `${gateway.name} · ${gateway.eui}` : gateway.name,
		value: gateway.id,
	})),
);
const regionOptions = computed(() => REGIONS.map((region) => ({ label: region, value: region })));
const macVersionOptions = computed(() => MAC_VERSIONS.map((version) => ({ label: version, value: version })));
const activationOptions = computed<{ label: string; value: LoraWanActivation }[]>(() => [
	{ label: t('connections.lorawan.otaa'), value: 'otaa' },
	{ label: t('connections.lorawan.abp'), value: 'abp' },
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
 * @param {Partial<LoraWanConnectionConfigProps['modelValue']>} partial - changed fields
 */
function patch(partial: Partial<LoraWanConnectionConfigProps['modelValue']>): void {
	emit('update:modelValue', { ...props.modelValue, ...partial });
}

/** LIFECYCLE HOOKS */
onMounted(() => {
	void gatewaysStore.fetch();
});
</script>

<style scoped lang="scss">
.lora-config {
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
