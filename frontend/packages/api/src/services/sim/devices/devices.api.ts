/** TYPE IMPORTS */
import type { AxiosInstance } from 'axios';
import type { Device, DeviceInput, DeviceEvent } from '@sim/schema';

/** SCHEMAS */
import { ZodDeviceInputSchema } from '@sim/schema';

/** HANDLERS */
import { createApiFactory } from '../../../common';

/**
 * Devices resource: full CRUD against /api/devices. Create and update validate the
 * body against the device input schema before it leaves.
 *
 * @param {AxiosInstance} http - the shared sidecar axios instance
 */
export function devicesApi(http: AxiosInstance) {
	const factory = createApiFactory(http);

	return factory({
		basePath: '/devices',
		methods: {
			list: {
				method: 'GET',
				path: '',
				responseType: {} as Device[],
			},
			create: {
				method: 'POST',
				path: '',
				bodyParams: {} as DeviceInput,
				bodySchema: ZodDeviceInputSchema,
				responseType: {} as Device,
			},
			update: {
				method: 'PUT',
				path: '/:id',
				pathParams: {} as { id: string },
				bodyParams: {} as DeviceInput,
				bodySchema: ZodDeviceInputSchema,
				responseType: {} as Device,
			},
			remove: {
				method: 'DELETE',
				path: '/:id',
				pathParams: {} as { id: string },
				responseType: undefined as void,
			},
			fire: {
				method: 'POST',
				path: '/:id/fire',
				pathParams: {} as { id: string },
				bodyParams: {} as { eventId?: string; event?: DeviceEvent },
				responseType: {} as { fired: boolean },
			},
		},
	});
}
