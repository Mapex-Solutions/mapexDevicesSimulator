import { z } from '../../schemas/common/primitives.schema';

import {
	ZodProtocolIdSchema,
	ZodKeyValueSchema,
	ZodHttpMethodSchema,
	ZodHttpAuthModeSchema,
	ZodMqttAuthModeSchema,
	ZodLoraWanActivationSchema,
	ZodLoraWanMacVersionSchema,
	ZodHttpConnectionConfigSchema,
	ZodMqttConnectionConfigSchema,
	ZodLoraWanConnectionConfigSchema,
	ZodBasicsStationConnectionConfigSchema,
	ZodProtocolConfigSchema,
} from '../../schemas/devices/protocol.schema';

import {
	ZodHttpBodyModeSchema,
	ZodRequestBodySchema,
	ZodHttpEventConfigSchema,
	ZodMqttQoSSchema,
	ZodMqttEventConfigSchema,
	ZodLoraWanEventConfigSchema,
	ZodEventScheduleUnitSchema,
	ZodEventScheduleSchema,
	ZodDeviceEventSchema,
} from '../../schemas/devices/event.schema';

import { ZodDeviceResponseSchema, ZodDeviceInputSchema } from '../../schemas/devices/device.schema';

export type ProtocolId = z.infer<typeof ZodProtocolIdSchema>;
export type KeyValue = z.infer<typeof ZodKeyValueSchema>;
export type HttpMethod = z.infer<typeof ZodHttpMethodSchema>;
export type HttpAuthMode = z.infer<typeof ZodHttpAuthModeSchema>;
export type MqttAuthMode = z.infer<typeof ZodMqttAuthModeSchema>;
export type LoraWanActivation = z.infer<typeof ZodLoraWanActivationSchema>;
export type LoraWanMacVersion = z.infer<typeof ZodLoraWanMacVersionSchema>;

export type HttpConnectionConfig = z.infer<typeof ZodHttpConnectionConfigSchema>;
export type MqttConnectionConfig = z.infer<typeof ZodMqttConnectionConfigSchema>;
export type LoraWanConnectionConfig = z.infer<typeof ZodLoraWanConnectionConfigSchema>;
export type BasicsStationConnectionConfig = z.infer<typeof ZodBasicsStationConnectionConfigSchema>;
export type ProtocolConfig = z.infer<typeof ZodProtocolConfigSchema>;

export type HttpBodyMode = z.infer<typeof ZodHttpBodyModeSchema>;
export type RequestBody = z.infer<typeof ZodRequestBodySchema>;
export type HttpEventConfig = z.infer<typeof ZodHttpEventConfigSchema>;
export type MqttQoS = z.infer<typeof ZodMqttQoSSchema>;
export type MqttEventConfig = z.infer<typeof ZodMqttEventConfigSchema>;
export type LoraWanEventConfig = z.infer<typeof ZodLoraWanEventConfigSchema>;
export type EventScheduleUnit = z.infer<typeof ZodEventScheduleUnitSchema>;
export type EventSchedule = z.infer<typeof ZodEventScheduleSchema>;
export type DeviceEvent = z.infer<typeof ZodDeviceEventSchema>;

export type Device = z.infer<typeof ZodDeviceResponseSchema>;
export type DeviceInput = z.infer<typeof ZodDeviceInputSchema>;
