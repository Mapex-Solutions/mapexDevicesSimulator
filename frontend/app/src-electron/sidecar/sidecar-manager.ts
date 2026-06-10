/**
 * Go sidecar lifecycle: locate the binary, pick a free port, spawn it, wait for
 * it to become healthy, and stop it on shutdown. Kept out of electron-main so
 * the orchestration stays small and the resolution logic is testable.
 */

import { type ChildProcess, spawn } from 'node:child_process';
import { existsSync } from 'node:fs';
import http from 'node:http';
import net from 'node:net';
import path from 'node:path';
import { app } from 'electron';

const BINARY_BASENAME = 'simulatord';

let sidecar: ChildProcess | null = null;
let stopped = false;

/**
 * Platform-specific binary file name.
 */
function binaryName(): string {
	return process.platform === 'win32' ? `${BINARY_BASENAME}.exe` : BINARY_BASENAME;
}

/**
 * Per-platform subfolder the backend build writes into.
 */
function platformDir(): string {
	return `${process.platform}-${process.arch}`;
}

/**
 * Resolve the sidecar binary path across dev and packaged layouts. Returns the
 * first candidate that exists, or null when no build is present yet.
 */
export function resolveSidecarBinary(): string | null {
	const name = binaryName();
	const platDir = platformDir();

	const candidates = app.isPackaged
		? [
				path.join(process.resourcesPath, 'bin', platDir, name),
				path.join(process.resourcesPath, 'bin', name),
			]
		: [
				path.join(app.getAppPath(), 'src-electron', 'sidecar', 'bin', platDir, name),
				path.join(process.cwd(), 'src-electron', 'sidecar', 'bin', platDir, name),
				path.join(process.cwd(), 'src-electron', 'sidecar', 'bin', name),
			];

	return candidates.find((candidate) => existsSync(candidate)) ?? null;
}

/**
 * Ask the OS for a free TCP port by binding to port 0 and reading it back.
 */
export function pickFreePort(): Promise<number> {
	return new Promise((resolve, reject) => {
		const server = net.createServer();
		server.unref();
		server.on('error', reject);
		server.listen(0, '127.0.0.1', () => {
			const address = server.address();
			if (address && typeof address === 'object') {
				const { port } = address;
				server.close(() => resolve(port));
			} else {
				server.close(() => reject(new Error('failed to acquire a free port')));
			}
		});
	});
}

/**
 * Spawn the sidecar bound to the given port. Returns false when no binary is
 * available so the app can still open and surface a disconnected engine state.
 * @param {number} port - the localhost port the sidecar must listen on
 */
export function startSidecar(port: number): boolean {
	const binary = resolveSidecarBinary();
	if (!binary) {
		console.warn('[sidecar] binary not found; start the backend build to enable the engine');
		return false;
	}

	stopped = false;
	sidecar = spawn(binary, ['--addr', '127.0.0.1', '--port', String(port)], {
		stdio: ['ignore', 'pipe', 'pipe'],
	});

	sidecar.stdout?.on('data', (chunk: Buffer) => console.log('[sidecar]', chunk.toString().trimEnd()));
	sidecar.stderr?.on('data', (chunk: Buffer) => console.error('[sidecar]', chunk.toString().trimEnd()));
	sidecar.on('exit', (code) => {
		sidecar = null;
		if (!stopped) console.error(`[sidecar] exited unexpectedly with code ${code ?? 'null'}`);
	});

	return true;
}

/**
 * Poll the sidecar health endpoint until it responds 200 or the timeout passes.
 * @param {number} port - the localhost port to probe
 * @param {number} timeoutMs - total time budget before giving up
 */
export async function waitForHealth(port: number, timeoutMs = 10000): Promise<boolean> {
	const deadline = Date.now() + timeoutMs;
	const url = `http://127.0.0.1:${port}/api/health`;

	while (Date.now() < deadline) {
		const ok = await probe(url);
		if (ok) return true;
		await delay(150);
	}

	return false;
}

/**
 * Single health probe resolving to true only on a 200 response.
 * @param {string} url - the health endpoint
 */
function probe(url: string): Promise<boolean> {
	return new Promise((resolve) => {
		const req = http.get(url, (res) => {
			res.resume();
			resolve(res.statusCode === 200);
		});
		req.on('error', () => resolve(false));
		req.setTimeout(1000, () => {
			req.destroy();
			resolve(false);
		});
	});
}

/**
 * Stop the sidecar gracefully, escalating to SIGKILL if it lingers.
 */
export function stopSidecar(): void {
	if (!sidecar) return;
	stopped = true;

	const child = sidecar;
	sidecar = null;
	child.kill('SIGTERM');

	const killTimer = setTimeout(() => {
		if (!child.killed) child.kill('SIGKILL');
	}, 5000);
	killTimer.unref?.();
}

function delay(ms: number): Promise<void> {
	return new Promise((resolve) => setTimeout(resolve, ms));
}

// Best-effort cleanup if Electron is force-terminated.
app.on('before-quit', stopSidecar);
process.on('exit', stopSidecar);
