/** TYPE IMPORTS */
import type { Ref } from 'vue';
import type { ConsoleMessage } from './messages.types';

/** VUE IMPORTS */
import { computed } from 'vue';

/**
 * Derived reads over the message stream: the active filter and the selected
 * message.
 * @param {object} ctx - reactive state bag
 */
export function createMessagesGetters(ctx: {
	items: Ref<ConsoleMessage[]>;
	selectedId: Ref<string | null>;
	deviceFilter: Ref<string | null>;
}) {
	const filtered = computed(() => {
		if (!ctx.deviceFilter.value) return ctx.items.value;
		return ctx.items.value.filter((message) => message.deviceId === ctx.deviceFilter.value);
	});

	const selected = computed(() => ctx.items.value.find((message) => message.id === ctx.selectedId.value) ?? null);

	const count = computed(() => ctx.items.value.length);

	return { filtered, selected, count };
}
