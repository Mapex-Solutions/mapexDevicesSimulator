/** TYPE IMPORTS */
import type { ProtocolId } from '@services/sim';
import type { FilterField } from './message-filters.interface';

/**
 * Filters available for every protocol: a free-text search and the message
 * direction.
 */
export const COMMON_MESSAGE_FILTERS: FilterField[] = [
	{ key: 'q', labelKey: 'console.filter.search', type: 'text', source: 'any' },
	{
		key: 'direction',
		labelKey: 'console.filter.direction',
		type: 'select',
		source: 'field',
		options: [
			{ value: 'up', labelKey: 'console.dirUp' },
			{ value: 'down', labelKey: 'console.dirDown' },
			{ value: 'system', labelKey: 'console.dirSystem' },
		],
	},
	{
		key: 'protocol',
		labelKey: 'console.filter.protocol',
		type: 'select',
		source: 'field',
		options: [
			{ value: 'http', labelKey: 'protocol.http' },
			{ value: 'mqtt', labelKey: 'protocol.mqtt' },
			{ value: 'lorawan', labelKey: 'protocol.lorawan' },
			{ value: 'basicstation', labelKey: 'protocol.basicstation' },
		],
	},
	{
		key: 'kind',
		labelKey: 'console.filter.kind',
		type: 'select',
		source: 'field',
		options: [
			{ value: 'data', labelKey: 'console.kind.data' },
			{ value: 'auth', labelKey: 'console.kind.auth' },
			{ value: 'join', labelKey: 'console.kind.join' },
			{ value: 'downlink', labelKey: 'console.kind.downlink' },
			{ value: 'status', labelKey: 'console.kind.status' },
		],
	},
];

/**
 * Per-protocol filter fields. Adding a protocol is just adding its entry here;
 * the console renders and applies whatever this exposes.
 */
export const PROTOCOL_MESSAGE_FILTERS: Partial<Record<ProtocolId, FilterField[]>> = {
	http: [
		{ key: 'status', labelKey: 'console.filter.status', type: 'text', source: 'field' },
		{ key: 'endpoint', labelKey: 'console.filter.endpoint', type: 'text', source: 'meta' },
	],
	mqtt: [
		{ key: 'topic', labelKey: 'console.filter.topic', type: 'text', source: 'meta' },
		{
			key: 'qos',
			labelKey: 'console.filter.qos',
			type: 'select',
			source: 'meta',
			options: [
				{ value: '0', label: '0' },
				{ value: '1', label: '1' },
				{ value: '2', label: '2' },
			],
		},
	],
	lorawan: [
		{ key: 'fPort', labelKey: 'console.filter.fport', type: 'text', source: 'meta' },
		{ key: 'gateway', labelKey: 'console.filter.gateway', type: 'text', source: 'meta' },
		{
			key: 'version',
			labelKey: 'console.filter.version',
			type: 'select',
			source: 'meta',
			options: [
				{ value: '1.0.3', label: '1.0.x' },
				{ value: '1.1', label: '1.1' },
			],
		},
	],
};

/**
 * Resolve the filter fields for a protocol context: the common fields plus the
 * protocol-specific ones. With no protocol, only the common fields apply.
 * @param {ProtocolId | null} protocol - the active protocol context
 * @returns {FilterField[]} the fields to render and apply
 */
export function getMessageFilterFields(protocol: ProtocolId | null): FilterField[] {
	if (!protocol) return COMMON_MESSAGE_FILTERS;
	return [...COMMON_MESSAGE_FILTERS, ...(PROTOCOL_MESSAGE_FILTERS[protocol] ?? [])];
}
