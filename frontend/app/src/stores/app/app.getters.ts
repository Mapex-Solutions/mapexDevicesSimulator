/** VUE IMPORTS */
import { computed } from 'vue';

/** TYPE IMPORTS */
import type { AppStateContext } from './app.types';

/**
 * Derived state for the sidecar connection.
 * @param {AppStateContext} ctx - reactive state bag
 */
export function createAppGetters(ctx: AppStateContext) {
	const isConnected = computed(() => ctx.sidecarStatus.value === 'connected');
	const isChecking = computed(() => ctx.sidecarStatus.value === 'checking');

	return { isConnected, isChecking };
}
