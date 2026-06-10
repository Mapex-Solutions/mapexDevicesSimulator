/* eslint-disable */

declare namespace NodeJS {
	interface ProcessEnv {
		NODE_ENV: string;
		VUE_ROUTER_MODE: 'hash' | 'history' | 'abstract' | undefined;
		VUE_ROUTER_BASE: string | undefined;
	}
}

/**
 * Bridge exposed by the Electron preload script. Carries the Go sidecar base
 * URLs so the renderer can reach the control API and the live stream even when
 * the SPA is served cross-origin by the Vite dev server.
 */
interface SimulatorBridge {
	apiBase?: string;
	wsBase?: string;
	platform?: string;
	appVersion?: string;
}

interface Window {
	__SIM__?: SimulatorBridge;
}
