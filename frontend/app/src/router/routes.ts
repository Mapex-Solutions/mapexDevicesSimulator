/** TYPE IMPORTS */
import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
	{
		path: '/',
		component: () => import('layouts/main/MainLayout.vue'),
		children: [
			{
				path: '',
				name: 'console',
				component: () => import('pages/console/consolePage/ConsolePage.vue'),
			},
			{
				path: 'devices',
				name: 'devices',
				component: () => import('pages/devices/deviceListPage/DeviceListPage.vue'),
			},
			{
				path: 'devices/new',
				name: 'device-new',
				component: () => import('pages/devices/createEditDevicePage/CreateEditDevicePage.vue'),
			},
			{
				path: 'devices/:id/edit',
				name: 'device-edit',
				component: () => import('pages/devices/createEditDevicePage/CreateEditDevicePage.vue'),
			},
			{
				path: 'marketplace',
				name: 'marketplace',
				component: () => import('pages/marketplace/marketplaceListPage/MarketplaceListPage.vue'),
			},
			{
				path: 'logs',
				name: 'logs',
				component: () => import('pages/logs/logListPage/LogListPage.vue'),
			},
			{
				path: 'gateways',
				name: 'gateways',
				component: () => import('pages/gateways/gatewayListPage/GatewayListPage.vue'),
			},
			{
				path: 'gateways/new',
				name: 'gateway-new',
				component: () => import('pages/gateways/createEditGatewayPage/CreateEditGatewayPage.vue'),
			},
			{
				path: 'gateways/:id/edit',
				name: 'gateway-edit',
				component: () => import('pages/gateways/createEditGatewayPage/CreateEditGatewayPage.vue'),
			},
		],
	},

	{
		path: '/:catchAll(.*)*',
		component: () => import('pages/errors/NotFoundPage.vue'),
	},
];

export default routes;
