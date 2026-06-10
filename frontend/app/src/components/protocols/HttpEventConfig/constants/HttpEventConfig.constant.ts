/** TYPE IMPORTS */
import type { HttpEventConfig } from '@services/sim';

/**
 * A blank HTTP event seeded with a templated body so the placeholders are
 * discoverable.
 */
export function defaultHttpEvent(): HttpEventConfig {
	return {
		method: 'POST',
		path: '/',
		headers: [],
		bodyMode: 'raw',
		bodyFields: [],
		body: '{\n  "value": {{randInt(10,30)}}\n}',
	};
}
