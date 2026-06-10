import { z } from '../../schemas/common/primitives.schema';
import {
	ZodGatewayRegionSchema,
	ZodGatewayLinkProtocolSchema,
	ZodGatewayLinkSchema,
	ZodGatewayResponseSchema,
	ZodGatewayInputSchema,
} from '../../schemas/gateways/gateway.schema';

export type GatewayRegion = z.infer<typeof ZodGatewayRegionSchema>;
export type GatewayLinkProtocol = z.infer<typeof ZodGatewayLinkProtocolSchema>;
export type GatewayLink = z.infer<typeof ZodGatewayLinkSchema>;
export type Gateway = z.infer<typeof ZodGatewayResponseSchema>;
export type GatewayInput = z.infer<typeof ZodGatewayInputSchema>;
