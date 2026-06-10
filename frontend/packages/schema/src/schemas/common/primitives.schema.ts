import { z } from 'zod';

/**
 * Re-export zod's `z` (value plus its `infer` namespace, preserved by re-exporting
 * the binding rather than aliasing it to a const) so schema files import `z` from
 * here and type files can use `z.infer`. Shared primitives keep the per-field
 * intent terse and consistent across every schema.
 */
export { z };

export const IsString = z.string();
export const IsBoolean = z.boolean();
export const IsNumber = z.number();
export const IsRecord = z.record(z.string());

/** A non-empty string (required text fields). */
export const StringAndNotBeEmpty = z.string().min(1);

/** A non-empty string that may also be omitted. */
export const StringAndNotBeEmptyOrOptional = z.string().min(1).optional();

/** A positive integer (pagination, ports, counts). */
export const NumberIntAndPositive = z.number().int().positive();

/** A non-negative integer (offsets). */
export const NumberIntAndNonNegative = z.number().int().min(0);
