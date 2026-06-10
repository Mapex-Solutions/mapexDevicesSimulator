import { z } from '../common/primitives.schema';
import { ZodProtocolIdSchema } from '../devices/protocol.schema';
import { ZodLogDirectionSchema, ZodLogKindSchema } from '../logs/log.schema';

/**
 * One frame on the live console stream (WebSocket): every uplink, downlink and
 * auth/join handshake the engine emits, across every protocol. `ts` is the
 * frame's event time on the live stream, not a persisted entity's created.
 */
export const ZodConsoleMessageSchema = z.object({
	id: z.string(),
	ts: z.string(),
	protocol: ZodProtocolIdSchema,
	deviceId: z.string(),
	deviceName: z.string(),
	direction: ZodLogDirectionSchema,
	kind: ZodLogKindSchema,
	summary: z.string(),
	payload: z.string(),
	status: z.string().optional(),
	meta: z.record(z.string()).optional(),
});
