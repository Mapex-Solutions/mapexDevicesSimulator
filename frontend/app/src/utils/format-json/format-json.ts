/**
 * Pretty-print a value when it is a JSON object or array, otherwise return it
 * unchanged. Used to make console payloads and HTTP responses readable without
 * mangling non-JSON content (a LoRaWAN hex payload, a plain string, a bare code).
 *
 * Only object/array bodies are reformatted: those are the JSON worth indenting,
 * and the guard keeps a hex string that happens to parse as a number untouched.
 * @param {string} text - the raw payload or response text
 * @returns {string} the indented JSON, or the original text when it is not JSON
 */
export function formatJson(text: string): string {
	const trimmed = text.trim();
	if (trimmed === '') return text;

	const first = trimmed[0];
	if (first !== '{' && first !== '[') return text;

	try {
		return JSON.stringify(JSON.parse(trimmed), null, 2);
	} catch {
		return text;
	}
}
