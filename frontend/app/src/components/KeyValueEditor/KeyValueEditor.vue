<template>
	<div class="kv-editor">
		<div v-for="(row, index) in rows" :key="index" class="kv-editor__row">
			<q-input
				class="kv-editor__field"
				:model-value="row.key"
				:placeholder="keyLabel ?? t('kv.key')"
				dense
				outlined
				hide-bottom-space
				@update:model-value="(v) => setRow(index, { key: String(v ?? '') })"
			/>
			<q-input
				class="kv-editor__field"
				:model-value="row.value"
				:placeholder="valueLabel ?? t('kv.value')"
				dense
				outlined
				hide-bottom-space
				@update:model-value="(v) => setRow(index, { value: String(v ?? '') })"
			/>
			<q-btn flat dense round icon="mdi-close" @click="removeRow(index)" />
		</div>

		<div v-if="!rows.length" class="kv-editor__empty">{{ emptyLabel ?? t('kv.empty') }}</div>

		<div>
			<q-btn flat dense no-caps icon="mdi-plus" :label="addLabel" @click="addRow" />
		</div>
	</div>
</template>

<script setup lang="ts">
/** TYPE IMPORTS */
import type { KeyValue } from '@services/sim';
import type { KeyValueEditorEmits, KeyValueEditorProps } from './interfaces';

/** COMPOSABLES */
import { useTranslations } from '@composables/i18n';

/** PROPS & EMITS */
const props = defineProps<KeyValueEditorProps>();
const emit = defineEmits<KeyValueEditorEmits>();

/** COMPOSABLES & STORES */
const { t } = useTranslations();

/** FUNCTIONS */

/**
 * Emit the rows with one entry's key or value changed.
 * @param {number} index - the row index
 * @param {Partial<KeyValue>} partial - the changed key or value
 */
function setRow(index: number, partial: Partial<KeyValue>): void {
	emit('update', props.rows.map((row, idx) => (idx === index ? { ...row, ...partial } : row)));
}

/**
 * Emit the rows with an empty entry appended.
 */
function addRow(): void {
	emit('update', [...props.rows, { key: '', value: '' }]);
}

/**
 * Emit the rows with one entry removed.
 * @param {number} index - the row index to remove
 */
function removeRow(index: number): void {
	emit('update', props.rows.filter((_, idx) => idx !== index));
}
</script>

<style scoped lang="scss">
.kv-editor {
	display: flex;
	flex-direction: column;
	gap: var(--mapex-spacing-xs);

	&__row {
		display: flex;
		align-items: center;
		gap: var(--mapex-spacing-sm);
	}

	&__field {
		flex: 1 1 0;
		min-width: 0;
	}

	&__empty {
		padding: var(--mapex-spacing-xs) 0;
		font-size: var(--mapex-font-xs);
		color: var(--mapex-text-muted);
		font-style: italic;
	}
}
</style>
