<template>
	<transition name="splash-fade">
		<div v-if="visible" class="splash">
			<img class="splash__watermark" src="/iot-illustration.png" alt="" aria-hidden="true" />

			<div class="splash__content column items-center">
				<div class="splash__tile">
					<q-img src="/only-logo.png" width="84px" height="84px" fit="contain" no-spinner />
				</div>

				<h1 class="splash__title">{{ t('app.name') }}</h1>
				<p class="splash__tagline">{{ t('splash.tagline') }}</p>

				<q-linear-progress
					indeterminate
					rounded
					color="white"
					track-color="rgba(255,255,255,0.25)"
					class="splash__loader"
				/>
				<span class="splash__loading">{{ t('splash.loading') }}</span>
			</div>
		</div>
	</transition>
</template>

<script setup lang="ts">
/** VUE IMPORTS */
import { onMounted, ref } from 'vue';

/** COMPOSABLES */
import { useTranslations } from '@composables/i18n';

/** LOCAL IMPORTS */
import { SPLASH_DURATION_MS } from './constants';

/** COMPOSABLES & STORES */
const { t } = useTranslations();

/** STATE */
const visible = ref(true);

/** LIFECYCLE HOOKS */
onMounted(() => {
	setTimeout(() => {
		visible.value = false;
	}, SPLASH_DURATION_MS);
});
</script>

<style scoped lang="scss">
// Brand-fixed splash (theme-independent) shown once on boot before the app.
.splash {
	position: fixed;
	inset: 0;
	z-index: 9999;
	display: flex;
	align-items: center;
	justify-content: center;
	overflow: hidden;
	background: linear-gradient(150deg, var(--mapex-brand-grad-from) 0%, var(--mapex-brand-grad-mid) 45%, var(--mapex-brand-grad-to) 100%);
}

.splash__watermark {
	position: absolute;
	right: -40px;
	bottom: -40px;
	width: 520px;
	max-width: 70vw;
	opacity: 0.08;
	pointer-events: none;
	user-select: none;
}

.splash__content {
	position: relative;
	gap: var(--mapex-spacing-md);
	text-align: center;
	padding: var(--mapex-spacing-xl);
}

.splash__tile {
	width: 132px;
	height: 132px;
	border-radius: 30px;
	background: var(--mapex-on-brand);
	display: flex;
	align-items: center;
	justify-content: center;
	box-shadow: 0 18px 50px rgba(0, 0, 0, 0.30);
	animation: splash-pulse 1.8s ease-in-out infinite;
}

.splash__title {
	margin: var(--mapex-spacing-sm) 0 0;
	font-size: 1.6rem;
	font-weight: var(--mapex-font-weight-semibold);
	color: var(--mapex-on-brand);
	letter-spacing: 0.3px;
}

.splash__tagline {
	margin: 0;
	max-width: 360px;
	color: var(--mapex-on-brand-muted);
	font-size: var(--mapex-font-md);
}

.splash__loader {
	width: 200px;
	margin-top: var(--mapex-spacing-lg);
}

.splash__loading {
	color: var(--mapex-on-brand-subtle);
	font-size: var(--mapex-font-xs);
	letter-spacing: 0.4px;
	text-transform: uppercase;
}

@keyframes splash-pulse {
	0%, 100% { transform: scale(1); box-shadow: 0 18px 50px rgba(0, 0, 0, 0.30); }
	50% { transform: scale(1.06); box-shadow: 0 22px 60px rgba(0, 0, 0, 0.38); }
}

.splash-fade-leave-active {
	transition: opacity 0.5s ease;
}

.splash-fade-leave-to {
	opacity: 0;
}
</style>
