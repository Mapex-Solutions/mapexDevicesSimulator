/** TYPE IMPORTS */
import type { AxiosInstance } from 'axios';
import type { LogPage, LogQuery } from '@sim/schema';

/** SCHEMAS */
import { ZodLogQuerySchema } from '@sim/schema';

/** HANDLERS */
import { createApiFactory } from '../../../common';

/**
 * Logs resource: the paginated, filterable history behind the console stream.
 *
 * @param {AxiosInstance} http - the shared sidecar axios instance
 */
export function logsApi(http: AxiosInstance) {
	const factory = createApiFactory(http);

	return factory({
		basePath: '/logs',
		methods: {
			list: {
				method: 'GET',
				path: '',
				queryParams: {} as LogQuery,
				querySchema: ZodLogQuerySchema,
				responseType: {} as LogPage,
			},
		},
	});
}
