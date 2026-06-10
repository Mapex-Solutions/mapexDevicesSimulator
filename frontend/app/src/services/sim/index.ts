/** SERVICES */
import { createSimApi } from '@sim/api';

/**
 * Aggregated typed client for the sidecar control API, built from the @sim/api
 * package. Components and stores import this single object; the per-resource
 * clients, the Zod-validated payloads and the response-envelope unwrap all live
 * in the workspace packages (@sim/api + @sim/schema).
 */
export const sim = createSimApi();

/** Live console stream helpers, re-exported for the console page. */
export { resolveWsBase, createConsoleStream } from '@sim/api';
export type { ConsoleStream } from '@sim/api';

/**
 * Re-export every schema-inferred type so app code keeps importing its contracts
 * from a single facade (`@services/sim`) while the source of truth is @sim/schema.
 */
export type * from '@sim/schema';
