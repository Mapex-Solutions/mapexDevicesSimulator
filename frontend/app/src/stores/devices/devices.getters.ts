/** TYPE IMPORTS */
import type { Ref } from 'vue';
import type { Device, ProtocolId } from '@services/sim';

/** VUE IMPORTS */
import { computed } from 'vue';

/**
 * Derived reads over the device list.
 * @param {{ items: Ref<Device[]> }} ctx - reactive state bag
 */
export function createDevicesGetters(ctx: { items: Ref<Device[]> }) {
	const count = computed(() => ctx.items.value.length);

	function byProtocol(protocolId: ProtocolId): Device[] {
		return ctx.items.value.filter((item) => item.protocolId === protocolId);
	}

	function getById(id: string): Device | undefined {
		return ctx.items.value.find((item) => item.id === id);
	}

	return { count, byProtocol, getById };
}
