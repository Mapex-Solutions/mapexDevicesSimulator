<template>
	<div class="events-editor column q-gutter-md">
		<!-- Editor form (always on top) -->
		<q-card flat bordered class="events-editor__form q-pa-md">
			<div class="events-editor__form-title">
				<q-icon name="bolt" color="primary" size="20px" />
				<span>{{ editingId ? t('deviceEvents.editingTitle') : t('deviceEvents.newTitle') }}</span>
			</div>

			<q-input
				v-model="form.name"
				:label="t('deviceEvents.name')"
				outlined
				dense
				stack-label
				hide-bottom-space
				class="q-mt-md"
			/>

			<component
				:is="configComponent"
				:model-value="activeConfig"
				class="q-mt-md"
				@update:model-value="onConfigUpdate"
			>
				<template #tabs>
					<q-tab name="schedule" :label="t('deviceEvents.scheduleTab')" />
				</template>
				<template #panel="{ active }">
					<div v-if="active === 'schedule'" class="events-editor__schedule">
						<q-checkbox
							:model-value="form.schedule.enabled"
							:label="t('deviceEvents.repeat')"
							@update:model-value="(v) => (form.schedule.enabled = v)"
						/>
						<div class="events-editor__schedule-hint">{{ t('deviceEvents.scheduleHint') }}</div>
						<div v-if="form.schedule.enabled" class="events-editor__schedule-row">
							<span class="events-editor__every">{{ t('deviceEvents.every') }}</span>
							<q-input
								class="events-editor__interval"
								type="number"
								min="1"
								:model-value="form.schedule.every"
								dense
								outlined
								hide-bottom-space
								@update:model-value="(v) => (form.schedule.every = clampEvery(v))"
							/>
							<q-select
								class="events-editor__unit"
								:model-value="form.schedule.unit"
								:options="unitOptions"
								dense
								outlined
								emit-value
								map-options
								@update:model-value="(v) => (form.schedule.unit = v)"
								hide-bottom-space
							/>
						</div>
					</div>
				</template>
			</component>

			<div class="row justify-end q-gutter-sm q-mt-md">
				<q-btn v-if="editingId" flat no-caps :label="t('common.cancel')" @click="resetForm" />
				<q-btn
					color="primary"
					no-caps
					:icon="editingId ? 'mdi-content-save' : 'mdi-plus'"
					:label="editingId ? t('deviceEvents.save') : t('deviceEvents.add')"
					:disable="!canSubmit"
					@click="submit"
				/>
			</div>
		</q-card>

		<!-- Added events (table below) -->
		<div class="events-editor__list">
			<div class="events-editor__list-title">{{ t('deviceEvents.listTitle', { count: modelValue.length }) }}</div>

			<q-markup-table v-if="modelValue.length" flat bordered dense wrap-cells>
				<thead>
					<tr>
						<th class="text-left">{{ t('deviceEvents.colName') }}</th>
						<th class="text-left">{{ t('deviceEvents.colRequest') }}</th>
						<th class="text-left">{{ t('deviceEvents.colCadence') }}</th>
						<th class="text-right">{{ t('deviceEvents.colActions') }}</th>
					</tr>
				</thead>
				<tbody>
					<tr v-for="event in modelValue" :key="event.id" :class="{ 'events-editor__row--active': event.id === editingId }">
						<td class="text-left">{{ event.name || t('deviceEvents.unnamed') }}</td>
						<td class="text-left"><code class="events-editor__summary">{{ summaryOf(event) }}</code></td>
						<td class="text-left">
							<span :class="event.schedule?.enabled ? 'events-editor__cadence--on' : 'events-editor__cadence--off'">
								{{ cadenceOf(event) }}
							</span>
						</td>
						<td class="text-right no-wrap">
							<q-btn flat dense round size="sm" icon="mdi-pencil" @click="editEvent(event)" />
							<q-btn flat dense round size="sm" icon="mdi-delete" color="negative" @click="removeEvent(event.id)" />
						</td>
					</tr>
				</tbody>
			</q-markup-table>

			<div v-else class="events-editor__empty">{{ t('deviceEvents.empty') }}</div>
		</div>
	</div>
</template>

<script setup lang="ts">
/** TYPE IMPORTS */
import type { Component } from 'vue';
import type { DeviceEvent, EventSchedule, EventScheduleUnit, HttpEventConfig as HttpEvent, LoraWanEventConfig as LoraEvent, MqttEventConfig as MqttEvent } from '@services/sim';
import type { DeviceEventsEditorEmits, DeviceEventsEditorProps } from './interfaces';

/** VUE IMPORTS */
import { computed, markRaw, reactive, ref } from 'vue';

/** COMPONENTS */
import { HttpEventConfig, defaultHttpEvent } from '@components/protocols/HttpEventConfig';
import { MqttEventConfig, defaultMqttEvent } from '@components/protocols/MqttEventConfig';
import { LoraWanEventConfig, defaultLoraWanEvent } from '@components/protocols/LoraWanEventConfig';

/** COMPOSABLES */
import { useTranslations } from '@composables/i18n';

const HTTP_CONFIG = markRaw(HttpEventConfig);
const MQTT_CONFIG = markRaw(MqttEventConfig);
const LORA_CONFIG = markRaw(LoraWanEventConfig);

/** Short suffix per time unit, used in the cadence summary. */
const UNIT_ABBR: Record<EventScheduleUnit, string> = { seconds: 's', minutes: 'min', hours: 'h', days: 'd' };

/**
 * A fresh, disabled schedule for a new event.
 * @returns {EventSchedule} the default schedule
 */
function defaultSchedule(): EventSchedule {
	return { enabled: false, every: 30, unit: 'seconds' };
}

/** PROPS & EMITS */
const props = defineProps<DeviceEventsEditorProps>();
const emit = defineEmits<DeviceEventsEditorEmits>();

/** COMPOSABLES & STORES */
const { t } = useTranslations();

/** STATE */
const editingId = ref<string | null>(null);
const form = reactive<{ name: string; http: HttpEvent; mqtt: MqttEvent; lorawan: LoraEvent; schedule: EventSchedule }>({
	name: '',
	http: defaultHttpEvent(),
	mqtt: defaultMqttEvent(),
	lorawan: defaultLoraWanEvent(),
	schedule: defaultSchedule(),
});
let seq = props.modelValue.reduce((max, event) => {
	const n = Number(event.id.replace(/\D/g, ''));
	return Number.isFinite(n) && n > max ? n : max;
}, 0);

/** COMPUTED */
const isHttp = computed(() => props.protocolId === 'http');
const isMqtt = computed(() => props.protocolId === 'mqtt');
const isLora = computed(() => props.protocolId === 'lorawan' || props.protocolId === 'basicstation');
const canSubmit = computed(() => form.name.trim().length > 0);

const configComponent = computed<Component>(() => {
	if (isMqtt.value) return MQTT_CONFIG;
	if (isLora.value) return LORA_CONFIG;
	return HTTP_CONFIG;
});
const activeConfig = computed<HttpEvent | MqttEvent | LoraEvent>(() => {
	if (isMqtt.value) return form.mqtt;
	if (isLora.value) return form.lorawan;
	return form.http;
});

const unitOptions = computed<{ label: string; value: EventScheduleUnit }[]>(() => [
	{ label: t('deviceEvents.units.seconds'), value: 'seconds' },
	{ label: t('deviceEvents.units.minutes'), value: 'minutes' },
	{ label: t('deviceEvents.units.hours'), value: 'hours' },
	{ label: t('deviceEvents.units.days'), value: 'days' },
]);

/** FUNCTIONS */

/**
 * Deep clone an HTTP event config so the form never mutates a stored event.
 * @param {HttpEvent} config - the config to clone
 * @returns {HttpEvent} a detached copy
 */
function cloneHttp(config: HttpEvent): HttpEvent {
	return {
		...config,
		headers: config.headers.map((row) => ({ ...row })),
		bodyFields: config.bodyFields.map((row) => ({ ...row })),
	};
}

/**
 * Deep clone an MQTT event config so the form never mutates a stored event.
 * @param {MqttEvent} config - the config to clone
 * @returns {MqttEvent} a detached copy
 */
function cloneMqtt(config: MqttEvent): MqttEvent {
	return {
		...config,
		bodyFields: config.bodyFields.map((row) => ({ ...row })),
	};
}

/**
 * Apply a config change from the active protocol component to the form.
 * @param {HttpEvent | MqttEvent} value - the updated config
 */
function onConfigUpdate(value: HttpEvent | MqttEvent | LoraEvent): void {
	if (isMqtt.value) form.mqtt = value as MqttEvent;
	else if (isLora.value) form.lorawan = value as LoraEvent;
	else form.http = value as HttpEvent;
}

/**
 * One-line summary of an event's request for the list.
 * @param {DeviceEvent} event - the event to summarize
 * @returns {string} the summary text
 */
function summaryOf(event: DeviceEvent): string {
	if (event.http) return `${event.http.method} ${event.http.path}`;
	if (event.mqtt) return `QoS${event.mqtt.qos} ${event.mqtt.topic}`;
	if (event.lorawan) return `FPort ${event.lorawan.fport} · ${event.lorawan.payloadHex}`;
	return '—';
}

/**
 * Cadence label for the list: the repeat interval or "Manual".
 * @param {DeviceEvent} event - the event to describe
 * @returns {string} the cadence text
 */
function cadenceOf(event: DeviceEvent): string {
	const schedule = event.schedule;
	if (!schedule?.enabled) return t('deviceEvents.manual');
	return `${t('deviceEvents.every')} ${schedule.every}${UNIT_ABBR[schedule.unit]}`;
}

/**
 * Clamp a raw interval input to a positive integer.
 * @param {string | number | null} value - the raw field value
 * @returns {number} the clamped interval
 */
function clampEvery(value: string | number | null): number {
	return Math.max(1, Math.floor(Number(value) || 1));
}

/**
 * Build an event from the current form for the given id.
 * @param {string} id - the event id
 * @returns {DeviceEvent} the assembled event
 */
function buildEvent(id: string): DeviceEvent {
	const base = { id, name: form.name.trim(), schedule: { ...form.schedule } };
	if (isMqtt.value) return { ...base, mqtt: cloneMqtt(form.mqtt) };
	if (isLora.value) return { ...base, lorawan: { ...form.lorawan } };
	return { ...base, http: cloneHttp(form.http) };
}

/**
 * Clear the form back to a blank new event.
 */
function resetForm(): void {
	editingId.value = null;
	form.name = '';
	form.http = defaultHttpEvent();
	form.mqtt = defaultMqttEvent();
	form.lorawan = defaultLoraWanEvent();
	form.schedule = defaultSchedule();
}

/**
 * Commit the form as a new event or as an update to the one being edited.
 */
function submit(): void {
	if (!canSubmit.value) return;
	if (editingId.value) {
		const id = editingId.value;
		emit('update:modelValue', props.modelValue.map((event) => (event.id === id ? buildEvent(id) : event)));
	} else {
		seq += 1;
		emit('update:modelValue', [...props.modelValue, buildEvent(`evt-${seq}`)]);
	}
	resetForm();
}

/**
 * Load an existing event into the form for editing.
 * @param {DeviceEvent} event - the event to edit
 */
function editEvent(event: DeviceEvent): void {
	editingId.value = event.id;
	form.name = event.name;
	form.http = event.http ? cloneHttp(event.http) : defaultHttpEvent();
	form.mqtt = event.mqtt ? cloneMqtt(event.mqtt) : defaultMqttEvent();
	form.lorawan = event.lorawan ? { ...event.lorawan } : defaultLoraWanEvent();
	form.schedule = event.schedule ? { ...event.schedule } : defaultSchedule();
}

/**
 * Remove an event, resetting the form if it was the one being edited.
 * @param {string} id - the event id to remove
 */
function removeEvent(id: string): void {
	emit('update:modelValue', props.modelValue.filter((event) => event.id !== id));
	if (editingId.value === id) resetForm();
}
</script>

<style scoped lang="scss">
.events-editor {
	&__form-title {
		display: flex;
		align-items: center;
		gap: var(--mapex-spacing-xs);
		font-weight: var(--mapex-font-weight-medium);
		color: var(--mapex-text-primary);
	}

	&__list-title {
		font-size: var(--mapex-font-xs);
		font-weight: var(--mapex-font-weight-medium);
		color: var(--mapex-text-secondary);
		text-transform: uppercase;
		letter-spacing: 0.04em;
		margin-bottom: var(--mapex-spacing-xs);
	}

	&__summary {
		font-family: var(--mapex-mono-font);
		font-size: var(--mapex-font-xs);
		color: var(--mapex-text-secondary);
	}

	&__schedule {
		padding-top: var(--mapex-spacing-xs);
	}

	&__schedule-hint {
		font-size: var(--mapex-font-xs);
		color: var(--mapex-text-muted);
		margin-top: var(--mapex-spacing-2xs);
	}

	&__schedule-row {
		display: flex;
		align-items: center;
		gap: var(--mapex-spacing-sm);
		margin-top: var(--mapex-spacing-sm);
	}

	&__every {
		font-size: var(--mapex-font-sm);
		color: var(--mapex-text-secondary);
	}

	&__interval {
		width: 90px;
	}

	&__unit {
		min-width: 130px;
	}

	&__cadence {
		&--on {
			font-size: var(--mapex-font-xs);
			font-weight: var(--mapex-font-weight-medium);
			color: var(--mapex-primary);
		}

		&--off {
			font-size: var(--mapex-font-xs);
			color: var(--mapex-text-muted);
		}
	}

	&__row--active {
		background: var(--mapex-surface-sunken);
	}

	&__empty {
		padding: var(--mapex-spacing-md);
		border: 1px dashed var(--mapex-card-border);
		border-radius: var(--mapex-radius-md);
		text-align: center;
		font-size: var(--mapex-font-sm);
		color: var(--mapex-text-muted);
	}
}
</style>
