<template>
	<router-view />
	<SplashScreen />
</template>

<script setup lang="ts">
/** COMPONENTS */
import { SplashScreen } from '@components/SplashScreen';

/**
 * Root component. Intentionally minimal: theme is initialized by useTheme inside
 * MainLayout, sidecar health is polled by the pages that show it, and the splash
 * overlays the app once on boot then fades itself out.
 */
</script>
