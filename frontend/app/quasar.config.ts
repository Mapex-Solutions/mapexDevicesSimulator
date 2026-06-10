// Quasar app configuration
// https://v2.quasar.dev/quasar-cli-vite/quasar-config-file

import { defineConfig } from '#q-app/wrappers';
import { fileURLToPath } from 'node:url';

export default defineConfig((ctx: any) => {
	return {
		// app boot files (/src/boot) — part of "main.js"
		boot: [
			'i18n',
		],

		css: [
			'app.scss',
			'page.scss',
		],

		extras: [
			'mdi-v7',
			'roboto-font',
			'material-icons',
		],

		build: {
			target: {
				browser: ['es2022', 'firefox115', 'chrome115', 'safari14'],
				node: 'node20',
			},

			typescript: {
				strict: true,
				vueShim: true,
			},

			vueRouterMode: 'history',

			alias: {
				'@src': fileURLToPath(new URL('./src', import.meta.url)),
				'@stores': fileURLToPath(new URL('./src/stores', import.meta.url)),
				'@components': fileURLToPath(new URL('./src/components', import.meta.url)),
				'@composables': fileURLToPath(new URL('./src/composables', import.meta.url)),
				'@utils': fileURLToPath(new URL('./src/utils', import.meta.url)),
				'@services': fileURLToPath(new URL('./src/services', import.meta.url)),
				'@interfaces': fileURLToPath(new URL('./src/interfaces', import.meta.url)),
				// Workspace packages (siblings under frontend/packages) consumed from TS
				// source with no build step: the Zod schemas (single source of truth)
				// and the typed sidecar client.
				'@sim/schema': fileURLToPath(new URL('../packages/schema/src', import.meta.url)),
				'@sim/api': fileURLToPath(new URL('../packages/api/src', import.meta.url)),
			},

			vitePlugins: [
				['@intlify/unplugin-vue-i18n/vite', {
					ssr: ctx.modeName === 'ssr',
					include: [fileURLToPath(new URL('./src/i18n', import.meta.url))],
				}],

				['vite-plugin-checker', {
					vueTsc: true,
				}, { server: false }],
			],
		},

		devServer: {
			port: 9100,
			open: false,
			// In dev the SPA is served by Vite while the Go engine runs separately;
			// proxy the control API and live stream to it. When the engine is down
			// these simply fail to connect (handled gracefully by the stores),
			// instead of Vite returning index.html for /api routes.
			proxy: {
				'/api': { target: 'http://127.0.0.1:5055', changeOrigin: true },
				'/ws': { target: 'http://127.0.0.1:5055', ws: true, changeOrigin: true },
			},
		},

		framework: {
			config: {},
			plugins: ['Notify', 'Dialog', 'LocalStorage'],
		},

		animations: [],

		// https://v2.quasar.dev/quasar-cli-vite/developing-electron-apps/configuring-electron
		electron: {
			preloadScripts: ['electron-preload'],

			inspectPort: 5858,

			bundler: 'packager',

			packager: {
				appBundleId: 'com.mapex.devices-simulator',
				// The Go sidecar binary is shipped alongside the app and resolved at
				// runtime under process.resourcesPath/sidecar so the desktop build can
				// serve the SPA and the control API without an external process.
				extraResource: [
					'src-electron/sidecar/bin',
				],
			},
		},
	};
});
