import type { MenuItem } from '../interfaces';

import { protocolIcon } from '@components/protocols/ProtocolRegistry';

/**
 * Build the sidebar menu. The first version is intentionally small: the console
 * (home) for live logs and firing events, device registration grouped by
 * protocol, and gateways (needed for LoRaWAN). This is the one place to grow the
 * navigation.
 *
 * @param {(key: string) => string} t - translation function
 * @returns {MenuItem[]} translated menu items
 */
export function buildMenuList(t: (key: string) => string): MenuItem[] {
	return [
		{ icon: 'terminal', label: t('nav.console'), to: '/' },

		{
			icon: 'memory',
			label: t('nav.devices'),
			children: [
				{ icon: 'apps', label: t('nav.devicesAll'), to: '/devices' },
				{ icon: protocolIcon('http'), label: t('protocol.http'), to: '/devices?protocol=http' },
				{ icon: protocolIcon('mqtt'), label: t('protocol.mqtt'), to: '/devices?protocol=mqtt' },
				{ icon: protocolIcon('lorawan'), label: t('protocol.lorawan'), to: '/devices?protocol=lorawan' },
			],
		},

		{ icon: 'router', label: t('nav.gateways'), to: '/gateways' },
		{ icon: 'receipt_long', label: t('nav.logs'), to: '/logs' },

		// Marketplace sits last and, when the sidebar is collapsed, is set a little
		// apart from the rest (see the mini-state rule in AppSidebar).
		{ icon: 'storefront', label: t('nav.marketplace'), to: '/marketplace' },
	];
}
