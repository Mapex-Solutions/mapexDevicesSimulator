/** The active marketplace filters sent to the catalog listing endpoint. */
export interface MarketplaceQuery {
	protocol?: string;
	readingType?: string;
	manufacturer?: string;
	search?: string;
	/** Active app locale (e.g. "pt-BR") so the catalog returns localized card descriptions. */
	lang?: string;
	page?: number;
	perPage?: number;
}
