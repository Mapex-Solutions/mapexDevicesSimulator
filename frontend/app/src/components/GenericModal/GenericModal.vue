<template>
	<q-dialog v-model="open">
		<q-card class="generic-modal" :style="cardStyle">
			<q-card-section class="generic-modal__header row items-center no-wrap">
				<q-icon v-if="icon" :name="icon" size="sm" color="primary" class="q-mr-sm" />
				<div class="generic-modal__title">{{ title }}</div>
				<q-space />
				<q-btn v-close-popup flat round dense icon="close" color="grey-7" />
			</q-card-section>

			<q-separator />

			<q-card-section class="generic-modal__body">
				<slot />
			</q-card-section>

			<q-card-actions v-if="$slots.footer" align="right" class="generic-modal__footer">
				<slot name="footer" />
			</q-card-actions>
		</q-card>
	</q-dialog>
</template>

<script setup lang="ts">
/** TYPE IMPORTS */
import type { GenericModalProps } from './interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';

/** PROPS & EMITS */
const props = withDefaults(defineProps<GenericModalProps>(), {
	minWidth: '520px',
	maxWidth: '760px',
});

const open = defineModel<boolean>({ required: true });

/** COMPUTED */
const cardStyle = computed(() => ({ minWidth: props.minWidth, maxWidth: props.maxWidth }));
</script>

<style scoped lang="scss">
.generic-modal {
	border-radius: var(--mapex-radius-lg);
	max-width: 90vw;
	background: var(--mapex-surface-elevated);

	&__header {
		padding: var(--mapex-spacing-md) var(--mapex-spacing-lg);
	}

	&__title {
		font-size: var(--mapex-font-lg);
		font-weight: var(--mapex-font-weight-semibold);
		color: var(--mapex-text-primary);
	}

	&__body {
		padding: var(--mapex-spacing-lg);
	}

	&__footer {
		padding: var(--mapex-spacing-md) var(--mapex-spacing-lg);
		border-top: 1px solid var(--mapex-divider);
	}
}
</style>
