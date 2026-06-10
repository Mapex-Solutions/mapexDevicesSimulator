import { z } from '../common/primitives.schema';
import { ZodGatewayRegionSchema } from '../gateways/gateway.schema';

/** The protocol a device speaks; also the discriminant of its config. */
export const ZodProtocolIdSchema = z.enum(['http', 'mqtt', 'lorawan', 'basicstation']);

/** A key/value pair used for HTTP headers and form fields. */
export const ZodKeyValueSchema = z.object({
	key: z.string(),
	value: z.string(),
});

export const ZodHttpMethodSchema = z.enum(['POST', 'PUT']);
export const ZodHttpAuthModeSchema = z.enum(['none', 'apiKey', 'basic']);
export const ZodMqttAuthModeSchema = z.enum(['none', 'userpass', 'tls']);

/** LoRaWAN activation flavour and MAC versions. */
export const ZodLoraWanActivationSchema = z.enum(['otaa', 'abp']);
export const ZodLoraWanMacVersionSchema = z.enum(['1.0.2', '1.0.3', '1.0.4', '1.1.0']);

/**
 * HTTP target configuration. Credentials are sent only at run time and never
 * logged by the engine; headers carry Content-Type and any custom entries.
 */
export const ZodHttpConnectionConfigSchema = z.object({
	kind: z.literal('http'),
	url: z.string(),
	method: ZodHttpMethodSchema,
	headers: z.array(ZodKeyValueSchema),
	authMode: ZodHttpAuthModeSchema,
	apiKeyHeader: z.string(),
	apiKey: z.string(),
	basicUser: z.string(),
	basicPass: z.string(),
});

/**
 * MQTT target configuration. Auth is either username/password or a TLS client
 * certificate; credentials and PEM material are sent only at run time.
 */
export const ZodMqttConnectionConfigSchema = z.object({
	kind: z.literal('mqtt'),
	brokerUrl: z.string(),
	clientId: z.string(),
	baseTopic: z.string(),
	authMode: ZodMqttAuthModeSchema,
	username: z.string(),
	password: z.string(),
	tlsCertPem: z.string(),
	tlsKeyPem: z.string(),
	tlsCaPem: z.string(),
});

/**
 * LoRaWAN device configuration. The device transmits through an attached gateway
 * (referenced by id) which carries the LNS link. Keys are flat; OTAA fields apply
 * for the join flow, ABP fields for a pre-provisioned session.
 */
export const ZodLoraWanConnectionConfigSchema = z.object({
	kind: z.literal('lorawan'),
	gatewayId: z.string(),
	region: ZodGatewayRegionSchema,
	macVersion: ZodLoraWanMacVersionSchema,
	activation: ZodLoraWanActivationSchema,
	devEui: z.string(),
	joinEui: z.string(),
	appKey: z.string(),
	nwkKey: z.string(),
	devAddr: z.string(),
	nwkSKey: z.string(),
	appSKey: z.string(),
});

/**
 * Basics Station device configuration: a LoRaWAN node carrying its own Basics
 * Station link to the LNS (its own embedded gateway), instead of attaching to a
 * separate gateway. Shares the LoRaWAN identity fields.
 */
export const ZodBasicsStationConnectionConfigSchema = z.object({
	kind: z.literal('basicstation'),
	lnsUri: z.string(),
	gatewayEui: z.string(),
	region: ZodGatewayRegionSchema,
	macVersion: ZodLoraWanMacVersionSchema,
	activation: ZodLoraWanActivationSchema,
	devEui: z.string(),
	joinEui: z.string(),
	appKey: z.string(),
	nwkKey: z.string(),
	devAddr: z.string(),
	nwkSKey: z.string(),
	appSKey: z.string(),
});

/**
 * Per-protocol target configuration, discriminated by `kind`. Widened with a new
 * member as each protocol ships.
 */
export const ZodProtocolConfigSchema = z.discriminatedUnion('kind', [
	ZodHttpConnectionConfigSchema,
	ZodMqttConnectionConfigSchema,
	ZodLoraWanConnectionConfigSchema,
	ZodBasicsStationConnectionConfigSchema,
]);
