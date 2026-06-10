/** TYPE IMPORTS */
import type { Gateway } from '@services/sim';

/**
 * Sample gateways used as a fallback while the simulator engine is offline, so
 * the UI is usable in development. Once the engine serves real data this seed is
 * ignored.
 * @returns {Gateway[]} the seed gateways
 */
export function seedGateways(): Gateway[] {
	return [
		{
			id: 'gw-seed-1',
			name: 'Rooftop Gateway',
			eui: '0016C001F1500001',
			enabled: true,
			region: 'EU868',
			description: 'Basics Station link to the local LNS.',
			link: { protocol: 'basicstation', lnsUri: 'wss://127.0.0.1:1887', host: '127.0.0.1', port: 1700 },
			created: '2026-05-20T11:00:00.000Z',
		},
		{
			id: 'gw-seed-2',
			name: 'Warehouse UDP',
			eui: '0016C001F1500002',
			enabled: false,
			region: 'US915',
			description: 'Semtech UDP packet forwarder.',
			link: { protocol: 'udp', lnsUri: '', host: '127.0.0.1', port: 1700 },
			created: '2026-06-01T16:30:00.000Z',
		},
	];
}
