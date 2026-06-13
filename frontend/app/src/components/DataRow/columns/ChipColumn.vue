<template>
  <div :class="containerClass">
    <q-chip
        dense
        outline
        size="sm"
        text-color="white"
        class="text-weight-medium"
        style="width: fit-content; max-width: 100%;"
        :color="getColor()"
        :label="displayValue"
    >
      <AppTooltip :content="displayValue" />
    </q-chip>
    <div v-if="secondaryText" class="text-caption text-grey-6 ellipsis q-mt-xs">
      {{ secondaryText }}
    </div>
  </div>
</template>

<script setup lang="ts">
/** TYPE IMPORTS */
import type { DataRowColumnProps } from '../interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPONENTS */
import { AppTooltip } from '@components/AppTooltip';

defineOptions({
  name: 'ChipColumn'
});

const props = defineProps<DataRowColumnProps>();

const displayValue = computed(() => {
  if (props.column.format) {
    return props.column.format(props.value, props.row);
  }
  return (props.value as string) || 'N/A';
});

const secondaryText = computed(() => {
  if (props.column.secondary) {
    return props.column.secondary(props.value, props.row);
  }
  return '';
});

function getColor() {
  if (typeof props.column.color === 'function') {
    return props.column.color(props.value, props.row);
  }
  return props.column.color || 'primary';
}

const containerClass = computed(() => {
  const classes = ['flex'];

  // Stack chip + secondary caption vertically when a caption is present.
  if (secondaryText.value) {
    classes.push('column');
  }

  // Apply alignment from column config
  if (props.column.align === 'center') {
    classes.push(secondaryText.value ? 'items-center' : 'justify-center');
  } else if (props.column.align === 'right') {
    classes.push(secondaryText.value ? 'items-end' : 'justify-end');
  } else {
    classes.push(secondaryText.value ? 'items-start' : 'justify-start');
  }

  return classes.join(' ');
});
</script>
