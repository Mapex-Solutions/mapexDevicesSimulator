/** TYPE IMPORTS */
import type { LoraWanEventConfig } from '@services/sim';

/**
 * A blank LoRaWAN uplink seeded with a templated hex payload so the
 * placeholders are discoverable.
 */
export function defaultLoraWanEvent(): LoraWanEventConfig {
	return {
		fport: 1,
		confirmed: false,
		payloadHex: '01{{randInt(0,255)}}',
	};
}
