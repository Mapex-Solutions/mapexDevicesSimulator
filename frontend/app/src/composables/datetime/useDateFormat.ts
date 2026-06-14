/** COMPOSABLES */
import { useTranslations } from '@composables/i18n';

/**
 * Date and time formatting bound to the app's selected language, so timestamps
 * follow the chosen locale (and switch with it) rather than the browser's. One
 * source of truth keeps every list and detail view on the same format.
 */
export function useDateFormat() {
	const { locale } = useTranslations();

	/**
	 * Format an ISO timestamp as a short localized date.
	 * @param {string | null | undefined} iso - the ISO timestamp
	 * @returns {string} the localized date, or an em dash when absent
	 */
	function formatDate(iso?: string | null): string {
		if (!iso) return '—';
		const date = new Date(iso);
		return Number.isNaN(date.getTime()) ? '—' : date.toLocaleDateString(locale.value, { dateStyle: 'short' });
	}

	/**
	 * Format an ISO timestamp as a short localized date and time.
	 * @param {string | null | undefined} iso - the ISO timestamp
	 * @returns {string} the localized date and time, or an em dash when absent
	 */
	function formatDateTime(iso?: string | null): string {
		if (!iso) return '—';
		const date = new Date(iso);
		return Number.isNaN(date.getTime()) ? '—' : date.toLocaleString(locale.value, { dateStyle: 'short', timeStyle: 'short' });
	}

	return { formatDate, formatDateTime };
}
