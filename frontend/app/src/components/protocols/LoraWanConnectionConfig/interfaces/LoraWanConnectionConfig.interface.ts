/** TYPE IMPORTS */
import type { LoraWanConnectionConfig } from '@services/sim';

/** Props for the LoRaWAN connection config. */
export interface LoraWanConnectionConfigProps {
	modelValue: LoraWanConnectionConfig;
}

/** Events emitted by the LoRaWAN connection config. */
export interface LoraWanConnectionConfigEmits {
	(event: 'update:modelValue', value: LoraWanConnectionConfig): void;
}
