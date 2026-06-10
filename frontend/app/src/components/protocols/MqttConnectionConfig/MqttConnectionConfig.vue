<template>
	<div class="mqtt-config">
		<q-input
			:model-value="modelValue.brokerUrl"
			:label="`${t('connections.mqtt.brokerUrl')} *`"
			:placeholder="t('connections.mqtt.brokerUrlPlaceholder')"
			outlined
			dense
			stack-label
			lazy-rules
			:rules="[requiredRule]"
			@update:model-value="(val) => patch({ brokerUrl: String(val ?? '') })"
			hide-bottom-space
		/>

		<div class="mqtt-config__pair">
			<q-input
				class="mqtt-config__field"
				:model-value="modelValue.clientId"
				:label="t('connections.mqtt.clientId')"
				outlined
				dense
				stack-label
				hide-bottom-space
				@update:model-value="(val) => patch({ clientId: String(val ?? '') })"
			/>
			<q-input
				class="mqtt-config__field"
				:model-value="modelValue.baseTopic"
				:label="t('connections.mqtt.baseTopic')"
				:placeholder="t('connections.mqtt.baseTopicPlaceholder')"
				outlined
				dense
				stack-label
				hide-bottom-space
				@update:model-value="(val) => patch({ baseTopic: String(val ?? '') })"
			/>
		</div>

		<q-select
			:model-value="modelValue.authMode"
			:options="authOptions"
			:label="t('connections.mqtt.authMode')"
			outlined
			dense
			stack-label
			hide-bottom-space
			emit-value
			map-options
			@update:model-value="(val) => patch({ authMode: val })"
		/>

		<template v-if="modelValue.authMode === 'userpass'">
			<q-input
				:model-value="modelValue.username"
				:label="`${t('connections.mqtt.username')} *`"
				outlined
				dense
				stack-label
				lazy-rules
				:rules="[requiredRule]"
				@update:model-value="(val) => patch({ username: String(val ?? '') })"
				hide-bottom-space
			/>
			<q-input
				:model-value="modelValue.password"
				:label="t('connections.mqtt.password')"
				type="password"
				outlined
				dense
				stack-label
				hide-bottom-space
				@update:model-value="(val) => patch({ password: String(val ?? '') })"
			/>
		</template>

		<template v-else-if="modelValue.authMode === 'tls'">
			<q-input
				:model-value="modelValue.tlsCertPem"
				:label="`${t('connections.mqtt.tlsCert')} *`"
				type="textarea"
				autogrow
				outlined
				dense
				stack-label
				class="mqtt-config__pem"
				lazy-rules
				:rules="[requiredRule]"
				@update:model-value="(val) => patch({ tlsCertPem: String(val ?? '') })"
				hide-bottom-space
			/>
			<q-input
				:model-value="modelValue.tlsKeyPem"
				:label="t('connections.mqtt.tlsKey')"
				type="textarea"
				autogrow
				outlined
				dense
				stack-label
				class="mqtt-config__pem"
				@update:model-value="(val) => patch({ tlsKeyPem: String(val ?? '') })"
				hide-bottom-space
			/>
			<q-input
				:model-value="modelValue.tlsCaPem"
				:label="t('connections.mqtt.tlsCa')"
				type="textarea"
				autogrow
				outlined
				dense
				stack-label
				class="mqtt-config__pem"
				@update:model-value="(val) => patch({ tlsCaPem: String(val ?? '') })"
				hide-bottom-space
			/>
		</template>

		<div v-else class="mqtt-config__none">{{ t('connections.mqtt.authNoneHint') }}</div>
	</div>
</template>

<script setup lang="ts">
/** TYPE IMPORTS */
import type { MqttAuthMode } from '@services/sim';
import type { MqttConnectionConfigEmits, MqttConnectionConfigProps } from './interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPOSABLES */
import { useTranslations } from '@composables/i18n';

/** PROPS & EMITS */
const props = defineProps<MqttConnectionConfigProps>();
const emit = defineEmits<MqttConnectionConfigEmits>();

/** COMPOSABLES & STORES */
const { t } = useTranslations();

/** COMPUTED */
const authOptions = computed<{ label: string; value: MqttAuthMode }[]>(() => [
	{ label: t('connections.mqtt.authNone'), value: 'none' },
	{ label: t('connections.mqtt.authUserPass'), value: 'userpass' },
	{ label: t('connections.mqtt.authTls'), value: 'tls' },
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
 * @param {Partial<MqttConnectionConfigProps['modelValue']>} partial - changed fields
 */
function patch(partial: Partial<MqttConnectionConfigProps['modelValue']>): void {
	emit('update:modelValue', { ...props.modelValue, ...partial });
}
</script>

<style scoped lang="scss">
.mqtt-config {
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

	&__pem {
		font-family: var(--mapex-mono-font);
	}

	&__none {
		padding: var(--mapex-spacing-md);
		text-align: center;
		font-size: var(--mapex-font-sm);
		color: var(--mapex-text-muted);
		font-style: italic;
	}
}
</style>
