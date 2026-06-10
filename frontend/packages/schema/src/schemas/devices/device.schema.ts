import { z, StringAndNotBeEmpty } from '../common/primitives.schema';
import { ZodProtocolIdSchema, ZodProtocolConfigSchema } from './protocol.schema';
import { ZodDeviceEventSchema } from './event.schema';

/**
 * The wire shape returned for a simulated device. `created` is the ISO-8601
 * creation time (null until persisted). `config` carries the per-protocol target
 * and `events` the pre-registered events; the engine, not this module, interprets
 * them, so on the wire they pass through as JSON object / array.
 */
export const ZodDeviceResponseSchema = z.object({
	id: z.string(),
	created: z.string().nullable(),
	name: z.string(),
	deviceId: z.string(),
	protocolId: ZodProtocolIdSchema,
	enabled: z.boolean(),
	storeLogs: z.boolean(),
	config: ZodProtocolConfigSchema,
	attributes: z.record(z.string()),
	events: z.array(ZodDeviceEventSchema),
});

/**
 * The create/update body. id and created are server-assigned and excluded.
 */
export const ZodDeviceInputSchema = z.object({
	name: StringAndNotBeEmpty,
	deviceId: StringAndNotBeEmpty,
	protocolId: ZodProtocolIdSchema,
	enabled: z.boolean(),
	storeLogs: z.boolean(),
	config: ZodProtocolConfigSchema,
	attributes: z.record(z.string()),
	events: z.array(ZodDeviceEventSchema),
});
