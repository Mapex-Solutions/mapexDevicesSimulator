<template>
	<q-expansion-item
		class="template-hints"
		dense
		dense-toggle
		icon="mdi-code-braces"
		:label="t('templateHints.title')"
		:caption="t('templateHints.expand')"
		header-class="template-hints__header"
	>
		<div class="template-hints__body">
			<div class="template-hints__note">
				<q-icon name="mdi-cursor-default-click-outline" size="14px" />
				<span>{{ t('templateHints.hint') }}</span>
			</div>
			<div class="template-hints__grid">
				<button
					v-for="placeholder in TEMPLATE_PLACEHOLDERS"
					:key="placeholder.token"
					type="button"
					class="template-hints__item"
					@click="copy(placeholder.token)"
				>
					<div class="template-hints__item-top">
						<code class="template-hints__token">{{ placeholder.token }}</code>
						<q-icon name="mdi-content-copy" size="14px" class="template-hints__copy" />
					</div>
					<span class="template-hints__desc">{{ t(`templateHints.desc.${placeholder.descriptionKey}`) }}</span>
				</button>
			</div>
		</div>
	</q-expansion-item>
</template>

<script setup lang="ts">
/** COMPOSABLES */
import { useTranslations } from '@composables/i18n';

/** UTILS */
import { notifySuccess, notifyWarning } from '@utils/alert';
import { TEMPLATE_PLACEHOLDERS } from '@utils/template';

/** COMPOSABLES & STORES */
const { t } = useTranslations();

/** FUNCTIONS */

/**
 * Copy a placeholder token to the clipboard and confirm with a toast.
 * @param {string} token - the placeholder token to copy
 */
async function copy(token: string): Promise<void> {
	try {
		await navigator.clipboard.writeText(token);
		notifySuccess({ message: t('templateHints.copied', { token }), timeout: 1200 });
	} catch {
		notifyWarning({ message: t('templateHints.copyFailed') });
	}
}
</script>

<style scoped lang="scss">
.template-hints {
	flex: 0 0 auto;
	border: 1px solid var(--mapex-card-border);
	border-radius: var(--mapex-radius-md);
	background: var(--mapex-surface-sunken);
	overflow: hidden;

	:deep(.template-hints__header) {
		min-height: 0;
		padding: var(--mapex-spacing-xs) var(--mapex-spacing-md);

		.q-item__label {
			font-size: var(--mapex-font-xs);
			font-weight: var(--mapex-font-weight-medium);
			color: var(--mapex-text-secondary);
			text-transform: uppercase;
			letter-spacing: 0.04em;
		}

		.q-item__label--caption {
			text-transform: none;
			letter-spacing: 0;
			color: var(--mapex-text-muted);
		}

		.q-icon {
			color: var(--mapex-primary);
		}
	}

	&__body {
		display: flex;
		flex-direction: column;
		gap: var(--mapex-spacing-sm);
		padding: var(--mapex-spacing-sm) var(--mapex-spacing-md) var(--mapex-spacing-md);
	}

	&__note {
		display: flex;
		align-items: center;
		gap: var(--mapex-spacing-2xs);
		font-size: var(--mapex-font-xs);
		color: var(--mapex-text-muted);
	}

	&__grid {
		/* Fixed column count so the panel height stays constant after it
		   opens. An auto-fill grid reflows when a scrollbar appears, which
		   made the expansion jump taller a beat after expanding. */
		display: grid;
		grid-template-columns: repeat(3, minmax(0, 1fr));
		gap: var(--mapex-spacing-sm);
	}

	@media (max-width: 600px) {
		&__grid {
			grid-template-columns: repeat(2, minmax(0, 1fr));
		}
	}

	&__item {
		display: flex;
		flex-direction: column;
		gap: var(--mapex-spacing-2xs);
		padding: var(--mapex-spacing-sm) var(--mapex-spacing-md);
		text-align: left;
		border: 1px solid var(--mapex-card-border);
		border-radius: var(--mapex-radius-sm);
		background: var(--mapex-surface-elevated);
		cursor: pointer;
		transition: var(--mapex-transition-base);

		&:hover {
			border-color: var(--mapex-primary);
			box-shadow: var(--mapex-shadow-sm);

			.template-hints__copy {
				opacity: 1;
			}
		}
	}

	&__item-top {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: var(--mapex-spacing-xs);
	}

	&__token {
		font-family: var(--mapex-mono-font);
		font-size: var(--mapex-font-sm);
		color: var(--mapex-primary);
	}

	&__copy {
		flex-shrink: 0;
		color: var(--mapex-text-muted);
		opacity: 0;
		transition: var(--mapex-transition-base);
	}

	&__desc {
		font-size: var(--mapex-font-xs);
		color: var(--mapex-text-muted);
		line-height: 1.3;
	}
}
</style>
