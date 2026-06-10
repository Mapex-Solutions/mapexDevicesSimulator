import { z, StringAndNotBeEmpty } from '../common/primitives.schema';

/** LoRaWAN regional frequency plan. */
export const ZodGatewayRegionSchema = z.enum([
	'EU868',
	'US915',
	'AU915',
	'AS923',
	'CN470',
	'IN865',
	'KR920',
	'RU864',
]);

/** Transport a gateway uses to reach the LNS. */
export const ZodGatewayLinkProtocolSchema = z.enum(['basicstation', 'udp']);

/**
 * How the gateway connects to the LNS. Basics Station uses a WebSocket LNS URI;
 * Semtech UDP uses a host and port. Both shapes are kept flat so a single object
 * covers either protocol.
 */
export const ZodGatewayLinkSchema = z.object({
	protocol: ZodGatewayLinkProtocolSchema,
	lnsUri: z.string(),
	host: z.string(),
	port: z.number(),
});

/**
 * A simulated LoRaWAN gateway: the Basics Station / UDP link to the LNS that
 * forwards frames from the simulated sensors. `created` is the ISO-8601 creation
 * time (null until persisted); `link` passes through as a JSON object.
 */
export const ZodGatewayResponseSchema = z.object({
	id: z.string(),
	created: z.string().nullable(),
	name: z.string(),
	eui: z.string(),
	enabled: z.boolean(),
	region: ZodGatewayRegionSchema,
	description: z.string(),
	link: ZodGatewayLinkSchema,
});

/** The create/update body. id and created are server-assigned and excluded. */
export const ZodGatewayInputSchema = z.object({
	name: StringAndNotBeEmpty,
	eui: StringAndNotBeEmpty,
	enabled: z.boolean(),
	region: ZodGatewayRegionSchema,
	description: z.string(),
	link: ZodGatewayLinkSchema,
});
