import type { AxiosInstance } from 'axios';
import type { ApiMethod, ApiClientConfig, InferResponseType } from '../../interfaces';
import { SchemaError, zodValidationError } from '../schemaError.handler';
import { replacePathParams } from './replacePathParams.handler';

/**
 * Derives each generated method's call signature from which of pathParams /
 * bodyParams / queryParams the definition sets (an omitted one infers as
 * `unknown`, which is excluded by `extends {}`), and its return from responseType.
 */
type ApiClient<T extends Record<string, ApiMethod<any, any, any, any>>> = {
	[K in keyof T]: T[K] extends ApiMethod<infer P, infer Q, infer B, infer R>
		? P extends {}
			? B extends {}
				? Q extends {}
					? (params: P, body: B, query: Q) => Promise<InferResponseType<R>>
					: (params: P, body: B) => Promise<InferResponseType<R>>
				: Q extends {}
					? (params: P, query: Q) => Promise<InferResponseType<R>>
					: (params: P) => Promise<InferResponseType<R>>
			: B extends {}
				? Q extends {}
					? (body: B, query: Q) => Promise<InferResponseType<R>>
					: (body: B) => Promise<InferResponseType<R>>
				: Q extends {}
					? (query: Q) => Promise<InferResponseType<R>>
					: () => Promise<InferResponseType<R>>
		: never;
};

/**
 * Builds typed API clients from a base path and an endpoint map. Each method
 * validates its argument against the matching Zod schema, substitutes path params,
 * serializes the query, then unwraps the Mapex `{ status, errors, data }` envelope
 * so callers receive `data` directly.
 *
 * @param http - the axios instance pointed at the sidecar control API
 */
export function createApiFactory(http: AxiosInstance) {
	return function createApiClient<T extends Record<string, ApiMethod<any, any, any, any>>>(
		config: ApiClientConfig<T>,
	): ApiClient<T> {
		const client = {} as ApiClient<T>;

		for (const [methodName, methodConfig] of Object.entries(config.methods)) {
			const fullPath = config.basePath + methodConfig.path;

			client[methodName as keyof T] = (async (...args: any[]) => {
				let url = fullPath;
				let params: any;
				let body: any;
				let query: any;
				let headers: Record<string, string> = { 'Content-Type': 'application/json' };

				if (methodConfig?.pathParams) params = args.shift();
				if (methodConfig?.bodyParams) body = args.shift();
				if (methodConfig?.queryParams) query = args.shift();

				if (methodConfig?.paramSchema) {
					const parsed = await methodConfig.paramSchema.safeParseAsync(params);
					if (!parsed.success) throw new SchemaError(zodValidationError(parsed));
					params = parsed.data;
				}

				if (methodConfig?.bodySchema) {
					const parsed = await methodConfig.bodySchema.safeParseAsync(body);
					if (!parsed.success) throw new SchemaError(zodValidationError(parsed));
					body = parsed.data;
				}

				if (methodConfig?.querySchema) {
					const parsed = await methodConfig.querySchema.safeParseAsync(query);
					if (!parsed.success) throw new SchemaError(zodValidationError(parsed));
					query = parsed.data;
				}

				if (params) url = replacePathParams(url, params);

				/**
				 * Serialize the query, skipping null/undefined and emitting arrays as
				 * repeated keys to match fiber's slice binding.
				 */
				if (query) {
					const search = new URLSearchParams();
					for (const [key, value] of Object.entries(query as Record<string, unknown>)) {
						if (value === undefined || value === null || value === '') continue;
						if (Array.isArray(value)) {
							for (const item of value) {
								if (item === undefined || item === null) continue;
								search.append(key, String(item));
							}
						} else {
							search.append(key, String(value));
						}
					}
					const queryString = search.toString();
					if (queryString) url += `?${queryString}`;
				}

				if (methodConfig?.headers) headers = { ...headers, ...methodConfig.headers };

				const restData = await http.request({
					url,
					method: methodConfig.method,
					data: body,
					headers,
					...methodConfig.axiosConfig,
				});

				if (methodConfig?.afterRequest) return methodConfig.afterRequest(restData);
				return restData?.data?.data;
			}) as ApiClient<T>[keyof T];
		}

		return client;
	};
}
