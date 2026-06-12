#!/usr/bin/env node
/**
 * Builds the Go sidecar (`simulatord`) into bin/<platform>-<arch>/ so Electron's
 * sidecar manager can resolve it at runtime and quasar's extraResource ships it.
 *
 * Usage:
 *   node build-sidecar.mjs          # build for the current host only (fast, dev)
 *   node build-sidecar.mjs --all    # cross-compile every packaged target
 *
 * SQLite is modernc.org/sqlite (pure Go), so every target builds with
 * CGO_ENABLED=0 — no cross toolchains required.
 */
import { execFileSync } from 'node:child_process';
import { mkdirSync } from 'node:fs';
import { dirname, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

const scriptDir = dirname(fileURLToPath(import.meta.url));
const backendDir = resolve(scriptDir, '../../../../backend');
const binRoot = resolve(scriptDir, 'bin');
const mainPkg = './service/src';

// The folder name must match Electron's platformDir(): `${process.platform}-${process.arch}`.
const ALL_TARGETS = [
	{ dir: 'linux-x64', goos: 'linux', goarch: 'amd64' },
	{ dir: 'linux-arm64', goos: 'linux', goarch: 'arm64' },
	{ dir: 'darwin-x64', goos: 'darwin', goarch: 'amd64' },
	{ dir: 'darwin-arm64', goos: 'darwin', goarch: 'arm64' },
	{ dir: 'win32-x64', goos: 'windows', goarch: 'amd64' },
];

const GOOS_BY_PLATFORM = { linux: 'linux', darwin: 'darwin', win32: 'windows' };
const GOARCH_BY_ARCH = { x64: 'amd64', arm64: 'arm64' };

/**
 * Resolve the single host target from the current process platform/arch.
 * @returns {{dir: string, goos: string, goarch: string}} the host build target
 */
function hostTarget() {
	const goos = GOOS_BY_PLATFORM[process.platform];
	const goarch = GOARCH_BY_ARCH[process.arch];
	if (!goos || !goarch) {
		throw new Error(`unsupported host ${process.platform}-${process.arch}`);
	}
	return { dir: `${process.platform}-${process.arch}`, goos, goarch };
}

/**
 * Cross-compile one target into bin/<dir>/simulatord[.exe].
 * @param {{dir: string, goos: string, goarch: string}} target - the build target
 */
function build(target) {
	const name = target.goos === 'windows' ? 'simulatord.exe' : 'simulatord';
	const outDir = resolve(binRoot, target.dir);
	mkdirSync(outDir, { recursive: true });
	const out = resolve(outDir, name);
	process.stdout.write(`[sidecar] building ${target.goos}/${target.goarch} -> ${out}\n`);
	execFileSync('go', ['build', '-trimpath', '-ldflags', '-s -w', '-o', out, mainPkg], {
		cwd: backendDir,
		stdio: 'inherit',
		env: { ...process.env, GOOS: target.goos, GOARCH: target.goarch, CGO_ENABLED: '0' },
	});
}

const targets = process.argv.includes('--all') ? ALL_TARGETS : [hostTarget()];
for (const target of targets) build(target);
process.stdout.write(`[sidecar] done (${targets.length} target${targets.length > 1 ? 's' : ''})\n`);
