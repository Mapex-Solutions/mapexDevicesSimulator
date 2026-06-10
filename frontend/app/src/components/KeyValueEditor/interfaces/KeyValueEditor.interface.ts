/** TYPE IMPORTS */
import type { KeyValue } from '@services/sim';

/** Props for the key/value editor. */
export interface KeyValueEditorProps {
	rows: KeyValue[];
	addLabel: string;
	keyLabel?: string;
	valueLabel?: string;
	emptyLabel?: string;
}

/** Events emitted by the key/value editor. */
export interface KeyValueEditorEmits {
	(event: 'update', rows: KeyValue[]): void;
}
