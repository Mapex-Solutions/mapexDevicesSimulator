import { z } from '../common/primitives.schema';

/** The GET /api/health payload: liveness plus the running sidecar version. */
export const ZodHealthResponseSchema = z.object({
	status: z.string(),
	version: z.string(),
});
