import { z } from '../../schemas/common/primitives.schema';
import { ZodEnvelopeSchema } from '../../schemas/common/envelope.schema';

/** The decoded `{ status, errors, data }` response envelope. */
export type Envelope<T = unknown> = Omit<z.infer<typeof ZodEnvelopeSchema>, 'data'> & { data: T };
