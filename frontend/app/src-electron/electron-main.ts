/**
 * Electron main process. Spawns the Go sidecar on a free localhost port, waits
 * for it to report healthy, then opens a window. In dev the window loads the
 * Vite dev server; in production it loads the bundled SPA from disk. Either way
 * the renderer reaches the sidecar (API/WS) via the preload bridge — the sidecar
 * no longer serves the UI.
 */

import os from 'node:os';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import { app, BrowserWindow, shell } from 'electron';
import { pickFreePort, startSidecar, stopSidecar, waitForHealth } from './sidecar/sidecar-manager';

const currentDir = fileURLToPath(new URL('.', import.meta.url));
const platform = process.platform || os.platform();

let mainWindow: BrowserWindow | undefined;
let sidecarPort = 0;

const gotLock = app.requestSingleInstanceLock();

if (!gotLock) {
	app.quit();
} else {
	app.on('second-instance', () => {
		if (!mainWindow) return;
		if (mainWindow.isMinimized()) mainWindow.restore();
		mainWindow.focus();
	});

	void app.whenReady().then(bootstrap);
}

/**
 * Bring the engine up before showing any UI so the window never flashes a
 * disconnected state on a cold start.
 */
async function bootstrap(): Promise<void> {
	sidecarPort = await pickFreePort();

	const started = startSidecar(sidecarPort);
	if (started) {
		const healthy = await waitForHealth(sidecarPort);
		if (!healthy) console.error('[main] sidecar did not become healthy in time');
	}

	await createWindow();
}

/**
 * Create the application window and point it at the right SPA source.
 */
async function createWindow(): Promise<void> {
	const apiBase = `http://127.0.0.1:${sidecarPort}`;
	const wsBase = `ws://127.0.0.1:${sidecarPort}`;

	// 1366x768 is the minimum the UI is designed for; Electron enforces minWidth/
	// minHeight so the window can never be dragged below it. useContentSize makes
	// these the web viewport size (excludes the OS chrome). width/height are the
	// restore size when un-maximized; the window opens maximized (see below).
	mainWindow = new BrowserWindow({
		width: 1366,
		height: 800,
		minWidth: 1366,
		minHeight: 768,
		useContentSize: true,
		webPreferences: {
			contextIsolation: true,
			sandbox: false,
			additionalArguments: [
				`--sim-api-base=${apiBase}`,
				`--sim-ws-base=${wsBase}`,
				`--sim-app-version=${app.getVersion()}`,
			],
			preload: path.resolve(
				currentDir,
				path.join(
					process.env.QUASAR_ELECTRON_PRELOAD_FOLDER as string,
					`electron-preload${process.env.QUASAR_ELECTRON_PRELOAD_EXTENSION as string}`,
				),
			),
		},
	});

	// Open using the whole screen (respecting the OS taskbar) while keeping the
	// 1366x768 floor for when the user restores/resizes the window.
	mainWindow.maximize();

	// Open external http(s) links (e.g. the GitHub releases page) in the user's
	// default browser instead of a new Electron window.
	mainWindow.webContents.setWindowOpenHandler(({ url }) => {
		if (url.startsWith('http://') || url.startsWith('https://')) {
			void shell.openExternal(url);
		}
		return { action: 'deny' };
	});

	if (process.env.DEV) {
		await mainWindow.loadURL(process.env.APP_URL as string);
	} else {
		// Bundled SPA loaded from disk (file://), so the router runs in hash mode.
		// The sidecar is reached only for API/WS via the preload bridge.
		await mainWindow.loadFile(path.resolve(currentDir, 'index.html'));
	}

	if (process.env.DEBUGGING) {
		mainWindow.webContents.openDevTools();
	} else {
		mainWindow.webContents.on('devtools-opened', () => {
			mainWindow?.webContents.closeDevTools();
		});
	}

	mainWindow.on('closed', () => {
		mainWindow = undefined;
	});
}

app.on('window-all-closed', () => {
	stopSidecar();
	if (platform !== 'darwin') {
		app.quit();
	}
});

app.on('activate', () => {
	if (mainWindow === undefined) {
		void createWindow();
	}
});

app.on('before-quit', () => {
	stopSidecar();
});
