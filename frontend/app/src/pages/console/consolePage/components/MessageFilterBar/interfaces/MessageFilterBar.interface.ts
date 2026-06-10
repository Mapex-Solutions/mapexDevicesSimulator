/** TYPE IMPORTS */
import type { ProtocolId } from '@services/sim';
import type { FilterValues } from '@utils/message-filters';

export interface MessageFilterBarProps {
	protocol: ProtocolId | null;
	modelValue: FilterValues;
}

export interface MessageFilterBarEmits {
	(event: 'update:modelValue', value: FilterValues): void;
}
