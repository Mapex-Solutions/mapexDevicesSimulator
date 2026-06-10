/** TYPE IMPORTS */
import type { MqttEventConfig } from '@services/sim';

/** Props for the MQTT event config. */
export interface MqttEventConfigProps {
	modelValue: MqttEventConfig;
}

/** Events emitted by the MQTT event config. */
export interface MqttEventConfigEmits {
	(event: 'update:modelValue', value: MqttEventConfig): void;
}
