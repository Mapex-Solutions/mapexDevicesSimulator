import { z } from '../common/primitives.schema';
import { ZodDeviceInputSchema } from '../devices/device.schema';

/**
 * The device marketplace contracts, mirroring the mapexMarketplace Go service
 * (`/api/v1/devices`). The catalog is browsed here and installed by the local
 * engine; protocol and reading types are plain strings so the catalog can grow
 * without a schema change.
 */

/** One catalog card returned by the listing endpoint. */
export const ZodMarketplaceCatalogItemSchema = z.object({
	id: z.string(),
	vendor: z.string(),
	vendorName: z.string(),
	model: z.string(),
	slug: z.string(),
	name: z.string(),
	description: z.string(),
	protocol: z.string(),
	readingTypes: z.array(z.string()),
	tags: z.array(z.string()).default([]),
	icon: z.string().default(''),
	hasCodec: z.boolean().default(false),
	hasManual: z.boolean().default(false),
});

/** A page of catalog cards plus the total match count. */
export const ZodMarketplaceListResponseSchema = z.object({
	items: z.array(ZodMarketplaceCatalogItemSchema),
	total: z.number(),
	page: z.number(),
	perPage: z.number(),
});

/** A selectable filter option with its display label and icon. */
export const ZodMarketplaceFacetOptionSchema = z.object({
	value: z.string(),
	label: z.string(),
	icon: z.string().default(''),
});

/** The filter options the listing UI renders (present in the catalog). */
export const ZodMarketplaceFacetsSchema = z.object({
	protocols: z.array(ZodMarketplaceFacetOptionSchema),
	readingTypes: z.array(ZodMarketplaceFacetOptionSchema),
	manufacturers: z.array(ZodMarketplaceFacetOptionSchema),
});

/** Vendor identity shown on the detail sheet. */
export const ZodMarketplaceVendorSchema = z.object({
	slug: z.string(),
	name: z.string(),
	site: z.string().default(''),
	support: z.string().default(''),
});

/** The detail sheet (`device_information.json`): description, image and files. */
export const ZodMarketplaceInformationSchema = z.object({
	id: z.string(),
	model: z.string(),
	name: z.string(),
	description: z.string(),
	protocol: z.string(),
	readingTypes: z.array(z.string()),
	tags: z.array(z.string()).default([]),
	version: z.string().optional(),
	vendor: ZodMarketplaceVendorSchema,
	images: z.object({ device: z.string().default('') }).default({}),
	files: z.object({ datasheet: z.string().default(''), manual: z.string().default('') }).default({}),
});

/**
 * One payload codec available for a model. The same device may expose several
 * (vendor-official plus community ports for different network servers).
 */
export const ZodMarketplaceCodecSchema = z.object({
	id: z.string(),
	name: z.string(),
	source: z.string().default(''),
	official: z.boolean().default(false),
	default: z.boolean().default(false),
	target: z.string().default(''),
	language: z.string().default(''),
	sourceUrl: z.string().default(''),
	path: z.string(),
	decoderFile: z.string().default(''),
	encoderFile: z.string().default(''),
});

/**
 * The installable template (`device_simulator.json`): a device input minus the
 * deviceId, which the engine mints fresh on install.
 */
export const ZodMarketplaceSimulatorSchema = ZodDeviceInputSchema.omit({ deviceId: true });
