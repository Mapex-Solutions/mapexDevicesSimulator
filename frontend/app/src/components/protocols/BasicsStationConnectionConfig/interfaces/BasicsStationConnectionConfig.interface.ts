/** TYPE IMPORTS */
import type { BasicsStationConnectionConfig } from '@services/sim';

/** Props for the Basics Station connection config. */
export interface BasicsStationConnectionConfigProps {
	modelValue: BasicsStationConnectionConfig;
}

/** Events emitted by the Basics Station connection config. */
export interface BasicsStationConnectionConfigEmits {
	(event: 'update:modelValue', value: BasicsStationConnectionConfig): void;
}
