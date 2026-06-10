/** TYPE IMPORTS */
import type { RequestBody } from '@services/sim';
import type { TemplateContext } from './render-template';

/** UTILS */
import { renderTemplate } from './render-template';

/**
 * Coerce a rendered string to a JSON-friendly value (number/boolean/string).
 * @param {string} value - the rendered value
 * @returns {unknown} the coerced value
 */
function coerce(value: string): unknown {
	if (value === 'true') return true;
	if (value === 'false') return false;
	if (value !== '' && /^-?\d+(\.\d+)?$/.test(value)) return Number(value);
	return value;
}

/**
 * Build a request body from an event's body config, rendering templates. In
 * form mode the fields are assembled into a JSON object; in raw mode the body
 * template is rendered as-is; in none mode the body is empty.
 * @param {RequestBody} config - the body config (HTTP or MQTT event)
 * @param {TemplateContext} ctx - the render context
 * @returns {string} the rendered body
 */
export function buildHttpBody(config: RequestBody, ctx: TemplateContext = {}): string {
	if (config.bodyMode === 'none') return '';

	if (config.bodyMode === 'form') {
		const obj: Record<string, unknown> = {};
		for (const field of config.bodyFields) {
			const key = field.key.trim();
			if (key) obj[key] = coerce(renderTemplate(field.value, ctx));
		}
		return JSON.stringify(obj, null, 2);
	}

	return renderTemplate(config.body, ctx);
}
