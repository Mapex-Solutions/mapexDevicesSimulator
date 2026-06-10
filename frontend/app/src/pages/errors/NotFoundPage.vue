<template>
	<q-page class="column flex-center q-pa-xl">
		<q-icon name="mdi-map-marker-question-outline" size="72px" color="primary" />
		<h1 class="not-found__title">{{ t('notFound.title') }}</h1>
		<p class="not-found__message">{{ t('notFound.message') }}</p>
		<q-btn
			color="primary"
			icon="mdi-arrow-left"
			:label="t('notFound.back')"
			:to="{ name: 'dashboard' }"
		/>
	</q-page>
</template>

<script setup lang="ts">
/** COMPOSABLES */
import { useTranslations } from '@composables/i18n';

/** COMPOSABLES & STORES */
const { t } = useTranslations();
</script>

<style scoped lang="scss">
.not-found {
	&__title {
		font-size: 1.5rem;
		font-weight: var(--mapex-font-weight-semibold);
		color: var(--mapex-text-primary);
		margin: var(--mapex-spacing-lg) 0 var(--mapex-spacing-sm);
	}

	&__message {
		color: var(--mapex-text-secondary);
		margin-bottom: var(--mapex-spacing-lg);
	}
}
</style>
