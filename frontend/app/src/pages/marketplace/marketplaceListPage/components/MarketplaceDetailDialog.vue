<template>
	<q-dialog
		:model-value="modelValue"
		transition-show="jump-up"
		transition-hide="jump-down"
		@update:model-value="emit('update:modelValue', $event)"
	>
		<q-card class="detail-dialog column no-wrap">
			<!-- Header -->
			<div class="detail-header row items-center no-wrap q-pa-md">
				<div class="header-thumb row items-center justify-center q-mr-md">
					<q-img v-if="imageUrl" :src="imageUrl" fit="contain" class="fit" no-spinner />
					<q-icon v-else :name="item?.icon || 'mdi-chip'" size="28px" color="primary" />
				</div>
				<div class="col">
					<div class="text-subtitle1 text-weight-bold ellipsis">{{ item?.name }}</div>
					<div class="text-caption text-grey-7">{{ item?.vendorName }} · {{ item?.model }}</div>
				</div>
				<q-btn v-close-popup flat round dense icon="close" color="grey-7" @click="emit('update:modelValue', false)" />
			</div>

			<q-separator />

			<!-- Tabs -->
			<q-tabs
				v-model="tab"
				dense
				no-caps
				active-color="primary"
				indicator-color="primary"
				align="left"
				class="detail-tabs text-grey-7"
			>
				<q-tab name="overview" icon="info" :label="t('marketplace.tabOverview')" />
				<q-tab name="codecs" icon="data_object" :label="t('marketplace.tabCodecs')">
					<q-badge v-if="codecs.length" color="primary" floating>{{ codecs.length }}</q-badge>
				</q-tab>
				<q-tab name="files" icon="folder_open" :label="t('marketplace.tabFiles')" />
			</q-tabs>

			<q-separator />

			<!-- Body -->
			<div class="detail-body col">
				<div v-if="loading" class="column items-center justify-center full-height q-py-xl">
					<q-spinner color="primary" size="2.5em" />
				</div>

				<q-tab-panels v-else v-model="tab" animated class="fit">
					<!-- Overview -->
					<q-tab-panel name="overview" class="q-pa-lg">
						<div class="overview-image row items-center justify-center q-mb-lg">
							<q-img v-if="imageUrl" :src="imageUrl" fit="contain" class="fit" no-spinner />
							<q-icon v-else :name="item?.icon || 'mdi-chip'" size="72px" color="primary" />
						</div>

						<div class="row items-center q-gutter-xs q-mb-md">
							<q-chip dense square size="sm" color="primary" text-color="white">{{ item?.protocolLabel }}</q-chip>
							<q-chip
								v-for="reading in item?.readingMetas || []"
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

						<div class="text-body2 text-grey-8 detail-description">{{ description }}</div>

						<div v-if="info?.tags?.length" class="row items-center q-gutter-xs q-mt-md">
							<q-chip v-for="tag in info?.tags" :key="tag" dense size="sm" outline color="grey-7">{{ tag }}</q-chip>
						</div>

						<q-btn
							v-if="info?.vendor.site"
							flat
							dense
							no-caps
							size="sm"
							color="primary"
							icon="open_in_new"
							:label="t('marketplace.site')"
							class="q-mt-md q-pl-none"
							@click="openExternal(info?.vendor.site)"
						/>
						<div v-if="collectedLabel" class="text-caption text-grey-6 q-mt-md">
							{{ t('marketplace.collectedAt') }}: {{ collectedLabel }}
						</div>
					</q-tab-panel>

					<!-- Codecs -->
					<q-tab-panel name="codecs" class="q-pa-lg">
						<div class="text-caption text-grey-7 q-mb-md">{{ t('marketplace.codecsHint') }}</div>

						<div v-if="!codecs.length" class="text-body2 text-grey-6">{{ t('marketplace.codecNone') }}</div>

						<div v-for="codec in codecs" :key="codec.id" class="codec-card q-pa-md q-mb-sm" :class="{ 'codec-card--default': codec.default }">
							<div class="row items-center no-wrap">
								<div class="col">
									<div class="row items-center q-gutter-xs">
										<span class="text-body2 text-weight-medium">{{ codec.id === 'not-found' ? t('marketplace.codecNotFound') : codec.name }}</span>
										<q-chip v-if="codec.official" dense size="sm" color="positive" text-color="white" icon="verified">{{ t('marketplace.codecOfficial') }}</q-chip>
										<q-chip v-if="codec.default" dense size="sm" color="primary" text-color="white">{{ t('marketplace.codecDefault') }}</q-chip>
									</div>
									<div v-if="codec.id !== 'not-found'" class="text-caption text-grey-7 q-mt-xs">
										{{ t('marketplace.codecTarget') }}: {{ codec.target || '—' }} · {{ codec.source || 'community' }} · {{ codec.language }}
									</div>
								</div>
								<q-chip v-if="codec.id !== 'not-found'" dense square size="sm" color="grey-3" text-color="grey-9" class="q-ml-sm">{{ codec.target }}</q-chip>
							</div>
							<div class="row q-gutter-sm q-mt-sm">
								<q-btn
									v-if="codec.id !== 'not-found'"
									flat
									dense
									no-caps
									size="sm"
									color="primary"
									icon="visibility"
									:label="t('marketplace.codecView')"
									@click="openExternal(codecFileUrl(codec))"
								/>
								<q-btn
									v-if="codec.sourceUrl"
									flat
									dense
									no-caps
									size="sm"
									color="grey-8"
									icon="open_in_new"
									:label="t('marketplace.codecSource')"
									@click="openExternal(codec.sourceUrl)"
								/>
							</div>
						</div>
					</q-tab-panel>

					<!-- Files -->
					<q-tab-panel name="files" class="q-pa-lg">
						<div class="column q-gutter-sm">
							<q-btn
								v-if="datasheetDoc.present"
								outline
								no-caps
								align="left"
								color="primary"
								:icon="datasheetDoc.local ? 'description' : 'open_in_new'"
								:label="datasheetDoc.local ? t('marketplace.datasheet') : t('marketplace.datasheetOnline')"
								@click="openExternal(datasheetDoc.href)"
							/>
							<q-btn
								v-if="manualDoc.present"
								outline
								no-caps
								align="left"
								color="primary"
								:icon="manualDoc.local ? 'menu_book' : 'open_in_new'"
								:label="manualDoc.local ? t('marketplace.manual') : t('marketplace.manualOnline')"
								@click="openExternal(manualDoc.href)"
							/>
							<div v-if="!datasheetDoc.present && !manualDoc.present" class="text-body2 text-grey-6">
								{{ t('marketplace.noResources') }}
							</div>
						</div>
					</q-tab-panel>
				</q-tab-panels>
			</div>

			<q-separator />

			<!-- Footer -->
			<div class="detail-footer row items-center q-pa-md">
				<q-space />
				<q-btn v-close-popup flat no-caps :label="t('common.close')" class="q-mr-sm" />
				<q-btn
					unelevated
					no-caps
					color="primary"
					icon="add"
					:label="t('marketplace.addToDevices')"
					:loading="installing"
					:disable="!item"
					@click="item && emit('install', item)"
				/>
			</div>
		</q-card>
	</q-dialog>
</template>

<script setup lang="ts">
/** TYPE IMPORTS */
import type { MarketplaceInformation, MarketplaceCodec } from '@services/sim';
import type { MarketplaceCardItem } from '../interfaces/marketplaceListPage.interface';

/** VUE IMPORTS */
import { computed, ref, watch } from 'vue';

/** COMPOSABLES */
import { useTranslations } from '@composables/i18n';
import { useDateFormat } from '@composables/datetime';

/** SERVICES */
import { resolveMarketplaceAssetUrl } from '@services/sim';

/** PROPS & EMITS */
const props = defineProps<{
	modelValue: boolean;
	item: MarketplaceCardItem | null;
	info: MarketplaceInformation | null;
	codecs: MarketplaceCodec[];
	loading: boolean;
	installing: boolean;
}>();

const emit = defineEmits<{
	'update:modelValue': [value: boolean];
	install: [item: MarketplaceCardItem];
}>();

/** COMPOSABLES */
const { t, locale } = useTranslations();
const { formatDate } = useDateFormat();

/** STATE */
const tab = ref<'overview' | 'codecs' | 'files'>('overview');

/** COMPUTED */

/** Absolute URL of the device image, when the detail sheet declares one. */
const imageUrl = computed((): string => {
	if (!props.item || !props.info?.images.device) return '';
	return resolveMarketplaceAssetUrl(props.item.vendor, props.item.slug, props.info.images.device);
});

/** Datasheet/manual resolved to a local PDF asset URL if bundled, else the external link. */
const datasheetDoc = computed(() => resolveDoc(props.info?.files.datasheet, props.info?.files.datasheetUrl));
const manualDoc = computed(() => resolveDoc(props.info?.files.manual, props.info?.files.manualUrl));

/** Collection date formatted for display (empty when the entry has none). */
const collectedLabel = computed((): string => (props.info?.collectedAt ? formatDate(props.info.collectedAt) : ''));

/**
 * The detail description resolved for the active locale. `description` may be a
 * plain string (older curated entries) or a { locale: text } map; pick the active
 * locale, then en-US, then the card's server-resolved description.
 */
const description = computed((): string => {
	const value = props.info?.description;
	if (typeof value === 'string') return value;
	if (value) return value[locale.value] ?? value['en-US'] ?? '';
	return props.item?.description ?? '';
});

/** WATCHERS */

/** Reset to the first tab each time the dialog opens on a new device. */
watch(
	() => props.modelValue,
	(open) => {
		if (open) tab.value = 'overview';
	},
);

/** FUNCTIONS */

/**
 * Build the URL of a codec's decoder file, served as a bundle asset.
 * @param {MarketplaceCodec} codec - the codec to locate
 * @returns {string} the decoder file URL
 */
function codecFileUrl(codec: MarketplaceCodec): string {
	if (!props.item) return '';
	const file = codec.decoderFile || 'decoder.js';
	return resolveMarketplaceAssetUrl(props.item.vendor, props.item.slug, `${codec.path}/${file}`);
}

/**
 * Resolve a document for display: a local bundled PDF (served as an asset) takes
 * precedence; otherwise the external link (e.g. an online wiki manual) is used.
 * @param {string | undefined} localPath - bundled file path, when present
 * @param {string | undefined} externalUrl - external doc URL, when present
 * @returns {{ href: string; local: boolean; present: boolean }} the resolved doc
 */
function resolveDoc(
	localPath: string | undefined,
	externalUrl: string | undefined,
): { href: string; local: boolean; present: boolean } {
	if (localPath && props.item) {
		return { href: resolveMarketplaceAssetUrl(props.item.vendor, props.item.slug, localPath), local: true, present: true };
	}
	if (externalUrl) {
		return { href: externalUrl, local: false, present: true };
	}
	return { href: '', local: false, present: false };
}

/**
 * Open an external resource (vendor site, codec, datasheet, manual) in a new tab.
 * @param {string | undefined} url - the URL to open
 */
function openExternal(url: string | undefined): void {
	if (url) window.open(url, '_blank', 'noopener,noreferrer');
}
</script>

<style scoped lang="scss">
.detail-dialog {
	width: 720px;
	max-width: 92vw;
	height: 640px;
	max-height: 88vh;
	border-radius: var(--mapex-radius-xl);
	background: var(--mapex-surface-elevated);
	overflow: hidden;
}

.detail-header {
	background: var(--mapex-header-bg);

	.header-thumb {
		width: 52px;
		height: 52px;
		padding: 4px;
		border-radius: var(--mapex-radius-md);
		// Constant light plate so transparent device SVGs stay visible in dark mode.
		background: var(--mapex-device-media-bg);
		border: 1px solid var(--mapex-device-media-border);
		flex-shrink: 0;
		overflow: hidden;
	}
}

.detail-tabs {
	padding: 0 var(--mapex-spacing-sm);
}

.detail-body {
	min-height: 0;
	overflow-y: auto;
}

.overview-image {
	width: 100%;
	height: 200px;
	padding: var(--mapex-spacing-md);
	border-radius: var(--mapex-radius-lg);
	// Constant light plate so transparent device SVGs stay visible in dark mode.
	background: var(--mapex-device-media-bg);
	border: 1px solid var(--mapex-device-media-border);
	overflow: hidden;
}

.detail-description {
	line-height: var(--mapex-line-height-relaxed);
}

.codec-card {
	border: 1px solid var(--mapex-card-border);
	border-radius: var(--mapex-radius-md);
	background: var(--mapex-surface-bg);
	transition: var(--mapex-transition-base);

	&--default {
		border-color: var(--mapex-active-border);
		background: var(--mapex-active-bg);
	}
}

.detail-footer {
	background: var(--mapex-header-bg);
}
</style>
