import { z } from '../../schemas/common/primitives.schema';
import { ZodHealthResponseSchema } from '../../schemas/health/health.schema';

export type HealthResponse = z.infer<typeof ZodHealthResponseSchema>;
