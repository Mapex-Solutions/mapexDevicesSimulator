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
				<q-select
					v-model="eventId"
					:options="eventOptions"
					:label="t('fireEvent.event')"
					outlined
					dense
					emit-value
					map-options
					hide-bottom-space
				/>

				<HttpEventConfig v-if="isHttp" v-model="httpConfig" />
				<MqttEventConfig v-else-if="isMqtt" v-model="mqttConfig" />
				<LoraWanEventConfig v-else v-model="loraConfig" />

				<div>
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
			<q-btn color="primary" icon="mdi-flash" :label="t('fireEvent.send')" :disable="!device || !isSupported" @click="onSend" />
		</template>
	</GenericModal>
</template>

<script setup lang="ts">
/** TYPE IMPORTS */
import type { HttpEventConfig as HttpEvent, LoraWanEventConfig as LoraEvent, MqttEventConfig as MqttEvent } from '@services/sim';

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
import { buildHttpBody, renderTemplate } from '@utils/template';

/** STORES */
import { useDevicesStore } from '@stores/devices';
import { useMessagesStore } from '@stores/messages';

/** PROPS & EMITS */
const props = defineProps<{ deviceId?: string | null }>();
const open = defineModel<boolean>({ required: true });

/** COMPOSABLES & STORES */
const { t } = useTranslations();
const $q = useQuasar();
const devicesStore = useDevicesStore();
const messagesStore = useMessagesStore();

/** STATE */
const selectedDeviceId = ref<string | null>(null);
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

const eventOptions = computed(() => [
	{ label: t('fireEvent.generic'), value: null },
	...(device.value?.events ?? [])
		.filter((event) => (isLora.value ? event.lorawan : isMqtt.value ? event.mqtt : event.http))
		.map((event) => ({ label: event.name, value: event.id })),
]);

const modalTitle = computed(() => (device.value ? `${t('fireEvent.title')} · ${device.value.name}` : t('fireEvent.title')));

const renderCtx = computed(() => ({
	deviceId: device.value?.deviceId || device.value?.id,
	deviceName: device.value?.name,
}));

const preview = computed(() =>
	isLora.value
		? renderTemplate(loraConfig.value.payloadHex, renderCtx.value)
		: buildHttpBody(isMqtt.value ? mqttConfig.value : httpConfig.value, renderCtx.value),
);

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
 * Reset both configs to blank generic events.
 */
function resetConfigs(): void {
	httpConfig.value = defaultHttpEvent();
	mqttConfig.value = defaultMqttEvent();
	loraConfig.value = defaultLoraWanEvent();
}

/**
 * Load the form from the selected pre-registered event, or a blank generic one.
 */
function loadEvent(): void {
	if (!eventId.value) {
		resetConfigs();
		return;
	}
	const event = device.value?.events.find((item) => item.id === eventId.value);
	httpConfig.value = event?.http ? cloneHttp(event.http) : defaultHttpEvent();
	mqttConfig.value = event?.mqtt ? cloneMqtt(event.mqtt) : defaultMqttEvent();
	loraConfig.value = event?.lorawan ? { ...event.lorawan } : defaultLoraWanEvent();
}

/**
 * Fire the event: render the template and record it on the console stream.
 */
function onSend(): void {
	if (!device.value || !isSupported.value) return;

	const ctx = { deviceId: device.value.deviceId || device.value.id, deviceName: device.value.name };
	const rendered = isLora.value
		? renderTemplate(loraConfig.value.payloadHex, ctx)
		: buildHttpBody(isMqtt.value ? mqttConfig.value : httpConfig.value, ctx);
	const eventName = eventId.value
		? (device.value.events.find((item) => item.id === eventId.value)?.name ?? t('fireEvent.generic'))
		: t('fireEvent.generic');
	const status = isLora.value
		? `FPort ${loraConfig.value.fport}${loraConfig.value.confirmed ? ' · confirmed' : ''}`
		: isMqtt.value
			? `${renderTopic(mqttConfig.value.topic, ctx)} (QoS${mqttConfig.value.qos})`
			: `${httpConfig.value.method} ${httpConfig.value.path}`;

	messagesStore.add({
		ts: new Date().toLocaleTimeString(),
		protocol: isLora.value ? 'lorawan' : isMqtt.value ? 'mqtt' : 'http',
		deviceId: device.value.id,
		deviceName: device.value.name,
		direction: 'down',
		kind: 'downlink',
		status,
		summary: eventName,
		payload: rendered,
	});

	$q.notify({ type: 'warning', message: t('fireEvent.offline') });
	open.value = false;
}

/**
 * Render the placeholders inside an MQTT topic for the console summary.
 * @param {string} topic - the topic template
 * @param {{ deviceId: string; deviceName: string }} ctx - the render context
 * @returns {string} the rendered topic
 */
function renderTopic(topic: string, ctx: { deviceId: string; deviceName: string }): string {
	return buildHttpBody({ bodyMode: 'raw', bodyFields: [], body: topic }, ctx);
}

/** WATCHERS */
watch(open, (isOpen) => {
	if (!isOpen) return;
	selectedDeviceId.value = props.deviceId ?? selectedDeviceId.value;
	eventId.value = null;
	resetConfigs();
});

watch(selectedDeviceId, () => {
	eventId.value = null;
	resetConfigs();
});

watch(eventId, loadEvent);

/** LIFECYCLE HOOKS */
onMounted(() => {
	void devicesStore.fetch();
});
</script>

<style scoped lang="scss">
.fire {
	min-width: 520px;

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
