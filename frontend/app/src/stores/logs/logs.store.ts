/** TYPE IMPORTS */
import type { Log, LogQuery } from '@services/sim';

/** VUE IMPORTS */
import { computed, ref } from 'vue';

/** SERVICES */
import { defineStore } from 'pinia';
import { sim } from '@services/sim';

/**
 * Persisted logs: the SQLite-backed history of device messages. Pagination is
 * cursor-based (keyset), so it stays stable while new rows arrive at the top and
 * never pays an offset/count cost on a large table. Previous is served from a
 * client-side stack of the cursors used to reach each page; Next follows the
 * cursor the engine returns. An offline engine yields an empty page.
 */
export const useLogsStore = defineStore('logs', () => {
	/** STATE */
	const items = ref<Log[]>([]);
	const loading = ref(false);
	const itemsPerPage = ref(15);
	const filters = ref<Record<string, string>>({});
	const lastUpdatedAt = ref<number | undefined>(undefined);

	// The first entry is the empty first-page cursor; each Next pushes the engine's
	// nextCursor, each Previous pops back to the page before it.
	const cursorStack = ref<string[]>(['']);
	const nextCursor = ref('');

	/** GETTERS */
	const hasPrev = computed(() => cursorStack.value.length > 1);
	const hasNext = computed(() => nextCursor.value !== '');

	/** ACTIONS */

	/**
	 * Load the page addressed by the cursor on top of the stack.
	 */
	async function load(): Promise<void> {
		loading.value = true;
		try {
			const cursor = cursorStack.value[cursorStack.value.length - 1] ?? '';
			const query: LogQuery = {
				limit: itemsPerPage.value,
				...(cursor ? { cursor } : {}),
				...filters.value,
			};
			const res = await sim.logs.list(query);
			items.value = res.items;
			nextCursor.value = res.nextCursor ?? '';
		} catch {
			items.value = [];
			nextCursor.value = '';
		} finally {
			loading.value = false;
			lastUpdatedAt.value = Date.now();
		}
	}

	/**
	 * Reload from the first page, resetting the cursor history.
	 */
	function fetch(): void {
		cursorStack.value = [''];
		void load();
	}

	/**
	 * Advance to the next page, if there is one.
	 */
	function next(): void {
		if (!nextCursor.value) return;
		cursorStack.value.push(nextCursor.value);
		void load();
	}

	/**
	 * Go back to the previous page, if there is one.
	 */
	function prev(): void {
		if (cursorStack.value.length <= 1) return;
		cursorStack.value.pop();
		void load();
	}

	/**
	 * Change the page size and reload from the first page.
	 * @param {number} value - the new page size
	 */
	function setItemsPerPage(value: number): void {
		itemsPerPage.value = value;
		fetch();
	}

	/**
	 * Replace the active filters and reload from the first page.
	 * @param {Record<string, string>} value - the filter map
	 */
	function setFilters(value: Record<string, string>): void {
		filters.value = value;
		fetch();
	}

	return {
		items, loading, itemsPerPage, filters, lastUpdatedAt, hasPrev, hasNext,
		fetch, next, prev, setItemsPerPage, setFilters,
	};
});
