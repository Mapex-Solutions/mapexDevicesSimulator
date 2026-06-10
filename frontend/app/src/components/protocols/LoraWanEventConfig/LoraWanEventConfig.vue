<template>
	<div class="lora-event">
		<div class="lora-event__line">
			<q-input
				class="lora-event__fport"
				:model-value="modelValue.fport"
				:label="t('lorawanEvent.fport')"
				type="number"
				min="1"
				max="223"
				outlined
				dense
				stack-label
				hide-bottom-space
				@update:model-value="(val) => patch({ fport: clampFport(val) })"
			/>
			<q-checkbox
				:model-value="modelValue.confirmed"
				:label="t('lorawanEvent.confirmed')"
				@update:model-value="(val) => patch({ confirmed: val })"
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
			<q-tab name="payload" :label="t('lorawanEvent.tabPayload')" />
			<slot name="tabs" />
		</q-tabs>
		<q-separator />

		<!-- Payload -->
		<div v-if="tab === 'payload'" class="lora-event__panel">
			<q-input
				:model-value="modelValue.payloadHex"
				:label="t('lorawanEvent.payloadHex')"
				:placeholder="t('lorawanEvent.payloadHexPlaceholder')"
				outlined
				dense
				stack-label
				hide-bottom-space
				class="lora-event__hex"
				@update:model-value="(val) => patch({ payloadHex: String(val ?? '') })"
			/>
			<TemplateHints />
		</div>

		<!-- Extra panels contributed by the host (for example a schedule). -->
		<slot name="panel" :active="tab" />
	</div>
</template>

<script setup lang="ts">
/** TYPE IMPORTS */
import type { LoraWanEventConfigEmits, LoraWanEventConfigProps } from './interfaces';

/** VUE IMPORTS */
import { ref } from 'vue';

/** COMPONENTS */
import { TemplateHints } from '@components/TemplateHints';

/** COMPOSABLES */
import { useTranslations } from '@composables/i18n';

/** PROPS & EMITS */
const props = defineProps<LoraWanEventConfigProps>();
const emit = defineEmits<LoraWanEventConfigEmits>();

/** COMPOSABLES & STORES */
const { t } = useTranslations();

/** STATE */
// A host may inject extra tabs (such as a schedule), so the value is widened.
const tab = ref<string>('payload');

/** FUNCTIONS */

/**
 * Clamp a raw FPort input to the valid application-port range (1-223).
 * @param {string | number | null} value - the raw field value
 * @returns {number} the clamped FPort
 */
function clampFport(value: string | number | null): number {
	return Math.min(223, Math.max(1, Math.floor(Number(value) || 1)));
}

/**
 * Emit a new event config with the changed fields merged in.
 * @param {Partial<LoraWanEventConfigProps['modelValue']>} partial - changed fields
 */
function patch(partial: Partial<LoraWanEventConfigProps['modelValue']>): void {
	emit('update:modelValue', { ...props.modelValue, ...partial });
}
</script>

<style scoped lang="scss">
.lora-event {
	display: flex;
	flex-direction: column;
	gap: var(--mapex-spacing-md);

	&__line {
		display: flex;
		align-items: center;
		gap: var(--mapex-spacing-md);
	}

	&__fport {
		width: 120px;
	}

	&__panel {
		display: flex;
		flex-direction: column;
		gap: var(--mapex-spacing-md);
		padding-top: var(--mapex-spacing-xs);
	}

	&__hex {
		font-family: var(--mapex-mono-font);
	}
}
</style>
