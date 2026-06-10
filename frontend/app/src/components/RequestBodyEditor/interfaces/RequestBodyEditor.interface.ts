/** TYPE IMPORTS */
import type { RequestBody } from '@services/sim';

/** Props for the request body editor. */
export interface RequestBodyEditorProps {
	modelValue: RequestBody;
}

/** Events emitted by the request body editor. */
export interface RequestBodyEditorEmits {
	(event: 'update:modelValue', value: RequestBody): void;
}
