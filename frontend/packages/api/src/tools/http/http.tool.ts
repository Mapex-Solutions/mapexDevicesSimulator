/** TYPE IMPORTS */
import type { AxiosInstance, CreateAxiosDefaults } from 'axios';
import type { ApiConfig, ApiInterceptors } from '../../common/interfaces';

/** SERVICES */
import axios from 'axios';

/** UTILS */
import { isEmpty } from '../../common/utils/isEmpty';

/**
 * Creates the axios instance for the sidecar control API: a typed transport with
 * the JSON content type and any optional interceptors wired. The envelope unwrap
 * lives in the api factory, not here, so this stays a plain transport.
 *
 * @param {ApiConfig} config - base URL, headers and optional interceptors
 * @param {ApiInterceptors} [interceptors] - interceptors overriding the config's
 */
export function createHttp(config: ApiConfig, interceptors?: ApiInterceptors): AxiosInstance {
	const { baseURL, headers = {} } = config;
	const active = !isEmpty(interceptors) ? interceptors : config.interceptors;

	const params: CreateAxiosDefaults = {
		baseURL,
		timeout: 15000,
		headers: { 'Content-Type': 'application/json', ...headers },
	};

	const instance = axios.create(params);

	if (active?.onRequest) instance.interceptors.request.use(active.onRequest);
	if (active?.onResponse) instance.interceptors.response.use(active.onResponse);
	if (active?.onError) instance.interceptors.response.use(null, active.onError);

	return instance;
}
