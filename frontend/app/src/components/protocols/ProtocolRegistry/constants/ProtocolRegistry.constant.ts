/** TYPE IMPORTS */
import type { ProtocolConfig, ProtocolId } from '@services/sim';
import type { ProtocolDefinition, ValidationResult } from '../interfaces';

/** VUE IMPORTS */
import { markRaw } from 'vue';

/** COMPONENTS */
import { HttpConnectionConfig } from '../../HttpConnectionConfig';
import { MqttConnectionConfig } from '../../MqttConnectionConfig';
import { LoraWanConnectionConfig } from '../../LoraWanConnectionConfig';
import { BasicsStationConnectionConfig } from '../../BasicsStationConnectionConfig';

/**
 * Validate a connection config. Each protocol checks the fields of its own
 * config kind; the discriminator keeps access type-safe.
 * @param {ProtocolConfig} config - the config to check
 */
function validateConfig(config: ProtocolConfig): ValidationResult {
	const errors: string[] = [];

	if (config.kind === 'http') {
		if (!config.url.trim()) errors.push('url');
		if (config.authMode === 'apiKey' && !config.apiKey.trim()) errors.push('apiKey');
		if (config.authMode === 'basic' && !config.basicUser.trim()) errors.push('basicUser');
	}

	if (config.kind === 'mqtt') {
		if (!config.brokerUrl.trim()) errors.push('brokerUrl');
		if (config.authMode === 'userpass' && !config.username.trim()) errors.push('username');
		if (config.authMode === 'tls' && !config.tlsCertPem.trim()) errors.push('tlsCert');
	}

	if (config.kind === 'lorawan') {
		if (!config.gatewayId) errors.push('gateway');
		if (config.activation === 'otaa') {
			if (!config.devEui.trim()) errors.push('devEui');
			if (!config.appKey.trim()) errors.push('appKey');
			if (config.macVersion.startsWith('1.1') && !config.nwkKey.trim()) errors.push('nwkKey');
		}
		if (config.activation === 'abp') {
			if (!config.devAddr.trim()) errors.push('devAddr');
			if (!config.nwkSKey.trim()) errors.push('nwkSKey');
			if (!config.appSKey.trim()) errors.push('appSKey');
		}
	}

	if (config.kind === 'basicstation') {
		if (!config.lnsUri.trim()) errors.push('lnsUri');
		if (!config.gatewayEui.trim()) errors.push('gatewayEui');
		if (config.activation === 'otaa') {
			if (!config.devEui.trim()) errors.push('devEui');
			if (!config.appKey.trim()) errors.push('appKey');
			if (config.macVersion.startsWith('1.1') && !config.nwkKey.trim()) errors.push('nwkKey');
		}
		if (config.activation === 'abp') {
			if (!config.devAddr.trim()) errors.push('devAddr');
			if (!config.nwkSKey.trim()) errors.push('nwkSKey');
			if (!config.appSKey.trim()) errors.push('appSKey');
		}
	}

	return { valid: errors.length === 0, errors };
}

/**
 * Protocol contributions keyed by id. Only entries with `enabled: true` are
 * surfaced in the UI; disabled protocols are added here as they ship.
 */
export const PROTOCOL_REGISTRY: Partial<Record<ProtocolId, ProtocolDefinition>> = {
	http: {
		id: 'http',
		labelKey: 'protocol.http',
		icon: 'mdi-web',
		enabled: true,
		configComponent: markRaw(HttpConnectionConfig),
		defaultConfig: () => ({
			kind: 'http',
			url: 'http://localhost:5001/v1/ingest',
			method: 'POST',
			headers: [{ key: 'Content-Type', value: 'application/json' }],
			authMode: 'apiKey',
			apiKeyHeader: 'X-API-Key',
			apiKey: '',
			basicUser: '',
			basicPass: '',
		}),
		validate: validateConfig,
	},
	mqtt: {
		id: 'mqtt',
		labelKey: 'protocol.mqtt',
		icon: 'mdi-transit-connection-variant',
		enabled: true,
		configComponent: markRaw(MqttConnectionConfig),
		defaultConfig: () => ({
			kind: 'mqtt',
			brokerUrl: 'mqtt://127.0.0.1:1883',
			clientId: 'sim-device',
			baseTopic: '',
			authMode: 'userpass',
			username: '',
			password: '',
			tlsCertPem: '',
			tlsKeyPem: '',
			tlsCaPem: '',
			receiveEnabled: false,
			subscriptions: [],
		}),
		validate: validateConfig,
	},
	lorawan: {
		id: 'lorawan',
		labelKey: 'protocol.lorawan',
		icon: 'mdi-access-point',
		enabled: true,
		configComponent: markRaw(LoraWanConnectionConfig),
		defaultConfig: () => ({
			kind: 'lorawan',
			gatewayId: '',
			region: 'EU868',
			macVersion: '1.0.4',
			class: 'A',
			activation: 'otaa',
			devEui: '',
			joinEui: '',
			appKey: '',
			nwkKey: '',
			devAddr: '',
			nwkSKey: '',
			appSKey: '',
		}),
		validate: validateConfig,
	},
	basicstation: {
		id: 'basicstation',
		labelKey: 'protocol.basicstation',
		icon: 'mdi-radio-tower',
		enabled: true,
		configComponent: markRaw(BasicsStationConnectionConfig),
		defaultConfig: () => ({
			kind: 'basicstation',
			lnsUri: 'wss://127.0.0.1:1887',
			gatewayEui: '0016C001F1500099',
			region: 'EU868',
			macVersion: '1.0.4',
			class: 'A',
			activation: 'otaa',
			devEui: '',
			joinEui: '',
			appKey: '',
			nwkKey: '',
			devAddr: '',
			nwkSKey: '',
			appSKey: '',
		}),
		validate: validateConfig,
	},
};

/**
 * Enabled protocol definitions in display order.
 */
export const ENABLED_PROTOCOLS: ProtocolDefinition[] = Object.values(PROTOCOL_REGISTRY)
	.filter((def): def is ProtocolDefinition => Boolean(def?.enabled));
