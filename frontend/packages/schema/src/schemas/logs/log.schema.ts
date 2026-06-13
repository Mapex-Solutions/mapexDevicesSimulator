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
	eventName: z.string().optional(),
	direction: ZodLogDirectionSchema,
	kind: ZodLogKindSchema,
	summary: z.string(),
	status: z.string().optional(),
	payload: z.string(),
	response: z.string().optional(),
});

/**
 * The cursor-paginated response for GET /api/logs. nextCursor is the opaque token
 * to pass back as `cursor` for the next page; absent when there are no more rows.
 */
export const ZodLogPageSchema = z.object({
	items: z.array(ZodLogResponseSchema),
	nextCursor: z.string().optional(),
});

/**
 * Query for GET /api/logs: a page limit plus optional filters. cursor is the
 * opaque keyset token of the previous page (absent for the first page). q is a
 * free-text match over summary, payload and device name; event matches the event
 * name; dateFrom/dateTo bound the message time.
 */
export const ZodLogQuerySchema = z.object({
	limit: z.number().int().min(0),
	cursor: z.string().optional(),
	protocol: z.string().optional(),
	kind: z.string().optional(),
	direction: z.string().optional(),
	device: z.string().optional(),
	event: z.string().optional(),
	dateFrom: z.string().optional(),
	dateTo: z.string().optional(),
	q: z.string().optional(),
});
