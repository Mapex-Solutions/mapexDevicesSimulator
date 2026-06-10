/**
 * Context available to template placeholders.
 */
export interface TemplateContext {
	deviceId?: string | undefined;
	deviceName?: string | undefined;
}

let counter = 0;

function randomUuid(): string {
	const c = globalThis.crypto;
	if (c && typeof c.randomUUID === 'function') return c.randomUUID();
	return `id-${Math.floor(Math.random() * 1e9)}`;
}

/**
 * Resolve a single placeholder expression to a string.
 * @param {string} expr - the inner expression (without braces)
 * @param {TemplateContext} ctx - the render context
 */
function resolve(expr: string, ctx: TemplateContext): string {
	if (expr === 'now') return new Date().toISOString();
	if (expr === 'nowMs') return String(Date.now());
	if (expr === 'deviceId') return ctx.deviceId ?? '';
	if (expr === 'deviceName') return ctx.deviceName ?? '';
	if (expr === 'counter') {
		counter += 1;
		return String(counter);
	}
	if (expr === 'uuid') return randomUuid();

	const intMatch = expr.match(/^randInt\(\s*(-?\d+)\s*,\s*(-?\d+)\s*\)$/);
	if (intMatch) {
		const a = Number(intMatch[1]);
		const b = Number(intMatch[2]);
		return String(Math.floor(a + Math.random() * (b - a + 1)));
	}

	const floatMatch = expr.match(/^randFloat\(\s*(-?\d*\.?\d+)\s*,\s*(-?\d*\.?\d+)\s*\)$/);
	if (floatMatch) {
		const a = Number(floatMatch[1]);
		const b = Number(floatMatch[2]);
		return (a + Math.random() * (b - a)).toFixed(2);
	}

	return '';
}

/**
 * Render a template string, replacing {{ ... }} placeholders. Unknown
 * placeholders resolve to an empty string. A plain string with no placeholders
 * is returned unchanged (fixed value).
 * @param {string} text - the template text
 * @param {TemplateContext} ctx - the render context
 * @returns {string} the rendered text
 */
export function renderTemplate(text: string, ctx: TemplateContext = {}): string {
	return text.replace(/\{\{\s*([^}]+?)\s*\}\}/g, (_match, expr: string) => resolve(expr.trim(), ctx));
}

/** A documented template placeholder shown in the hints UI. */
export interface TemplatePlaceholder {
	/** The mustache token, ready to paste. */
	token: string;
	/** i18n key suffix under `templateHints.desc.*` describing the token. */
	descriptionKey: string;
}

/** Placeholders offered to the user when authoring payloads and event bodies. */
export const TEMPLATE_PLACEHOLDERS: TemplatePlaceholder[] = [
	{ token: '{{randInt(10,30)}}', descriptionKey: 'randInt' },
	{ token: '{{randFloat(0,1)}}', descriptionKey: 'randFloat' },
	{ token: '{{now}}', descriptionKey: 'now' },
	{ token: '{{counter}}', descriptionKey: 'counter' },
	{ token: '{{deviceId}}', descriptionKey: 'deviceId' },
	{ token: '{{uuid}}', descriptionKey: 'uuid' },
];

/** Just the token strings, for compact inline rendering. */
export const TEMPLATE_HINTS: string[] = TEMPLATE_PLACEHOLDERS.map((placeholder) => placeholder.token);
