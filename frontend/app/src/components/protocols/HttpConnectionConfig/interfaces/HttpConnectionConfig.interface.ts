/** TYPE IMPORTS */
import type { HttpConnectionConfig as HttpConfig } from '@services/sim';

export interface HttpConnectionConfigProps {
	modelValue: HttpConfig;
}

export interface HttpConnectionConfigEmits {
	(event: 'update:modelValue', value: HttpConfig): void;
}
