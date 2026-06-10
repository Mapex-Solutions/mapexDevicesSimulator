/** TYPE IMPORTS */
import type { AxiosInstance } from 'axios';
import type { Gateway, GatewayInput } from '@sim/schema';

/** SCHEMAS */
import { ZodGatewayInputSchema } from '@sim/schema';

/** HANDLERS */
import { createApiFactory } from '../../../common';

/**
 * Gateways resource: full CRUD against /api/gateways (LoRaWAN gateways only).
 *
 * @param {AxiosInstance} http - the shared sidecar axios instance
 */
export function gatewaysApi(http: AxiosInstance) {
	const factory = createApiFactory(http);

	return factory({
		basePath: '/gateways',
		methods: {
			list: {
				method: 'GET',
				path: '',
				responseType: {} as Gateway[],
			},
			create: {
				method: 'POST',
				path: '',
				bodyParams: {} as GatewayInput,
				bodySchema: ZodGatewayInputSchema,
				responseType: {} as Gateway,
			},
			update: {
				method: 'PUT',
				path: '/:id',
				pathParams: {} as { id: string },
				bodyParams: {} as GatewayInput,
				bodySchema: ZodGatewayInputSchema,
				responseType: {} as Gateway,
			},
			remove: {
				method: 'DELETE',
				path: '/:id',
				pathParams: {} as { id: string },
				responseType: undefined as void,
			},
		},
	});
}
