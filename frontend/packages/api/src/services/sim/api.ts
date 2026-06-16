/** TYPE IMPORTS */
import type { ApiConfig } from '../../common';

/** TOOLS */
import { createHttp } from '../../tools/http/http.tool';

/** CLIENT */
import { resolveApiBase, resolveMarketplaceBase } from '../../client';

/** RESOURCES */
import { devicesApi } from './devices/devices.api';
import { gatewaysApi } from './gateways/gateways.api';
import { logsApi } from './logs/logs.api';
import { healthApi } from './health/health.api';
import { marketplaceApi } from './marketplace/marketplace.api';

/**
 * Builds the aggregated, typed sidecar client. Resolves the API base from the
 * preload bridge / page origin unless one is supplied, then exposes one client
 * per resource over a shared axios instance. Stores and components import the
 * returned object as their single entry point.
 *
 * @param {Partial<ApiConfig>} [config] - optional overrides (base URL, headers, interceptors)
 */
export function createSimApi(config?: Partial<ApiConfig>) {
	const http = createHttp({ baseURL: config?.baseURL ?? resolveApiBase(), ...config });

	// The marketplace is an independent online catalog, so it gets its own transport
	// pointed at the CDN base rather than the sidecar API.
	const marketplaceHttp = createHttp({ baseURL: resolveMarketplaceBase() });

	return {
		http,
		health: healthApi(http),
		devices: devicesApi(http),
		gateways: gatewaysApi(http),
		logs: logsApi(http),
		marketplace: marketplaceApi(marketplaceHttp),
	};
}

export type SimApi = ReturnType<typeof createSimApi>;
