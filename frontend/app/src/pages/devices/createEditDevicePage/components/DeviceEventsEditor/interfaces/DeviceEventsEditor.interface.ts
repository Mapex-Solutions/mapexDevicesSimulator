/** TYPE IMPORTS */
import type { DeviceEvent, ProtocolId } from '@services/sim';

/** Props for the device events editor. */
export interface DeviceEventsEditorProps {
	modelValue: DeviceEvent[];
	protocolId: ProtocolId;
}

/** Events emitted by the device events editor. */
export interface DeviceEventsEditorEmits {
	(event: 'update:modelValue', value: DeviceEvent[]): void;
}
