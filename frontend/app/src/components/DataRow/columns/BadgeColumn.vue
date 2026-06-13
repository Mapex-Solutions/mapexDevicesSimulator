<template>
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
</template>

<script setup lang="ts">
/** TYPE IMPORTS */
import type { DataRowColumnProps } from '../interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPONENTS */
import { AppTooltip } from '@components/AppTooltip';

defineOptions({
  name: 'BadgeColumn'
});

const props = defineProps<DataRowColumnProps>();

const displayValue = computed(() => {
  if (props.column.format) {
    return props.column.format(props.value, props.row);
  }
  return (props.value as string) || 'N/A';
});

function getColor() {
  if (typeof props.column.color === 'function') {
    return props.column.color(props.value, props.row);
  }
  return props.column.color || 'grey';
}
</script>
