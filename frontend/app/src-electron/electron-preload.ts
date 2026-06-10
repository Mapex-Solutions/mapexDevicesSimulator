/**
 * Preload bridge. Exposes the sidecar base URLs (passed as launch arguments by
 * the main process) to the renderer under a read-only window.__SIM__ object so
 * the SPA can reach the engine even when served cross-origin by the dev server.
 */

import { contextBridge } from 'electron';

/**
 * Read a `--flag=value` launch argument.
 * @param {string} flag - the argument name including leading dashes
 */
function argValue(flag: string): string | undefined {
	const prefix = `${flag}=`;
	const match = process.argv.find((arg) => arg.startsWith(prefix));
	return match ? match.slice(prefix.length) : undefined;
}

contextBridge.exposeInMainWorld('__SIM__', {
	apiBase: argValue('--sim-api-base'),
	wsBase: argValue('--sim-ws-base'),
	appVersion: argValue('--sim-app-version'),
	platform: process.platform,
});
