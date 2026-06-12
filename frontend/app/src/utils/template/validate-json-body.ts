/** Result of validating a templated JSON body. */
export interface JsonBodyValidation {
	/** Whether the body parses as JSON once placeholders are substituted. */
	valid: boolean;
	/** The raw parser message when invalid (empty when valid). */
	detail: string;
}

/**
 * Validate a raw event body as JSON, tolerating template placeholders.
 *
 * The body the user authors is not literal JSON: it carries {{ ... }} tokens
 * (e.g. {{deviceId}}, {{randInt(1,9)}}) resolved by the engine at send time.
 * To check the surrounding structure we replace every token with a neutral
 * literal (0) — valid in both string ("0") and numeric positions — then parse.
 * This catches the real authoring mistakes (missing comma, trailing comma,
 * unbalanced braces, unquoted keys) without faulting a correctly templated body.
 *
 * An empty body is treated as valid; emptiness is a separate concern from
 * malformed JSON.
 * @param {string} raw - the raw body text the user typed
 * @returns {JsonBodyValidation} the validation outcome
 */
export function validateJsonBody(raw: string): JsonBodyValidation {
	if (raw.trim() === '') return { valid: true, detail: '' };

	const substituted = raw.replace(/\{\{\s*[^}]+?\s*\}\}/g, '0');

	try {
		JSON.parse(substituted);
		return { valid: true, detail: '' };
	} catch (error) {
		return { valid: false, detail: error instanceof Error ? error.message : String(error) };
	}
}
