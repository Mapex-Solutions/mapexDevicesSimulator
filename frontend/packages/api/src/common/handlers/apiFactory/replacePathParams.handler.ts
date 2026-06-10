import { SchemaError } from '../schemaError.handler';

/**
 * Replaces `:name` path parameters with their URL-encoded values. Throws when a
 * placeholder has no matching value, so a malformed call fails loudly instead of
 * hitting the wrong URL.
 */
export function replacePathParams(path: string, pathParams: Record<string, unknown>): string {
	return path.replace(/:([a-zA-Z0-9_-]+)/g, (_, key) => {
		if (!(key in pathParams)) {
			throw new SchemaError([`missing path parameter: ${key}`]);
		}
		return encodeURIComponent(String(pathParams[key]));
	});
}
