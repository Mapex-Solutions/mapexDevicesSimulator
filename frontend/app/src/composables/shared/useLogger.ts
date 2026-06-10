/**
 * Minimal namespaced logger so renderer logs carry a consistent prefix and can
 * be silenced in one place later if needed.
 * @param {string} scope - short component or domain name
 */
export function useLogger(scope: string) {
	const prefix = `[${scope}]`;

	return {
		debug: (...args: unknown[]): void => console.debug(prefix, ...args),
		info: (...args: unknown[]): void => console.info(prefix, ...args),
		warn: (...args: unknown[]): void => console.warn(prefix, ...args),
		error: (...args: unknown[]): void => console.error(prefix, ...args),
	};
}
