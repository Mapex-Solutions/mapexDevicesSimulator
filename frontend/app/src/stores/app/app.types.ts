/** TYPE IMPORTS */
import type { Ref } from 'vue';

export type SidecarStatus = 'checking' | 'connected' | 'disconnected';

/**
 * Reactive state bag shared between the store definition and its getter/action
 * factories.
 */
export interface AppStateContext {
	sidecarStatus: Ref<SidecarStatus>;
	version: Ref<string>;
}
