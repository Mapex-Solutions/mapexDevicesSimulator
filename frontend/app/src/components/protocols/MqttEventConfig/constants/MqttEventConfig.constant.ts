/** TYPE IMPORTS */
import type { MqttEventConfig } from '@services/sim';

/**
 * A blank MQTT event seeded with a templated topic and body so the placeholders
 * are discoverable.
 */
export function defaultMqttEvent(): MqttEventConfig {
	return {
		topic: 'devices/{{deviceId}}/up',
		qos: 0,
		retain: false,
		bodyMode: 'raw',
		bodyFields: [],
		body: '{\n  "value": {{randInt(10,30)}}\n}',
	};
}
