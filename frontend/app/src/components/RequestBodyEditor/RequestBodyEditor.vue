<template>
	<div class="body-editor">
		<div class="body-editor__modes">
			<q-radio
				v-for="opt in bodyOptions"
				:key="opt.value"
				:model-value="modelValue.bodyMode"
				:val="opt.value"
				:label="opt.label"
				dense
				@update:model-value="(val) => patch({ bodyMode: val })"
			/>
		</div>

		<template v-if="modelValue.bodyMode === 'raw'">
			<q-input
				:model-value="modelValue.body"
				type="textarea"
				autogrow
				outlined
				dense
				class="body-editor__raw"
				input-style="min-height: 140px"
				:rules="[jsonRule]"
				lazy-rules
				@update:model-value="(val) => patch({ body: String(val ?? '') })"
			/>
			<TemplateHints />
		</template>

		<template v-else-if="modelValue.bodyMode === 'form'">
			<KeyValueEditor
				:rows="modelValue.bodyFields"
				:add-label="t('httpEvent.addField')"
				:key-label="t('httpEvent.fieldName')"
				:value-label="t('httpEvent.fieldValue')"
				@update="(rows) => patch({ bodyFields: rows })"
			/>
			<TemplateHints />
		</template>

		<div v-else class="body-editor__none">{{ t('httpEvent.bodyNoneHint') }}</div>
	</div>
</template>

<script setup lang="ts">
/** TYPE IMPORTS */
import type { HttpBodyMode } from '@services/sim';
import type { RequestBodyEditorEmits, RequestBodyEditorProps } from './interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPONENTS */
import { KeyValueEditor } from '@components/KeyValueEditor';
import { TemplateHints } from '@components/TemplateHints';

/** COMPOSABLES */
import { useTranslations } from '@composables/i18n';

/** UTILS */
import { validateJsonBody } from '@utils/template';

/** PROPS & EMITS */
const props = defineProps<RequestBodyEditorProps>();
const emit = defineEmits<RequestBodyEditorEmits>();

/** COMPOSABLES & STORES */
const { t } = useTranslations();

/** COMPUTED */
const bodyOptions = computed<{ label: string; value: HttpBodyMode }[]>(() => [
	{ label: t('httpEvent.bodyNone'), value: 'none' },
	{ label: t('httpEvent.bodyRaw'), value: 'raw' },
	{ label: t('httpEvent.bodyForm'), value: 'form' },
]);

/** FUNCTIONS */

/**
 * Emit the body config with the changed fields merged in.
 * @param {Partial<RequestBodyEditorProps['modelValue']>} partial - changed fields
 */
function patch(partial: Partial<RequestBodyEditorProps['modelValue']>): void {
	emit('update:modelValue', { ...props.modelValue, ...partial });
}

/**
 * Validate the raw body as JSON, tolerating {{ template }} placeholders. Returns
 * true when valid, or a localized error message (with the parser detail) when not.
 * @param {string | number | null} value - the textarea value
 * @returns {true | string} true when valid, otherwise the error message
 */
function jsonRule(value: string | number | null): true | string {
	const result = validateJsonBody(String(value ?? ''));
	return result.valid || `${t('validation.jsonInvalid')} — ${result.detail}`;
}
</script>

<style scoped lang="scss">
.body-editor {
	display: flex;
	flex-direction: column;
	gap: var(--mapex-spacing-md);

	&__modes {
		display: flex;
		align-items: center;
		gap: var(--mapex-spacing-lg);
	}

	&__raw {
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
