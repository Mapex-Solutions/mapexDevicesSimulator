/** TYPE IMPORTS */
import type { AxiosInstance } from 'axios';
import type { HealthResponse } from '@sim/schema';

/** HANDLERS */
import { createApiFactory } from '../../../common';

/**
 * Health resource: the liveness/version probe the app polls to show sidecar state.
 *
 * @param {AxiosInstance} http - the shared sidecar axios instance
 */
export function healthApi(http: AxiosInstance) {
	const factory = createApiFactory(http);

	return factory({
		basePath: '/health',
		methods: {
			get: {
				method: 'GET',
				path: '',
				responseType: {} as HealthResponse,
			},
		},
	});
}
