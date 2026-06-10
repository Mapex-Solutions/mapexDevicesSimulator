<template>
	<div class="detail">
		<div class="detail__header">
			<span class="detail__title">{{ t('console.details') }}</span>
		</div>

		<div v-if="!message" class="detail__empty">
			<q-icon name="mdi-cursor-default-click-outline" size="32px" />
			<span>{{ t('console.noSelection') }}</span>
		</div>

		<div v-else class="detail__body">
			<div class="row items-center q-gutter-xs q-mb-sm">
				<q-chip dense :color="dirColor" text-color="white" :icon="dirIcon" :label="dirLabel" />
				<q-chip dense outline :label="t(`protocol.${message.protocol}`)" />
				<q-space />
				<span class="detail__ts">{{ message.ts }}</span>
			</div>

			<div class="detail__device">{{ message.deviceName }}</div>
			<div class="detail__sub">
				{{ message.summary }}
				<q-badge v-if="message.status" color="grey-7" :label="message.status" class="q-ml-xs" />
			</div>

			<div class="detail__section-head row items-center justify-between">
				<span>{{ t('console.payload') }}</span>
				<q-btn flat dense size="sm" no-caps icon="mdi-content-copy" :label="t('console.copy')" @click="copy(message.payload)" />
			</div>
			<pre class="detail__payload">{{ message.payload }}</pre>

			<template v-if="metaEntries.length">
				<div class="detail__section-head">{{ t('console.meta') }}</div>
				<div class="detail__meta">
					<div v-for="[key, value] in metaEntries" :key="key" class="detail__meta-row">
						<span class="detail__meta-key">{{ key }}</span>
						<span class="detail__meta-value">{{ value }}</span>
					</div>
				</div>
			</template>
		</div>
	</div>
</template>

<script setup lang="ts">
/** TYPE IMPORTS */
import type { MessageDetailProps } from './interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPOSABLES */
import { useTranslations } from '@composables/i18n';

/** UTILS */
import { copyToClipboard, useQuasar } from 'quasar';

/** PROPS & EMITS */
const props = defineProps<MessageDetailProps>();

/** COMPOSABLES & STORES */
const { t } = useTranslations();
const $q = useQuasar();

/** COMPUTED */
const metaEntries = computed(() => Object.entries(props.message?.meta ?? {}));

const dirColor = computed(() => {
	if (props.message?.direction === 'up') return 'teal';
	if (props.message?.direction === 'down') return 'primary';
	return 'grey-7';
});

const dirIcon = computed(() => {
	if (props.message?.direction === 'up') return 'mdi-arrow-up';
	if (props.message?.direction === 'down') return 'mdi-arrow-down';
	return 'mdi-cog-outline';
});

const dirLabel = computed(() => {
	if (props.message?.direction === 'up') return t('console.dirUp');
	if (props.message?.direction === 'down') return t('console.dirDown');
	return t('console.dirSystem');
});

/** FUNCTIONS */

/**
 * Copy a value to the clipboard and confirm.
 * @param {string} text - the value to copy
 */
async function copy(text: string): Promise<void> {
	await copyToClipboard(text);
	$q.notify({ type: 'positive', message: t('console.copied') });
}
</script>

<style scoped lang="scss">
.detail {
	display: flex;
	flex-direction: column;
	height: 100%;
	background: var(--mapex-surface-elevated);

	&__header {
		padding: var(--mapex-spacing-md) var(--mapex-spacing-lg);
		border-bottom: 1px solid var(--mapex-divider);
	}

	&__title {
		font-size: var(--mapex-font-sm);
		font-weight: var(--mapex-font-weight-semibold);
		color: var(--mapex-text-primary);
	}

	&__empty {
		flex: 1;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: var(--mapex-spacing-md);
		color: var(--mapex-text-muted);
		padding: var(--mapex-spacing-2xl);
		text-align: center;
	}

	&__body {
		flex: 1;
		overflow-y: auto;
		padding: var(--mapex-spacing-lg);
	}

	&__ts {
		font-size: var(--mapex-font-xs);
		color: var(--mapex-text-muted);
		font-family: var(--mapex-mono-font);
	}

	&__device {
		font-size: var(--mapex-font-lg);
		font-weight: var(--mapex-font-weight-semibold);
		color: var(--mapex-text-primary);
	}

	&__sub {
		font-size: var(--mapex-font-sm);
		color: var(--mapex-text-secondary);
		margin-top: var(--mapex-spacing-2xs);
	}

	&__section-head {
		margin-top: var(--mapex-spacing-lg);
		margin-bottom: var(--mapex-spacing-xs);
		font-size: var(--mapex-font-xs);
		text-transform: uppercase;
		letter-spacing: 0.04em;
		color: var(--mapex-text-secondary);
	}

	&__payload {
		margin: 0;
		padding: var(--mapex-spacing-md);
		border-radius: var(--mapex-radius-sm);
		background: var(--mapex-surface-sunken);
		color: var(--mapex-text-primary);
		font-family: var(--mapex-mono-font);
		font-size: var(--mapex-font-xs);
		white-space: pre-wrap;
		word-break: break-word;
	}

	&__meta-row {
		display: flex;
		justify-content: space-between;
		gap: var(--mapex-spacing-md);
		padding: 4px 0;
		border-bottom: 1px solid var(--mapex-divider);
		font-size: var(--mapex-font-sm);
	}

	&__meta-key {
		color: var(--mapex-text-secondary);
	}

	&__meta-value {
		color: var(--mapex-text-primary);
		font-family: var(--mapex-mono-font);
	}
}
</style>
