/** VUE IMPORTS */
import { createMemoryHistory, createRouter, createWebHashHistory, createWebHistory } from 'vue-router';

/** SERVICES */
import { defineRouter } from '#q-app/wrappers';

/** LOCAL IMPORTS */
import routes from './routes';

export default defineRouter(function () {
	const createHistory = process.env.SERVER
		? createMemoryHistory
		: process.env.VUE_ROUTER_MODE === 'history'
			? createWebHistory
			: createWebHashHistory;

	const Router = createRouter({
		scrollBehavior: () => ({ left: 0, top: 0 }),
		routes,
		history: createHistory(process.env.VUE_ROUTER_BASE),
	});

	return Router;
});
