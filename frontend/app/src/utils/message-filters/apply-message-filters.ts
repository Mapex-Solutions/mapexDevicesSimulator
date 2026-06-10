/** TYPE IMPORTS */
import type { ConsoleMessage } from '@stores/messages';
import type { FilterField, FilterValues } from './message-filters.interface';

/**
 * Read a single message field as a string for matching.
 * @param {ConsoleMessage} message - the message
 * @param {string} key - the field key
 */
function readField(message: ConsoleMessage, key: string): string {
	const value = (message as unknown as Record<string, unknown>)[key];
	return value === undefined || value === null ? '' : String(value);
}

/**
 * Whether a message satisfies one active field filter.
 * @param {ConsoleMessage} message - the message
 * @param {FilterField} field - the field definition
 * @param {string} raw - the active value
 */
function matches(message: ConsoleMessage, field: FilterField, raw: string): boolean {
	const value = raw.trim().toLowerCase();

	if (field.source === 'any') {
		return `${message.summary} ${message.payload} ${message.deviceName}`.toLowerCase().includes(value);
	}

	const target = (
		field.source === 'meta' ? (message.meta?.[field.key] ?? '') : readField(message, field.key)
	).toLowerCase();

	return field.type === 'text' ? target.includes(value) : target === value;
}

/**
 * Filter a message list by the active values for the given fields. Fields with
 * an empty value are ignored; a message must satisfy every active field.
 * @param {ConsoleMessage[]} messages - the messages to filter
 * @param {FilterField[]} fields - the fields in effect
 * @param {FilterValues} values - the active values keyed by field key
 * @returns {ConsoleMessage[]} the filtered messages
 */
export function applyMessageFilters(
	messages: ConsoleMessage[],
	fields: FilterField[],
	values: FilterValues,
): ConsoleMessage[] {
	const active = fields.filter((field) => (values[field.key] ?? '').trim() !== '');
	if (!active.length) return messages;

	return messages.filter((message) => active.every((field) => matches(message, field, values[field.key] ?? '')));
}
