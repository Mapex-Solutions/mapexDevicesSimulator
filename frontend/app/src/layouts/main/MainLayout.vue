<script setup lang="ts">
/** VUE IMPORTS */
import { computed, ref } from 'vue';

/** COMPONENTS */
import AppHeader from './components/AppHeader.vue';
import AppSidebar from './components/AppSidebar.vue';

/** COMPOSABLES */
import { useTranslations } from '@composables/i18n';
import { useTheme } from '@composables/theme';

/** LOCAL IMPORTS */
import { buildMenuList } from './constants';

const githubUrl = 'https://github.com/Mapex-Solutions/mapexDevicesSimulator';

/** COMPOSABLES & STORES */
const { t } = useTranslations();
useTheme();

/** STATE */
const drawerOpen = ref(false);
const miniState = ref(true);

/** COMPUTED */
const menuList = computed(() => buildMenuList((key) => t(key)));

/** FUNCTIONS */

/**
 * Toggle the drawer: collapse/expand on desktop, open/close on small screens.
 */
function toggleDrawer(): void {
  if (window.innerWidth <= 1024) {
    drawerOpen.value = !drawerOpen.value;
  } else {
    miniState.value = !miniState.value;
  }
}
</script>

<template>
  <q-layout view="lHh Lpr lFf">
    <AppHeader @toggle-drawer="toggleDrawer" />

    <AppSidebar
      v-model:is-open="drawerOpen"
      :mini-state="miniState"
      :menu-list="menuList"
    />

    <q-page-container id="page-container" class="q-my-md">
      <div class="container">
        <router-view />
      </div>
    </q-page-container>

    <q-footer class="app-footer">
      <div class="app-footer__inner">
        {{ t('footer.prefix') }}
        <q-icon name="mdi-heart" class="app-footer__heart" />
        {{ t('footer.suffix') }}
        <a :href="githubUrl" target="_blank" rel="noopener noreferrer" class="app-footer__link">
          <q-icon name="mdi-github" size="16px" />
          Mapex Solutions · Thiago Anselmo
        </a>
      </div>
    </q-footer>
  </q-layout>
</template>

<style lang="scss">
.container {
  max-width: 1280px;
  width: 100%;
  margin: 0 auto;
  padding: 0 var(--mapex-spacing-lg);
}

@media (min-width: 1024px) {
  .container {
    padding: 0 var(--mapex-spacing-2xl);
  }
}

@media (min-width: 1440px) {
  .container {
    padding: 0 var(--mapex-spacing-3xl);
  }
}

body {
  background-color: transparent;
}

:deep(.menu-parent-item) {
  border-radius: 0 var(--mapex-radius-md) var(--mapex-radius-md) 0;
  margin: var(--mapex-spacing-xs) 0;
  transition: var(--mapex-transition-base);
  padding: var(--mapex-spacing-sm) var(--mapex-spacing-lg);

  &:hover {
    background: rgba(var(--q-primary-rgb), 0.05);
  }
}

:deep(.q-expansion-item__content) {
  background: var(--mapex-submenu-bg);
}

:deep(.q-item__section--side > .q-icon) {
  font-size: 24px;
}

:deep(.q-expansion-item__toggle-icon) {
  font-size: 24px;
}

.app-footer {
  background: var(--mapex-header-bg);
  color: var(--mapex-text-secondary);
  border-top: 1px solid var(--mapex-header-border);
  backdrop-filter: blur(8px);
}

.app-footer__inner {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  height: 36px;
  font-size: var(--mapex-font-xs);
}

.app-footer__heart {
  color: #e25555;
  animation: footer-heart 1.6s ease-in-out infinite;
}

.app-footer__link {
  display: inline-flex;
  align-items: center;
  gap: var(--mapex-spacing-xs);
  color: var(--mapex-primary);
  font-weight: var(--mapex-font-weight-medium);
  text-decoration: none;

  &:hover {
    text-decoration: underline;
  }
}

@keyframes footer-heart {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.18); }
}
</style>
