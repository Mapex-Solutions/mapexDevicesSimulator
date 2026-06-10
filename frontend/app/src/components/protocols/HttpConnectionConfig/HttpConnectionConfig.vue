<template>
	<div class="http-req">
		<!-- Request line -->
		<div class="http-req__line">
			<q-select
				class="http-req__method"
				:model-value="modelValue.method"
				:options="methodOptions"
				:label="t('connections.http.method')"
				outlined
				dense
				stack-label
				emit-value
				map-options
				@update:model-value="(val) => patch({ method: val })"
				hide-bottom-space
			/>
			<q-input
				class="http-req__url"
				:model-value="modelValue.url"
				:label="`${t('connections.http.url')} *`"
				:placeholder="t('connections.http.urlPlaceholder')"
				outlined
				dense
				stack-label
				lazy-rules
				:rules="[requiredRule]"
				@update:model-value="(val) => patch({ url: String(val ?? '') })"
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
			<q-tab name="auth" :label="t('connections.http.tabAuth')" />
			<q-tab name="headers" :label="`${t('connections.http.tabHeaders')} (${modelValue.headers.length})`" />
		</q-tabs>
		<q-separator />

		<!-- Authorization -->
		<div v-if="tab === 'auth'" class="http-req__panel">
			<q-select
				:model-value="modelValue.authMode"
				:options="authOptions"
				:label="t('connections.http.authMode')"
				outlined
				dense
				stack-label
				hide-bottom-space
				emit-value
				map-options
				@update:model-value="(val) => patch({ authMode: val })"
			/>

			<template v-if="modelValue.authMode === 'apiKey'">
				<q-input
					:model-value="modelValue.apiKeyHeader"
					:label="t('connections.http.apiKeyHeader')"
					outlined
					dense
					stack-label
					hide-bottom-space
					@update:model-value="(val) => patch({ apiKeyHeader: String(val ?? '') })"
				/>
				<q-input
					:model-value="modelValue.apiKey"
					:label="`${t('connections.http.apiKey')} *`"
					type="password"
					outlined
					dense
					stack-label
					lazy-rules
					:rules="[requiredRule]"
					@update:model-value="(val) => patch({ apiKey: String(val ?? '') })"
					hide-bottom-space
				/>
			</template>

			<template v-else-if="modelValue.authMode === 'basic'">
				<q-input
					:model-value="modelValue.basicUser"
					:label="`${t('connections.http.basicUser')} *`"
					outlined
					dense
					stack-label
					lazy-rules
					:rules="[requiredRule]"
					@update:model-value="(val) => patch({ basicUser: String(val ?? '') })"
					hide-bottom-space
				/>
				<q-input
					:model-value="modelValue.basicPass"
					:label="`${t('connections.http.basicPass')} *`"
					type="password"
					outlined
					dense
					stack-label
					lazy-rules
					:rules="[requiredRule]"
					@update:model-value="(val) => patch({ basicPass: String(val ?? '') })"
					hide-bottom-space
				/>
			</template>

			<div v-else class="http-req__none">{{ t('connections.http.authNoneHint') }}</div>
		</div>

		<!-- Headers -->
		<div v-else class="http-req__panel">
			<KeyValueEditor
				:rows="modelValue.headers"
				:add-label="t('connections.http.addHeader')"
				@update="(rows) => patch({ headers: rows })"
			/>
		</div>
	</div>
</template>

<script setup lang="ts">
/** TYPE IMPORTS */
import type { HttpAuthMode, HttpMethod } from '@services/sim';
import type { HttpConnectionConfigEmits, HttpConnectionConfigProps } from './interfaces';

/** VUE IMPORTS */
import { computed, ref } from 'vue';

/** COMPONENTS */
import { KeyValueEditor } from '@components/KeyValueEditor';

/** COMPOSABLES */
import { useTranslations } from '@composables/i18n';

/** PROPS & EMITS */
const props = defineProps<HttpConnectionConfigProps>();
const emit = defineEmits<HttpConnectionConfigEmits>();

/** COMPOSABLES & STORES */
const { t } = useTranslations();

/** STATE */
const tab = ref<'auth' | 'headers'>('auth');

/** COMPUTED */
const methodOptions = computed<{ label: string; value: HttpMethod }[]>(() => [
	{ label: 'POST', value: 'POST' },
	{ label: 'PUT', value: 'PUT' },
]);

const authOptions = computed<{ label: string; value: HttpAuthMode }[]>(() => [
	{ label: t('connections.http.authNone'), value: 'none' },
	{ label: t('connections.http.authApiKey'), value: 'apiKey' },
	{ label: t('connections.http.authBasic'), value: 'basic' },
]);

/** FUNCTIONS */

/**
 * Validation rule for a required credential field.
 * @param {string | null | undefined} value - the field value
 * @returns {true | string} true when valid, otherwise the error message
 */
function requiredRule(value: string | null | undefined): true | string {
	return (!!value && value.trim().length > 0) || t('validation.required');
}

/**
 * Emit a new config object with the changed fields merged in.
 * @param {Partial<HttpConnectionConfigProps['modelValue']>} partial - changed fields
 */
function patch(partial: Partial<HttpConnectionConfigProps['modelValue']>): void {
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

	&__panel {
		display: flex;
		flex-direction: column;
		gap: var(--mapex-spacing-md);
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
