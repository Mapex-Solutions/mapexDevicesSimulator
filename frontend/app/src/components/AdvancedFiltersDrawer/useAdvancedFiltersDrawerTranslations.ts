/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPOSABLES */
import { useTranslations } from '@composables/i18n';

/**
 * Translations for the AdvancedFiltersDrawer, shaped to match the component's
 * expectations (computed refs) while sourcing strings from the app i18n.
 */
export function useAdvancedFiltersDrawerTranslations() {
	const { t } = useTranslations();

	return {
		title: computed(() => t('advFilters.title')),
		closeTooltip: computed(() => t('advFilters.closeTooltip')),
		buttons: {
			reset: computed(() => t('advFilters.reset')),
			resetTooltip: computed(() => t('advFilters.resetTooltip')),
			apply: computed(() => t('advFilters.apply')),
			applyTooltip: computed(() => t('advFilters.applyTooltip')),
		},
		autocomplete: {
			noOption: computed(() => t('advFilters.noOption')),
		},
	};
}
