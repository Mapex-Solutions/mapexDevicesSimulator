/** TYPE IMPORTS */
import type { Log, LogQuery } from '@services/sim';

/** VUE IMPORTS */
import { computed, ref } from 'vue';

/** SERVICES */
import { defineStore } from 'pinia';
import { sim } from '@services/sim';

/**
 * Persisted logs: the SQLite-backed history of device messages. Reads call the
 * engine's paginated/filterable endpoint and reflect it exactly; an offline
 * engine yields an empty page (connectivity is shown globally via health).
 */
export const useLogsStore = defineStore('logs', () => {
	/** STATE */
	const items = ref<Log[]>([]);
	const total = ref(0);
	const loading = ref(false);
	const page = ref(1);
	const itemsPerPage = ref(15);
	const filters = ref<Record<string, string>>({});
	const lastUpdatedAt = ref<number | undefined>(undefined);

	/** COMPUTED */
	const totalPages = computed(() => Math.max(1, Math.ceil(total.value / itemsPerPage.value)));

	async function fetch(): Promise<void> {
		loading.value = true;
		try {
			const query: LogQuery = {
				limit: itemsPerPage.value,
				offset: (page.value - 1) * itemsPerPage.value,
				...filters.value,
			};
			const res = await sim.logs.list(query);
			items.value = res.items;
			total.value = res.total;
		} catch {
			items.value = [];
			total.value = 0;
		} finally {
			loading.value = false;
			lastUpdatedAt.value = Date.now();
		}
	}

	function setPage(value: number): void {
		page.value = value;
		void fetch();
	}

	function setItemsPerPage(value: number): void {
		itemsPerPage.value = value;
		page.value = 1;
		void fetch();
	}

	function setFilters(value: Record<string, string>): void {
		filters.value = value;
		page.value = 1;
		void fetch();
	}

	return {
		items, total, loading, page, itemsPerPage, filters, totalPages, lastUpdatedAt,
		fetch, setPage, setItemsPerPage, setFilters,
	};
});
