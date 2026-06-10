import { z } from '../common/primitives.schema';
import { ZodKeyValueSchema, ZodHttpMethodSchema, ZodMqttQoSSchema } from './protocol.schema';

/**
 * How a request body is authored: no body, raw text (JSON), or form fields
 * assembled into a JSON object. Shared by HTTP and MQTT events. Values may carry
 * template placeholders (for example {{randInt(10,30)}}) resolved at send time.
 */
export const ZodHttpBodyModeSchema = z.enum(['none', 'raw', 'form']);

export const ZodRequestBodySchema = z.object({
	bodyMode: ZodHttpBodyModeSchema,
	bodyFields: z.array(ZodKeyValueSchema),
	body: z.string(),
});

/** HTTP event configuration. Query parameters are authored in the path. */
export const ZodHttpEventConfigSchema = ZodRequestBodySchema.extend({
	method: ZodHttpMethodSchema,
	path: z.string(),
	headers: z.array(ZodKeyValueSchema),
});

/** MQTT event configuration: where and how the payload is published. */
export const ZodMqttEventConfigSchema = ZodRequestBodySchema.extend({
	topic: z.string(),
	qos: ZodMqttQoSSchema,
	retain: z.boolean(),
});

/**
 * LoRaWAN uplink event: an application-port frame carrying raw bytes. The payload
 * is authored as a hex string and may contain template placeholders.
 */
export const ZodLoraWanEventConfigSchema = z.object({
	fport: z.number(),
	confirmed: z.boolean(),
	payloadHex: z.string(),
});

/** Time unit for an event's repeat interval. */
export const ZodEventScheduleUnitSchema = z.enum(['seconds', 'minutes', 'hours', 'days']);

/**
 * Optional auto-fire schedule for an event. When enabled and the device is on,
 * the engine fires the event every `every` units.
 */
export const ZodEventScheduleSchema = z.object({
	enabled: z.boolean(),
	every: z.number(),
	unit: ZodEventScheduleUnitSchema,
});

/**
 * A pre-registered event on a device. Holds the protocol-specific config for the
 * device's protocol plus an optional repeat schedule.
 */
export const ZodDeviceEventSchema = z.object({
	id: z.string(),
	name: z.string(),
	http: ZodHttpEventConfigSchema.optional(),
	mqtt: ZodMqttEventConfigSchema.optional(),
	lorawan: ZodLoraWanEventConfigSchema.optional(),
	schedule: ZodEventScheduleSchema.optional(),
});
