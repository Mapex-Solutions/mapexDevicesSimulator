import type { AxiosResponse, InternalAxiosRequestConfig } from 'axios';

/**
 * Optional axios interceptors. The simulator is a local trusted tool, so there is
 * no auth or token refresh; interceptors exist for logging/diagnostics only.
 */
export interface ApiInterceptors {
	onRequest?: (config: InternalAxiosRequestConfig) => InternalAxiosRequestConfig | Promise<InternalAxiosRequestConfig>;
	onResponse?: (response: AxiosResponse) => AxiosResponse | Promise<AxiosResponse>;
	onError?: (error: unknown) => unknown;
}

/** Configuration for the sidecar http client. */
export interface ApiConfig {
	baseURL: string;
	headers?: Record<string, string>;
	interceptors?: ApiInterceptors;
}
