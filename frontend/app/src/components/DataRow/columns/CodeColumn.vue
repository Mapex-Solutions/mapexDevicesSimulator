<template>
  <code
      :class="props.column.ellipsis ? 'ellipsis' : ''"
      class="text-body2 text-weight-medium text-grey-8"
      style="font-family: 'Courier New', monospace; display: block; min-width: 0;"
  >
    {{ displayValue }}
    <AppTooltip v-if="props.column.ellipsis" :content="displayValue" />
  </code>
</template>

<script setup lang="ts">
/** TYPE IMPORTS */
import type { DataRowColumnProps } from '../interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPONENTS */
import { AppTooltip } from '@components/AppTooltip';

defineOptions({
  name: 'CodeColumn'
});

const props = defineProps<DataRowColumnProps>();

const displayValue = computed(() => {
  if (props.column.format) {
    return props.column.format(props.value, props.row);
  }
  return (props.value as string) || '-';
});
</script>
