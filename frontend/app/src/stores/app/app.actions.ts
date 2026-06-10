/** TYPE IMPORTS */
import type { AppStateContext } from './app.types';

/** SERVICES */
import { sim } from '@services/sim';

const HEALTH_INTERVAL_MS = 5000;

/**
 * Sidecar health polling. The poll timer lives in this closure because it is
 * not part of reactive state.
 * @param {AppStateContext} ctx - reactive state bag
 */
export function createAppActions(ctx: AppStateContext) {
	let healthTimer: ReturnType<typeof setInterval> | null = null;

	/**
	 * Probe the sidecar once and reflect the outcome in state.
	 */
	async function checkHealth(): Promise<void> {
		try {
			const res = await sim.health.get();
			ctx.sidecarStatus.value = res.status === 'ok' ? 'connected' : 'disconnected';
			ctx.version.value = res.version;
		} catch {
			ctx.sidecarStatus.value = 'disconnected';
			ctx.version.value = '';
		}
	}

	/**
	 * Run an immediate health check and keep polling on an interval.
	 */
	async function startHealthPolling(): Promise<void> {
		await checkHealth();
		if (healthTimer) return;
		healthTimer = setInterval(() => {
			void checkHealth();
		}, HEALTH_INTERVAL_MS);
	}

	function stopHealthPolling(): void {
		if (!healthTimer) return;
		clearInterval(healthTimer);
		healthTimer = null;
	}

	return { checkHealth, startHealthPolling, stopHealthPolling };
}
