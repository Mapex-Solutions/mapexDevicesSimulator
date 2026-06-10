import { z } from '../common/primitives.schema';
import { ZodProtocolIdSchema } from '../devices/protocol.schema';

export const ZodLogDirectionSchema = z.enum(['up', 'down', 'system']);
export const ZodLogKindSchema = z.enum(['data', 'auth', 'join', 'downlink', 'status']);

/**
 * One persisted device message (the SQLite-backed history behind the live console
 * stream). `created` is the message time (ISO-8601, null only on malformed rows).
 */
export const ZodLogResponseSchema = z.object({
	id: z.string(),
	created: z.string().nullable(),
	protocol: ZodProtocolIdSchema,
	deviceId: z.string(),
	deviceName: z.string(),
	direction: ZodLogDirectionSchema,
	kind: ZodLogKindSchema,
	summary: z.string(),
	status: z.string().optional(),
	payload: z.string(),
});

/** The paginated response for GET /api/logs. */
export const ZodLogPageSchema = z.object({
	items: z.array(ZodLogResponseSchema),
	total: z.number(),
});

/**
 * Query for GET /api/logs: pagination plus optional filters. Empty filters are
 * ignored; q is a free-text match over summary, payload and device name.
 */
export const ZodLogQuerySchema = z.object({
	limit: z.number().int().min(0),
	offset: z.number().int().min(0),
	protocol: z.string().optional(),
	kind: z.string().optional(),
	direction: z.string().optional(),
	device: z.string().optional(),
	q: z.string().optional(),
});
