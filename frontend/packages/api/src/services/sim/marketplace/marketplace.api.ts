/** TYPE IMPORTS */
import type { AxiosInstance, AxiosResponse } from 'axios';
import type {
	MarketplaceListResponse,
	MarketplaceFacets,
	MarketplaceInformation,
	MarketplaceSimulator,
	MarketplaceCodec,
} from '@sim/schema';

/** HANDLERS */
import { createApiFactory } from '../../../common';

/**
 * Device marketplace resource: read-only access to the online mapexMarketplace
 * service (`/api/v1/devices`). Responses use the Mapex `{ status, errors, data }`
 * envelope, so each method unwraps `data`. The catalog is browsed here; installing
 * a chosen template into SQLite goes through the devices resource.
 *
 * @param {AxiosInstance} http - an axios instance pointed at the marketplace base
 */
export function marketplaceApi(http: AxiosInstance) {
	const factory = createApiFactory(http);

	return factory({
		basePath: '/devices',
		methods: {
			list: {
				method: 'GET',
				path: '',
				queryParams: {} as {
					protocol?: string;
					readingType?: string;
					manufacturer?: string;
					search?: string;
					page?: number;
					perPage?: number;
				},
				responseType: {} as MarketplaceListResponse,
				afterRequest: (res: AxiosResponse) => res.data.data as MarketplaceListResponse,
			},
			facets: {
				method: 'GET',
				path: '/facets',
				responseType: {} as MarketplaceFacets,
				afterRequest: (res: AxiosResponse) => res.data.data as MarketplaceFacets,
			},
			information: {
				method: 'GET',
				path: '/:vendor/:slug',
				pathParams: {} as { vendor: string; slug: string },
				responseType: {} as MarketplaceInformation,
				afterRequest: (res: AxiosResponse) => res.data.data as MarketplaceInformation,
			},
			simulator: {
				method: 'GET',
				path: '/:vendor/:slug/simulator',
				pathParams: {} as { vendor: string; slug: string },
				responseType: {} as MarketplaceSimulator,
				afterRequest: (res: AxiosResponse) => res.data.data as MarketplaceSimulator,
			},
			codecs: {
				method: 'GET',
				path: '/:vendor/:slug/codecs',
				pathParams: {} as { vendor: string; slug: string },
				responseType: {} as MarketplaceCodec[],
				afterRequest: (res: AxiosResponse) => res.data.data as MarketplaceCodec[],
			},
		},
	});
}
