// Quasar app configuration
// https://v2.quasar.dev/quasar-cli-vite/quasar-config-file

import { defineConfig } from '#q-app/wrappers';
import { fileURLToPath } from 'node:url';
import { readFileSync } from 'node:fs';

// Single source of truth for the version shown in the UI: package.json.
const pkg = JSON.parse(readFileSync(fileURLToPath(new URL('./package.json', import.meta.url)), 'utf8'));

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
			// Exposed to the renderer as process.env.APP_VERSION (replaced at build).
			env: {
				APP_VERSION: pkg.version,
			},

			target: {
				browser: ['es2022', 'firefox115', 'chrome115', 'safari14'],
				node: 'node20',
			},

			typescript: {
				strict: true,
				vueShim: true,
			},

			// hash mode: the packaged Electron app loads the SPA from file://, where
			// history mode can't match the path and would 404. Hash works everywhere.
			vueRouterMode: 'hash',

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

			bundler: 'builder',

			// electron-builder: produces the distributables. Quasar calls
			// electron-builder for the HOST OS only, so each OS block below activates
			// when you run the build on that OS (Linux → deb/rpm, macOS → dmg,
			// Windows → nsis). The Go sidecar ships under resources/bin/<plat>/ via
			// extraResources, where the sidecar manager resolves it at runtime; the
			// UI is the bundled SPA rendered by Electron.
			builder: {
				appId: 'com.mapex.devices-simulator',
				// No spaces: on Linux the install path becomes /opt/<productName>, and
				// spaces there break Chromium's zygote/sandbox child launch
				// ("failed to execvp: /opt/Mapex"). The friendly label lives in the
				// per-OS display names below (desktop.Name / CFBundleDisplayName /
				// nsis.shortcutName).
				productName: 'MapexDeviceSimulator',
				// electron is hoisted to the workspace root node_modules, so builder
				// cannot auto-detect it from app/ — pin the installed version.
				electronVersion: '33.4.11',

				// ---- Linux: Debian/Ubuntu (.deb) + RedHat/Fedora (.rpm) ----
				linux: {
					target: ['deb', 'rpm'],
					icon: 'src-electron/icons/icon.png',
					category: 'Utility',
					// Friendly name shown in the app menu / taskbar.
					desktop: { Name: 'Mapex Device Simulator' },
				},
				deb: {
					// Force the SUID chrome-sandbox (see the script) so the app runs
					// sandboxed on Ubuntu 24.04 without --no-sandbox.
					afterInstall: 'src-electron/deb/after-install.sh',
				},

				// ---- macOS (.dmg) — only builds when run ON macOS ----
				mac: {
					target: ['dmg'],
					icon: 'src-electron/icons/icon.png', // electron-builder derives .icns
					category: 'public.app-category.developer-tools',
					// Pretty name in Finder / menu bar (path stays MapexDeviceSimulator).
					extendInfo: { CFBundleDisplayName: 'Mapex Device Simulator' },
				},

				// ---- Windows (.exe / NSIS installer) — only builds when run ON Windows ----
				win: {
					target: ['nsis'],
					icon: 'src-electron/icons/icon.png', // electron-builder derives .ico
				},
				nsis: {
					oneClick: false,
					perMachine: false,
					allowToChangeInstallationDirectory: true,
					createDesktopShortcut: true,
					createStartMenuShortcut: true,
					shortcutName: 'Mapex Device Simulator',
				},

				extraResources: [
					{ from: 'src-electron/sidecar/bin', to: 'bin' },
				],
			},
		},
	};
});
