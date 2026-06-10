<template>
	<div class="http-req">
		<!-- Request line -->
		<div class="http-req__line">
			<q-select
				class="http-req__method"
				:model-value="modelValue.method"
				:options="methodOptions"
				outlined
				dense
				emit-value
				map-options
				@update:model-value="(val) => patch({ method: val })"
				hide-bottom-space
			/>
			<q-input
				class="http-req__url"
				:model-value="modelValue.path"
				:placeholder="t('httpEvent.pathPlaceholder')"
				outlined
				dense
				@update:model-value="(val) => patch({ path: String(val ?? '') })"
				hide-bottom-space
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
			class="http-req__tabs text-grey-7"
		>
			<q-tab name="body" :label="t('httpEvent.tabBody')" />
			<q-tab name="headers" :label="`${t('httpEvent.tabHeaders')} (${modelValue.headers.length})`" />
			<slot name="tabs" />
		</q-tabs>
		<q-separator />

		<!-- Body -->
		<div v-if="tab === 'body'" class="http-req__panel">
			<RequestBodyEditor
				:model-value="{ bodyMode: modelValue.bodyMode, bodyFields: modelValue.bodyFields, body: modelValue.body }"
				@update:model-value="(body) => patch(body)"
			/>
		</div>

		<!-- Headers -->
		<div v-else-if="tab === 'headers'" class="http-req__panel">
			<KeyValueEditor
				:rows="modelValue.headers"
				:add-label="t('httpEvent.addHeader')"
				@update="(rows) => patch({ headers: rows })"
			/>
		</div>

		<!-- Extra panels contributed by the host (for example a schedule). -->
		<slot name="panel" :active="tab" />
	</div>
</template>

<script setup lang="ts">
/** TYPE IMPORTS */
import type { HttpMethod } from '@services/sim';
import type { HttpEventConfigEmits, HttpEventConfigProps } from './interfaces';

/** VUE IMPORTS */
import { computed, ref } from 'vue';

/** COMPONENTS */
import { KeyValueEditor } from '@components/KeyValueEditor';
import { RequestBodyEditor } from '@components/RequestBodyEditor';

/** COMPOSABLES */
import { useTranslations } from '@composables/i18n';

/** PROPS & EMITS */
const props = defineProps<HttpEventConfigProps>();
const emit = defineEmits<HttpEventConfigEmits>();

/** COMPOSABLES & STORES */
const { t } = useTranslations();

/** STATE */
// A host may inject extra tabs (such as a schedule), so the value is widened.
const tab = ref<string>('body');

/** COMPUTED */
const methodOptions = computed<{ label: string; value: HttpMethod }[]>(() => [
	{ label: 'POST', value: 'POST' },
	{ label: 'PUT', value: 'PUT' },
]);

/** FUNCTIONS */

/**
 * Emit a new event config with the changed fields merged in.
 * @param {Partial<HttpEventConfigProps['modelValue']>} partial - changed fields
 */
function patch(partial: Partial<HttpEventConfigProps['modelValue']>): void {
	emit('update:modelValue', { ...props.modelValue, ...partial });
}
</script>

<style scoped lang="scss">
.http-req {
	display: flex;
	flex-direction: column;
	gap: var(--mapex-spacing-md);

	&__line {
		display: flex;
		gap: var(--mapex-spacing-sm);
	}

	&__method {
		flex: 0 0 auto;
		min-width: 110px;
	}

	&__url {
		flex: 1 1 0;
		min-width: 0;
	}

	&__body {
		font-family: var(--mapex-mono-font);
	}

	&__panel {
		padding-top: var(--mapex-spacing-xs);
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
