/** TYPE IMPORTS */
import type { MqttConnectionConfig } from '@services/sim';

/** Props for the MQTT connection config. */
export interface MqttConnectionConfigProps {
	modelValue: MqttConnectionConfig;
}

/** Events emitted by the MQTT connection config. */
export interface MqttConnectionConfigEmits {
	(event: 'update:modelValue', value: MqttConnectionConfig): void;
}
