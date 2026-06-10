/** TYPE IMPORTS */
import type { HttpEventConfig as HttpEvent } from '@services/sim';

export interface HttpEventConfigProps {
	modelValue: HttpEvent;
}

export interface HttpEventConfigEmits {
	(event: 'update:modelValue', value: HttpEvent): void;
}
