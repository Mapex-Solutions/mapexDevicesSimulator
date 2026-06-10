import { z } from './primitives.schema';

/**
 * The platform's standard response envelope: every REST endpoint wraps its
 * payload as `{ status, errors, data }`. The api package unwraps `data` before it
 * reaches a caller, so this schema documents the transport rather than gating
 * call sites.
 */
export const ZodEnvelopeSchema = z.object({
	status: z.number(),
	errors: z.array(z.string()).nullable(),
	data: z.unknown(),
});
