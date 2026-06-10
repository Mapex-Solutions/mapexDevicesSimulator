<template>
	<div class="filter-bar">
		<div class="filter-bar__grid">
			<div v-for="field in fields" :key="field.key" class="filter-bar__field">
				<label class="filter-bar__label">{{ t(field.labelKey) }}</label>

				<q-select
					v-if="field.type === 'select'"
					:model-value="modelValue[field.key] ?? ''"
					:options="optionsFor(field)"
					dense
					outlined
					clearable
					emit-value
					map-options
					@update:model-value="(val) => patch(field.key, val)"
					hide-bottom-space
				/>
				<q-input
					v-else
					:model-value="modelValue[field.key] ?? ''"
					dense
					outlined
					clearable
					@update:model-value="(val) => patch(field.key, val)"
					hide-bottom-space
				/>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
/** TYPE IMPORTS */
import type { FilterField } from '@utils/message-filters';
import type { MessageFilterBarEmits, MessageFilterBarProps } from './interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPOSABLES */
import { useTranslations } from '@composables/i18n';

/** UTILS */
import { getMessageFilterFields } from '@utils/message-filters';

/** PROPS & EMITS */
const props = defineProps<MessageFilterBarProps>();
const emit = defineEmits<MessageFilterBarEmits>();

/** COMPOSABLES & STORES */
const { t } = useTranslations();

/** COMPUTED */
const fields = computed(() => getMessageFilterFields(props.protocol));

/** FUNCTIONS */

/**
 * Map a field's options to Quasar select options with resolved labels.
 * @param {FilterField} field - the select field
 */
function optionsFor(field: FilterField): { label: string; value: string }[] {
	return (field.options ?? []).map((option) => ({
		label: option.labelKey ? t(option.labelKey) : (option.label ?? option.value),
		value: option.value,
	}));
}

/**
 * Emit a new values object with one field changed.
 * @param {string} key - the field key
 * @param {string | number | null} value - the new value (null clears)
 */
function patch(key: string, value: string | number | null): void {
	emit('update:modelValue', { ...props.modelValue, [key]: value == null ? '' : String(value) });
}
</script>

<style scoped lang="scss">
.filter-bar {
	&__grid {
		display: grid;
		grid-template-columns: repeat(2, minmax(0, 1fr));
		gap: var(--mapex-spacing-md);
	}

	&__field {
		display: flex;
		flex-direction: column;
		gap: var(--mapex-spacing-2xs);
	}

	&__label {
		font-size: var(--mapex-font-xs);
		font-weight: var(--mapex-font-weight-medium);
		color: var(--mapex-text-secondary);
	}
}
</style>
