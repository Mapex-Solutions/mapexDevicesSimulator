<template>
	<div class="mqtt-event">
		<q-input
			:model-value="modelValue.topic"
			:label="`${t('mqttEvent.topic')} *`"
			:placeholder="t('mqttEvent.topicPlaceholder')"
			outlined
			dense
			stack-label
			@update:model-value="(val) => patch({ topic: String(val ?? '') })"
			hide-bottom-space
		/>

		<div class="mqtt-event__line">
			<q-select
				class="mqtt-event__qos"
				:model-value="modelValue.qos"
				:options="qosOptions"
				:label="t('mqttEvent.qos')"
				outlined
				dense
				stack-label
				emit-value
				map-options
				@update:model-value="(val) => patch({ qos: val })"
				hide-bottom-space
			/>
			<q-checkbox
				:model-value="modelValue.retain"
				:label="t('mqttEvent.retain')"
				@update:model-value="(val) => patch({ retain: val })"
			/>
		</div>

		<!-- Tabs -->
		<q-tabs
			v-model="tab"
			dense
			no-caps
			align="left"
			active-color="primary"
			indicator-color="primary"
			class="text-grey-7"
		>
			<q-tab name="body" :label="t('mqttEvent.tabPayload')" />
			<slot name="tabs" />
		</q-tabs>
		<q-separator />

		<!-- Payload -->
		<div v-if="tab === 'body'" class="mqtt-event__panel">
			<RequestBodyEditor
				:model-value="{ bodyMode: modelValue.bodyMode, bodyFields: modelValue.bodyFields, body: modelValue.body }"
				@update:model-value="(body) => patch(body)"
			/>
		</div>

		<!-- Extra panels contributed by the host (for example a schedule). -->
		<slot name="panel" :active="tab" />
	</div>
</template>

<script setup lang="ts">
/** TYPE IMPORTS */
import type { MqttQoS } from '@services/sim';
import type { MqttEventConfigEmits, MqttEventConfigProps } from './interfaces';

/** VUE IMPORTS */
import { computed, ref } from 'vue';

/** COMPONENTS */
import { RequestBodyEditor } from '@components/RequestBodyEditor';

/** COMPOSABLES */
import { useTranslations } from '@composables/i18n';

/** PROPS & EMITS */
const props = defineProps<MqttEventConfigProps>();
const emit = defineEmits<MqttEventConfigEmits>();

/** COMPOSABLES & STORES */
const { t } = useTranslations();

/** STATE */
// A host may inject extra tabs (such as a schedule), so the value is widened.
const tab = ref<string>('body');

/** COMPUTED */
const qosOptions = computed<{ label: string; value: MqttQoS }[]>(() => [
	{ label: '0', value: 0 },
	{ label: '1', value: 1 },
	{ label: '2', value: 2 },
]);

/** FUNCTIONS */

/**
 * Emit a new event config with the changed fields merged in.
 * @param {Partial<MqttEventConfigProps['modelValue']>} partial - changed fields
 */
function patch(partial: Partial<MqttEventConfigProps['modelValue']>): void {
	emit('update:modelValue', { ...props.modelValue, ...partial });
}
</script>

<style scoped lang="scss">
.mqtt-event {
	display: flex;
	flex-direction: column;
	gap: var(--mapex-spacing-md);

	&__line {
		display: flex;
		align-items: center;
		gap: var(--mapex-spacing-md);
	}

	&__qos {
		min-width: 120px;
	}

	&__panel {
		padding-top: var(--mapex-spacing-xs);
	}
}
</style>
