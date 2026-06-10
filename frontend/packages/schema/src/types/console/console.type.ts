import { z } from '../../schemas/common/primitives.schema';
import { ZodConsoleMessageSchema } from '../../schemas/console/console.schema';

export type ConsoleMessage = z.infer<typeof ZodConsoleMessageSchema>;
