import type { AxiosRequestConfig, AxiosResponse } from 'axios';

/** Resolve an explicit responseType, falling back to `any` when left as `{}`. */
export type InferResponseType<R> = {} extends R ? any : R;

/**
 * A single endpoint definition. Which of pathParams / bodyParams / queryParams
 * are set determines the generated method's call signature; the matching Zod
 * schema validates that argument before the request leaves.
 */
export interface ApiMethod<P = {}, Q = {}, B = {}, R = {}> {
	pathParams?: P;
	queryParams?: Q;
	bodyParams?: B;

	path: string;
	method: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';
	headers?: Record<string, string>;
	axiosConfig?: AxiosRequestConfig;

	responseType?: R;

	afterRequest?: (response: AxiosResponse) => Promise<R> | R;

	paramSchema?: { safeParseAsync: (data: unknown) => Promise<any> };
	bodySchema?: { safeParseAsync: (data: unknown) => Promise<any> };
	querySchema?: { safeParseAsync: (data: unknown) => Promise<any> };
}

/** The configuration passed to a factory: a base path plus its endpoint map. */
export interface ApiClientConfig<T extends Record<string, ApiMethod<any, any, any, any>>> {
	basePath: string;
	methods: T;
}
