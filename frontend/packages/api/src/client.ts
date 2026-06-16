/**
 * The Electron preload bridge exposes the dynamically-chosen sidecar origin so a
 * cross-origin dev SPA can reach it; in the packaged app the page origin already
 * points at the sidecar. We read it through a local cast rather than augmenting
 * the global Window, so this package never collides with the app's own typing of
 * `window.__SIM__`.
 */
interface SidecarBridge {
	apiBase?: string;
	wsBase?: string;
	marketplaceBase?: string;
}

/** Read the preload-bridged sidecar origins, if running inside the desktop shell. */
function bridge(): SidecarBridge | undefined {
	if (typeof window === 'undefined') return undefined;
	return (window as unknown as { __SIM__?: SidecarBridge }).__SIM__;
}

/**
 * Resolve the sidecar control API base. Prefers the preload-bridged origin, then
 * the page origin, then a localhost fallback for non-browser contexts.
 */
export function resolveApiBase(): string {
	const origin = bridge()?.apiBase ?? (typeof window !== 'undefined' ? window.location.origin : 'http://127.0.0.1:5055');
	return `${origin.replace(/\/$/, '')}/api`;
}

/**
 * Resolve the device marketplace base — the online mapexMarketplace service,
 * independent of the local sidecar. Prefers the preload-bridged origin so the
 * packaged app can be pointed at the deployed catalog; in development it falls
 * back to the local marketplace server on :6060.
 */
export function resolveMarketplaceBase(): string {
	const bridged = bridge()?.marketplaceBase;
	if (bridged) return bridged.replace(/\/$/, '');
	return 'http://127.0.0.1:6060/api/v1';
}

/**
 * Build the URL of a marketplace bundle asset (device image, codec file) for a
 * model, served by the catalog under `/devices/:vendor/:slug/assets/*`.
 */
export function resolveMarketplaceAssetUrl(vendor: string, slug: string, path: string): string {
	const clean = path.replace(/^\//, '');
	return `${resolveMarketplaceBase()}/devices/${vendor}/${slug}/assets/${clean}`;
}

/** Resolve the sidecar WebSocket base used by the live console stream. */
export function resolveWsBase(): string {
	const bridged = bridge()?.wsBase;
	if (bridged) return bridged.replace(/\/$/, '');

	const origin = typeof window !== 'undefined' ? window.location.origin : 'http://127.0.0.1:5055';
	return origin.replace(/^http/, 'ws').replace(/\/$/, '');
}
