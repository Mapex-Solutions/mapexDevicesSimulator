/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPOSABLES */
import { useTranslations } from '@composables/i18n';

/**
 * Translations for the ListHeaderMenu, shaped to match the component's
 * expectations (computed refs + relative-time formatter functions).
 */
export function useListHeaderMenuTranslations() {
	const { t } = useTranslations();

	return {
		refresh: computed(() => t('listHeader.refresh')),
		itemsPerPage: computed(() => t('listHeader.itemsPerPage')),
		visibleColumns: computed(() => t('listHeader.visibleColumns')),
		filtered: computed(() => t('listHeader.filtered')),
		lastUpdatedNow: computed(() => t('listHeader.now')),
		lastUpdatedSeconds: (n: number) => t('listHeader.secondsAgo', { n }),
		lastUpdatedMinutes: (n: number) => t('listHeader.minutesAgo', { n }),
		lastUpdatedHours: (n: number) => t('listHeader.hoursAgo', { n }),
	};
}
