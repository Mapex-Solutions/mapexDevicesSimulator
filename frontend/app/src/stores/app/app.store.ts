/** VUE IMPORTS */
import { ref } from 'vue';

/** TYPE IMPORTS */
import type { AppStateContext, SidecarStatus } from './app.types';

/** SERVICES */
import { defineStore } from 'pinia';

/** LOCAL IMPORTS */
import { createAppGetters } from './app.getters';
import { createAppActions } from './app.actions';

/**
 * Sidecar connectivity store: tracks whether the local engine is reachable and
 * its version. Theme and navigation live in their own stores.
 */
export const useAppStore = defineStore('app', () => {
	/** STATE */
	const sidecarStatus = ref<SidecarStatus>('checking');
	const version = ref<string>('');

	const ctx: AppStateContext = { sidecarStatus, version };

	const getters = createAppGetters(ctx);
	const actions = createAppActions(ctx);

	return { sidecarStatus, version, ...getters, ...actions };
});
