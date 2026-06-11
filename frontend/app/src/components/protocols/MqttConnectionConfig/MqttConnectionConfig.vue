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

		<q-separator class="mqtt-config__sep" />

		<q-toggle
			:model-value="modelValue.receiveEnabled ?? false"
			:label="t('connections.mqtt.receiveEnabled')"
			@update:model-value="(val) => patch({ receiveEnabled: !!val })"
		/>
		<div class="mqtt-config__none">{{ t('connections.mqtt.receiveHint') }}</div>

		<template v-if="modelValue.receiveEnabled">
			<div v-for="(sub, index) in subscriptions" :key="index" class="mqtt-config__sub">
				<q-input
					class="mqtt-config__field"
					:model-value="sub.name"
					:label="t('connections.mqtt.subName')"
					outlined
					dense
					stack-label
					hide-bottom-space
					@update:model-value="(val) => updateSub(index, { name: String(val ?? '') })"
				/>
				<q-input
					class="mqtt-config__field"
					:model-value="sub.topic"
					:label="`${t('connections.mqtt.subTopic')} *`"
					outlined
					dense
					stack-label
					hide-bottom-space
					@update:model-value="(val) => updateSub(index, { topic: String(val ?? '') })"
				/>
				<q-select
					class="mqtt-config__qos"
					:model-value="sub.qos"
					:options="qosOptions"
					:label="t('connections.mqtt.subQos')"
					outlined
					dense
					stack-label
					hide-bottom-space
					emit-value
					map-options
					@update:model-value="(val) => updateSub(index, { qos: val })"
				/>
				<q-btn flat round dense icon="mdi-close" :aria-label="t('common.remove')" @click="removeSub(index)" />
			</div>
			<q-btn flat dense icon="mdi-plus" :label="t('connections.mqtt.subAdd')" @click="addSub" />
		</template>
	</div>
</template>

<script setup lang="ts">
/** TYPE IMPORTS */
import type { MqttAuthMode, MqttQoS, MqttSubscription } from '@services/sim';
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

const subscriptions = computed<MqttSubscription[]>(() => props.modelValue.subscriptions ?? []);

const qosOptions: { label: string; value: MqttQoS }[] = [
	{ label: '0', value: 0 },
	{ label: '1', value: 1 },
	{ label: '2', value: 2 },
];

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

/**
 * Append a blank downlink subscription row.
 */
function addSub(): void {
	patch({ subscriptions: [...subscriptions.value, { name: '', topic: '', qos: 0 }] });
}

/**
 * Merge changed fields into one subscription row.
 * @param {number} index - the row index
 * @param {Partial<MqttSubscription>} partial - changed fields
 */
function updateSub(index: number, partial: Partial<MqttSubscription>): void {
	patch({ subscriptions: subscriptions.value.map((sub, i) => (i === index ? { ...sub, ...partial } : sub)) });
}

/**
 * Remove one subscription row.
 * @param {number} index - the row index
 */
function removeSub(index: number): void {
	patch({ subscriptions: subscriptions.value.filter((_, i) => i !== index) });
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

	&__sep {
		margin: var(--mapex-spacing-xs) 0;
	}

	&__sub {
		display: flex;
		gap: var(--mapex-spacing-sm);
		align-items: flex-start;
	}

	&__qos {
		flex: 0 0 120px;
	}
}
</style>
