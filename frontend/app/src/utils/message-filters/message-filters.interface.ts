export type FilterFieldType = 'text' | 'select';

/** Where the field value is read from on a console message. */
export type FilterSource = 'field' | 'meta' | 'any';

export interface FilterFieldOption {
	value: string;
	/** i18n key for the option label (takes precedence over `label`). */
	labelKey?: string;
	/** Literal label when no translation key applies. */
	label?: string;
}

/**
 * One filter field. Protocol-specific fields read from `meta`; common fields
 * read from a message field or match across all text (`any`).
 */
export interface FilterField {
	key: string;
	labelKey: string;
	type: FilterFieldType;
	source: FilterSource;
	options?: FilterFieldOption[];
}

/** Active filter values keyed by field key. Empty string means inactive. */
export type FilterValues = Record<string, string>;
