/** The active marketplace filters sent to the catalog listing endpoint. */
export interface MarketplaceQuery {
	protocol?: string;
	readingType?: string;
	manufacturer?: string;
	search?: string;
	page?: number;
	perPage?: number;
}
