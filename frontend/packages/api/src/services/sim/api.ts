/** TYPE IMPORTS */
import type { ApiConfig } from '../../common';

/** TOOLS */
import { createHttp } from '../../tools/http/http.tool';

/** CLIENT */
import { resolveApiBase } from '../../client';

/** RESOURCES */
import { devicesApi } from './devices/devices.api';
import { gatewaysApi } from './gateways/gateways.api';
import { logsApi } from './logs/logs.api';
import { healthApi } from './health/health.api';

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

	return {
		http,
		health: healthApi(http),
		devices: devicesApi(http),
		gateways: gatewaysApi(http),
		logs: logsApi(http),
	};
}

export type SimApi = ReturnType<typeof createSimApi>;
