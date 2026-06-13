<template>
  <div style="min-width: 0; display: flex; flex-direction: column; gap: var(--mapex-spacing-2xs);">
    <!-- Primary Text -->
    <div
        :class="props.column.ellipsis ? 'ellipsis' : ''"
        :style="mobile ? 'min-width: 0; flex: 1;' : 'min-width: 0;'"
        class="text-body2 text-weight-medium"
    >
      {{ displayValue }}
      <AppTooltip v-if="props.column.ellipsis" :content="displayValue" />
    </div>

    <!-- Secondary Text (if secondaryKey exists) -->
    <div
        v-if="secondaryValue"
        class="text-caption text-grey-6 ellipsis"
        style="min-width: 0;"
    >
      {{ secondaryValue }}
      <AppTooltip :content="secondaryValue" />
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
  name: 'TextColumn'
});

const props = defineProps<DataRowColumnProps>();

const displayValue = computed(() => {
  if (props.column.format) {
    return props.column.format(props.value, props.row);
  }
  return (props.value as string) || '-';
});

const secondaryValue = computed(() => {
  // Support secondary as function (like format)
  if (props.column.secondary) {
    return props.column.secondary(props.value, props.row);
  }

  // Support secondaryKey as property path
  if (props.column.secondaryKey) {
    const keys = props.column.secondaryKey.split('.');
    let value: unknown = props.row;

    for (const key of keys) {
      value = (value as Record<string, unknown> | undefined)?.[key];
    }

    return (value as string) || null;
  }

  return null;
});
</script>
