<template>
	<q-card flat class="marketplace-card column" @click="emit('open', item)">
		<!-- Header: device icon, model and vendor -->
		<q-card-section class="row items-center no-wrap q-pb-sm">
			<div class="card-avatar row items-center justify-center q-mr-md">
				<q-img v-if="imageUrl" :src="imageUrl" fit="contain" class="card-photo" no-spinner>
					<template #error>
						<q-icon :name="item.icon || 'mdi-chip'" size="24px" color="primary" />
					</template>
				</q-img>
				<q-icon v-else :name="item.icon || 'mdi-chip'" size="24px" color="primary" />
			</div>
			<div class="col">
				<div class="text-subtitle2 text-weight-semibold ellipsis">{{ item.model }}</div>
				<div class="text-caption text-grey-7 ellipsis">{{ item.vendorName }}</div>
			</div>
			<q-chip dense square size="sm" color="primary" text-color="white" class="q-ml-sm">
				{{ item.protocolLabel }}
			</q-chip>
		</q-card-section>

		<!-- Description -->
		<q-card-section class="q-py-none col">
			<div class="text-body2 card-description">{{ item.description }}</div>
		</q-card-section>

		<!-- Reading types -->
		<q-card-section class="q-pt-sm q-pb-none">
			<div class="row items-center q-gutter-xs">
				<q-chip
					v-for="reading in item.readingMetas"
					:key="reading.value"
					dense
					outline
					size="sm"
					color="primary"
					:icon="reading.icon"
				>
					{{ reading.label }}
				</q-chip>
			</div>
		</q-card-section>

		<!-- Footer: resource badges + add -->
		<q-card-actions class="q-px-md q-pb-md q-pt-sm items-center">
			<q-icon v-if="item.hasManual" name="menu_book" size="xs" color="grey-6" class="q-mr-sm">
				<AppTooltip :content="manualLabel" />
			</q-icon>
			<q-icon v-if="item.hasCodec" name="data_object" size="xs" color="grey-6">
				<AppTooltip :content="codecLabel" />
			</q-icon>
			<q-space />
			<q-btn
				flat
				dense
				no-caps
				color="primary"
				icon="info"
				:label="detailsLabel"
				class="q-mr-xs"
				@click.stop="emit('open', item)"
			/>
			<q-btn
				unelevated
				dense
				no-caps
				color="primary"
				icon="add"
				:label="addLabel"
				:loading="installing"
				@click.stop="emit('install', item)"
			/>
		</q-card-actions>
	</q-card>
</template>

<script setup lang="ts">
/** TYPE IMPORTS */
import type { MarketplaceCardItem } from '../interfaces/marketplaceListPage.interface';

/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPONENTS */
import { AppTooltip } from '@components/AppTooltip';

/** SERVICES */
import { resolveMarketplaceAssetUrl } from '@services/sim';

/** PROPS & EMITS */
const props = defineProps<{
	item: MarketplaceCardItem;
	installing: boolean;
	addLabel: string;
	detailsLabel: string;
	manualLabel: string;
	codecLabel: string;
}>();

/** Absolute URL of the device photo, when the catalog item carries one. */
const imageUrl = computed((): string =>
	props.item.image ? resolveMarketplaceAssetUrl(props.item.vendor, props.item.slug, props.item.image) : '',
);

const emit = defineEmits<{
	open: [item: MarketplaceCardItem];
	install: [item: MarketplaceCardItem];
}>();
</script>

<style scoped lang="scss">
.marketplace-card {
	height: 100%;
	background: var(--mapex-surface-elevated);
	border: 1px solid var(--mapex-card-border);
	border-radius: var(--mapex-radius-lg);
	box-shadow: 0 2px 4px var(--mapex-elevation-shadow);
	transition: var(--mapex-transition-base);
	cursor: pointer;

	&:hover {
		border-color: var(--mapex-card-hover-border);
		box-shadow: 0 4px 12px var(--mapex-hover-shadow);
		transform: translateY(-2px);
	}
}

.card-avatar {
	width: 44px;
	height: 44px;
	border-radius: var(--mapex-radius-md);
	background: rgba(var(--mapex-primary-rgb), 0.1);
	flex-shrink: 0;
	overflow: hidden;
}

.card-photo {
	width: 100%;
	height: 100%;
}

.card-description {
	color: var(--mapex-text-secondary);
	display: -webkit-box;
	-webkit-line-clamp: 2;
	line-clamp: 2;
	-webkit-box-orient: vertical;
	overflow: hidden;
	min-height: 2.6em;
}
</style>
