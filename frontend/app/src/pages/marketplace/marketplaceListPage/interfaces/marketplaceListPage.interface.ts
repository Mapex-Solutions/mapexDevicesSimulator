/** TYPE IMPORTS */
import type { MarketplaceCatalogItem } from '@services/sim';

/** A resolved reading type (facet value → display label + icon). */
export interface MarketplaceReadingMeta {
	value: string;
	label: string;
	icon: string;
}

/**
 * A catalog item enriched with the display labels the card needs, so the card
 * stays presentational and the facet lookup happens once on the page.
 */
export interface MarketplaceCardItem extends MarketplaceCatalogItem {
	protocolLabel: string;
	readingMetas: MarketplaceReadingMeta[];
}

/** The active filter selection driving the catalog grid. */
export interface MarketplaceFilters {
	search: string;
	protocol: string | null;
	readingType: string | null;
	manufacturer: string | null;
}
