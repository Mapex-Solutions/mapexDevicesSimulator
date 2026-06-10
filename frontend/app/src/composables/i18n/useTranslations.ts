/** VUE IMPORTS */
import { useI18n } from 'vue-i18n';

/**
 * Single entry point for translations across the app. Wrapping useI18n keeps a
 * stable import surface so component scripts never reach for vue-i18n directly.
 */
export function useTranslations() {
	return useI18n({ useScope: 'global' });
}
