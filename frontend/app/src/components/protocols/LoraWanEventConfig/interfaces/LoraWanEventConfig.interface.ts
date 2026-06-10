/** TYPE IMPORTS */
import type { LoraWanEventConfig } from '@services/sim';

/** Props for the LoRaWAN event config. */
export interface LoraWanEventConfigProps {
	modelValue: LoraWanEventConfig;
}

/** Events emitted by the LoRaWAN event config. */
export interface LoraWanEventConfigEmits {
	(event: 'update:modelValue', value: LoraWanEventConfig): void;
}
