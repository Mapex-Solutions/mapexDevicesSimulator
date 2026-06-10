import { z } from '../../schemas/common/primitives.schema';
import {
	ZodLogDirectionSchema,
	ZodLogKindSchema,
	ZodLogResponseSchema,
	ZodLogPageSchema,
	ZodLogQuerySchema,
} from '../../schemas/logs/log.schema';

export type LogDirection = z.infer<typeof ZodLogDirectionSchema>;
export type LogKind = z.infer<typeof ZodLogKindSchema>;
export type Log = z.infer<typeof ZodLogResponseSchema>;
export type LogPage = z.infer<typeof ZodLogPageSchema>;
export type LogQuery = z.infer<typeof ZodLogQuerySchema>;
