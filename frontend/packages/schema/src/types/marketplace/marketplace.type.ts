import { z } from '../../schemas/common/primitives.schema';

import {
	ZodMarketplaceCatalogItemSchema,
	ZodMarketplaceListResponseSchema,
	ZodMarketplaceFacetOptionSchema,
	ZodMarketplaceFacetsSchema,
	ZodMarketplaceVendorSchema,
	ZodMarketplaceInformationSchema,
	ZodMarketplaceCodecSchema,
	ZodMarketplaceSimulatorSchema,
} from '../../schemas/marketplace/marketplace.schema';

export type MarketplaceCatalogItem = z.infer<typeof ZodMarketplaceCatalogItemSchema>;
export type MarketplaceListResponse = z.infer<typeof ZodMarketplaceListResponseSchema>;
export type MarketplaceFacetOption = z.infer<typeof ZodMarketplaceFacetOptionSchema>;
export type MarketplaceFacets = z.infer<typeof ZodMarketplaceFacetsSchema>;
export type MarketplaceVendor = z.infer<typeof ZodMarketplaceVendorSchema>;
export type MarketplaceInformation = z.infer<typeof ZodMarketplaceInformationSchema>;
export type MarketplaceCodec = z.infer<typeof ZodMarketplaceCodecSchema>;
export type MarketplaceSimulator = z.infer<typeof ZodMarketplaceSimulatorSchema>;
