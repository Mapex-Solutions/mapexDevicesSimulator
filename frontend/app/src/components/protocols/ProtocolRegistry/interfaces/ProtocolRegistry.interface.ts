/** TYPE IMPORTS */
import type { Component } from 'vue';
import type { ProtocolConfig, ProtocolId } from '@services/sim';

export interface ValidationResult {
	valid: boolean;
	errors: string[];
}

/**
 * One protocol's contribution to the UI. Turning a protocol on is a matter of
 * flipping `enabled` and supplying its config component, default config and
 * validator; pages and stores read this registry instead of hardcoding protocol
 * knowledge.
 */
export interface ProtocolDefinition {
	id: ProtocolId;
	labelKey: string;
	icon: string;
	enabled: boolean;
	configComponent: Component;
	defaultConfig: () => ProtocolConfig;
	validate: (config: ProtocolConfig) => ValidationResult;
}
