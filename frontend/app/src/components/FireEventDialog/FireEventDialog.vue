<template>
	<GenericModal v-model="open" :title="modalTitle" icon="mdi-flash" min-width="560px">
		<div class="column q-gutter-md fire">
			<q-select
				v-model="selectedDeviceId"
				:options="deviceOptions"
				:label="t('fireEvent.selectDevice')"
				outlined
				dense
				emit-value
				map-options
				hide-bottom-space
			/>

			<template v-if="device && isSupported">
				<!-- Mode picker: fire a registered event, or build a generic ad-hoc one. -->
				<div class="fire__modes">
					<button
						type="button"
						class="fire__mode"
						:class="{ 'fire__mode--active': mode === 'events' }"
						:disabled="!hasEvents"
						@click="mode = 'events'"
					>
						<q-icon name="mdi-playlist-play" size="22px" />
						<span class="fire__mode-label">{{ t('fireEvent.modeEvents') }}</span>
						<span class="fire__mode-hint">{{ t('fireEvent.modeEventsHint') }}</span>
					</button>
					<button
						type="button"
						class="fire__mode"
						:class="{ 'fire__mode--active': mode === 'generic' }"
						@click="mode = 'generic'"
					>
						<q-icon name="mdi-flash" size="22px" />
						<span class="fire__mode-label">{{ t('fireEvent.modeGeneric') }}</span>
						<span class="fire__mode-hint">{{ t('fireEvent.modeGenericHint') }}</span>
					</button>
				</div>

				<!-- Events: pick a pre-registered event and fire it as stored. -->
				<template v-if="mode === 'events'">
					<q-list v-if="hasEvents" bordered class="fire__events">
						<q-item
							v-for="evt in registeredEvents"
							:key="evt.id"
							clickable
							:active="eventId === evt.id"
							active-class="fire__event--active"
							@click="eventId = evt.id"
						>
							<q-item-section avatar>
								<q-icon :name="eventId === evt.id ? 'mdi-radiobox-marked' : 'mdi-radiobox-blank'" :color="eventId === evt.id ? 'primary' : 'grey-6'" />
							</q-item-section>
							<q-item-section>
								<q-item-label>{{ evt.name }}</q-item-label>
								<q-item-label caption>{{ summaryOf(evt) }}</q-item-label>
							</q-item-section>
						</q-item>
					</q-list>
					<q-banner v-else dense class="fire__unsupported">
						<template #avatar><q-icon name="mdi-information-outline" color="primary" /></template>
						{{ t('fireEvent.noEvents') }}
					</q-banner>
				</template>

				<!-- Generic: the protocol's own event editor (the connection comes from the device). -->
				<template v-else>
					<HttpEventConfig v-if="isHttp" v-model="httpConfig" />
					<MqttEventConfig v-else-if="isMqtt" v-model="mqttConfig" />
					<LoraWanEventConfig v-else v-model="loraConfig" />
				</template>

				<div v-if="preview">
					<div class="fire__preview-label">{{ t('fireEvent.preview') }}</div>
					<pre class="fire__preview">{{ preview }}</pre>
				</div>
			</template>

			<q-banner v-else-if="device" dense class="fire__unsupported">
				<template #avatar><q-icon name="mdi-information-outline" color="primary" /></template>
				{{ t('fireEvent.unsupported') }}
			</q-banner>
		</div>

		<template #footer>
			<q-btn flat :label="t('common.cancel')" @click="open = false" />
			<q-btn color="primary" icon="mdi-flash" :label="t('fireEvent.send')" :disable="!canSend" @click="onSend" />
		</template>
	</GenericModal>
</template>

<script setup lang="ts">
/** TYPE IMPORTS */
import type {
	DeviceEvent,
	HttpEventConfig as HttpEvent,
	LoraWanEventConfig as LoraEvent,
	MqttEventConfig as MqttEvent,
} from '@services/sim';

/** VUE IMPORTS */
import { computed, onMounted, ref, watch } from 'vue';

/** COMPONENTS */
import { GenericModal } from '@components/GenericModal';
import { HttpEventConfig, defaultHttpEvent } from '@components/protocols/HttpEventConfig';
import { MqttEventConfig, defaultMqttEvent } from '@components/protocols/MqttEventConfig';
import { LoraWanEventConfig, defaultLoraWanEvent } from '@components/protocols/LoraWanEventConfig';

/** COMPOSABLES */
import { useTranslations } from '@composables/i18n';

/** UTILS */
import { useQuasar } from 'quasar';
import { buildHttpBody, renderTemplate, validateJsonBody } from '@utils/template';
import { formatJson } from '@utils/format-json';

/** SERVICES */
import { sim } from '@services/sim';

/** STORES */
import { useDevicesStore } from '@stores/devices';

/** PROPS & EMITS */
const props = defineProps<{ deviceId?: string | null }>();
const open = defineModel<boolean>({ required: true });

/** COMPOSABLES & STORES */
const { t } = useTranslations();
const $q = useQuasar();
const devicesStore = useDevicesStore();

/** STATE */
const selectedDeviceId = ref<string | null>(null);
const mode = ref<'events' | 'generic'>('generic');
const eventId = ref<string | null>(null);
const httpConfig = ref<HttpEvent>(defaultHttpEvent());
const mqttConfig = ref<MqttEvent>(defaultMqttEvent());
const loraConfig = ref<LoraEvent>(defaultLoraWanEvent());

/** COMPUTED */
const deviceOptions = computed(() => devicesStore.items.map((device) => ({ label: device.name, value: device.id })));
const device = computed(() => devicesStore.items.find((item) => item.id === selectedDeviceId.value) ?? null);
const isHttp = computed(() => device.value?.protocolId === 'http');
const isMqtt = computed(() => device.value?.protocolId === 'mqtt');
const isLora = computed(() => device.value?.protocolId === 'lorawan' || device.value?.protocolId === 'basicstation');
const isSupported = computed(() => isHttp.value || isMqtt.value || isLora.value);

// Registered events that match the device's protocol — the source for Events mode.
const registeredEvents = computed(() =>
	(device.value?.events ?? []).filter((event) => (isLora.value ? event.lorawan : isMqtt.value ? event.mqtt : event.http)),
);
const hasEvents = computed(() => registeredEvents.value.length > 0);
const selectedEvent = computed(() => registeredEvents.value.find((event) => event.id === eventId.value) ?? null);

const modalTitle = computed(() => (device.value ? `${t('fireEvent.title')} · ${device.value.name}` : t('fireEvent.title')));

const renderCtx = computed(() => ({
	deviceId: device.value?.deviceId || device.value?.id,
	deviceName: device.value?.name,
}));

// What will actually be sent: the chosen registered event's payload in Events mode,
// or the generic editor's payload in Generic mode.
const preview = computed(() => {
	if (mode.value === 'events') return payloadOf(selectedEvent.value);
	return isLora.value
		? renderTemplate(loraConfig.value.payloadHex, renderCtx.value)
		: formatJson(buildHttpBody(isMqtt.value ? mqttConfig.value : httpConfig.value, renderCtx.value));
});

const bodyValid = computed(() => {
	if (isLora.value) return true;
	const config = isMqtt.value ? mqttConfig.value : httpConfig.value;
	if (config.bodyMode !== 'raw') return true;
	return validateJsonBody(config.body).valid;
});

const canSend = computed(() => {
	if (!device.value || !isSupported.value) return false;
	if (mode.value === 'events') return Boolean(eventId.value);
	return bodyValid.value;
});

/** WATCHERS */
watch(open, (isOpen) => {
	if (!isOpen) return;
	selectedDeviceId.value = props.deviceId ?? selectedDeviceId.value;
	resetForm();
});

watch(selectedDeviceId, resetForm);

/** FUNCTIONS */

/**
 * Deep clone an HTTP event config (its key/value arrays included).
 * @param {HttpEvent} config - the config to clone
 */
function cloneHttp(config: HttpEvent): HttpEvent {
	return {
		...config,
		headers: config.headers.map((row) => ({ ...row })),
		bodyFields: config.bodyFields.map((row) => ({ ...row })),
	};
}

/**
 * Deep clone an MQTT event config (its body fields included).
 * @param {MqttEvent} config - the config to clone
 */
function cloneMqtt(config: MqttEvent): MqttEvent {
	return {
		...config,
		bodyFields: config.bodyFields.map((row) => ({ ...row })),
	};
}

/**
 * Reset the form when the dialog opens or the device changes: blank generic
 * configs, and default to Events mode (first event preselected) when the device
 * has any, otherwise Generic.
 */
function resetForm(): void {
	httpConfig.value = defaultHttpEvent();
	mqttConfig.value = defaultMqttEvent();
	loraConfig.value = defaultLoraWanEvent();
	const events = registeredEvents.value;
	if (events.length) {
		mode.value = 'events';
		eventId.value = events[0]?.id ?? null;
	} else {
		mode.value = 'generic';
		eventId.value = null;
	}
}

/**
 * One-line summary of a registered event's request, for the Events list.
 * @param {DeviceEvent} event - the event to summarize
 * @returns {string} the summary text
 */
function summaryOf(event: DeviceEvent): string {
	if (event.http) return `${event.http.method} ${event.http.path}`;
	if (event.mqtt) return `QoS${event.mqtt.qos} ${event.mqtt.topic}`;
	if (event.lorawan) return `FPort ${event.lorawan.fport}`;
	return '';
}

/**
 * Render an event's payload for the preview (hex for LoRaWAN, body for HTTP/MQTT).
 * @param {DeviceEvent | null} event - the event to preview
 * @returns {string} the rendered payload
 */
function payloadOf(event: DeviceEvent | null): string {
	if (!event) return '';
	if (event.lorawan) return renderTemplate(event.lorawan.payloadHex, renderCtx.value);
	const body = event.mqtt ?? event.http;
	return body ? formatJson(buildHttpBody(body, renderCtx.value)) : '';
}

/**
 * Build the ad-hoc event for Generic mode from the protocol editor.
 * @returns {DeviceEvent} the ad-hoc event payload
 */
function buildGenericEvent(): DeviceEvent {
	const base = { id: '', name: t('fireEvent.generic') };
	if (isLora.value) return { ...base, lorawan: { ...loraConfig.value } };
	if (isMqtt.value) return { ...base, mqtt: cloneMqtt(mqttConfig.value) };
	return { ...base, http: cloneHttp(httpConfig.value) };
}

/**
 * Fire through the engine: a registered event by id in Events mode, or the ad-hoc
 * event in Generic mode. The engine renders + sends over the device's protocol;
 * the echo and any downlink stream back over the WebSocket.
 */
async function onSend(): Promise<void> {
	if (!device.value || !canSend.value) return;

	const body = mode.value === 'events' ? { eventId: eventId.value as string } : { event: buildGenericEvent() };

	try {
		await sim.devices.fire({ id: device.value.id }, body);
		$q.notify({ type: 'positive', message: t('fireEvent.sent') });
		open.value = false;
	} catch {
		$q.notify({ type: 'negative', message: t('fireEvent.failed') });
	}
}

/** LIFECYCLE HOOKS */
onMounted(() => {
	void devicesStore.fetch();
});
</script>

<style scoped lang="scss">
.fire {
	min-width: 520px;

	&__modes {
		display: flex;
		gap: var(--mapex-spacing-md);
	}

	&__mode {
		flex: 1;
		display: flex;
		flex-direction: column;
		align-items: flex-start;
		gap: var(--mapex-spacing-2xs);
		padding: var(--mapex-spacing-md);
		border: 1px solid var(--mapex-card-border);
		border-radius: var(--mapex-radius-md);
		background: var(--mapex-surface-elevated);
		color: var(--mapex-text-secondary);
		cursor: pointer;
		transition: var(--mapex-transition-fast);

		&:hover:not(:disabled) {
			border-color: var(--mapex-primary);
		}

		&:disabled {
			opacity: 0.5;
			cursor: not-allowed;
		}

		&--active {
			border-color: var(--mapex-primary);
			background: var(--mapex-active-bg);
			color: var(--mapex-primary);
		}
	}

	&__mode-label {
		font-size: var(--mapex-font-sm);
		font-weight: var(--mapex-font-weight-semibold);
		color: var(--mapex-text-primary);
	}

	&__mode-hint {
		font-size: var(--mapex-font-xs);
		color: var(--mapex-text-muted);
	}

	&__events {
		border-radius: var(--mapex-radius-md);
		max-height: 240px;
		overflow-y: auto;
	}

	&__event--active {
		background: var(--mapex-active-bg);
		color: var(--mapex-primary);
	}

	&__preview-label {
		font-size: var(--mapex-font-xs);
		text-transform: uppercase;
		letter-spacing: 0.04em;
		color: var(--mapex-text-secondary);
		margin-bottom: var(--mapex-spacing-2xs);
	}

	&__preview {
		margin: 0;
		padding: var(--mapex-spacing-md);
		border-radius: var(--mapex-radius-sm);
		background: var(--mapex-surface-sunken);
		color: var(--mapex-text-primary);
		font-family: var(--mapex-mono-font);
		font-size: var(--mapex-font-xs);
		white-space: pre-wrap;
		word-break: break-word;
	}

	&__unsupported {
		border-radius: var(--mapex-radius-sm);
		background: var(--mapex-surface-sunken);
		color: var(--mapex-text-secondary);
	}
}
</style>
