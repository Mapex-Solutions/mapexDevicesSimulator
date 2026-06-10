/** TYPE IMPORTS */
import type { Ref } from 'vue';
import type { Gateway } from '@services/sim';

/** VUE IMPORTS */
import { computed } from 'vue';

/**
 * Derived reads over the gateway list.
 * @param {{ items: Ref<Gateway[]> }} ctx - reactive state bag
 */
export function createGatewaysGetters(ctx: { items: Ref<Gateway[]> }) {
	const count = computed(() => ctx.items.value.length);

	function getById(id: string): Gateway | undefined {
		return ctx.items.value.find((item) => item.id === id);
	}

	return { count, getById };
}
