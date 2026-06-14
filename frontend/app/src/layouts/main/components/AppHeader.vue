<script setup lang="ts">
/** VUE IMPORTS */
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

/** COMPOSABLES */
import { useTheme } from '@composables/theme';
import { useBreadcrumbs } from '../composables';

/** UTILS */
import { notifyInfo } from '@utils/alert';

/** LOCAL IMPORTS */
import { buildMenuList } from '../constants';

/** EMITS */
const emit = defineEmits<{
  'toggle-drawer': [];
}>();

/** COMPOSABLES & STORES */
const { t, locale } = useI18n();
const { isDark, toggleTheme } = useTheme();
const translatedMenu = computed(() => buildMenuList((key) => t(key)));
const { breadcrumbs } = useBreadcrumbs(translatedMenu);

/** COMPUTED */

/**
 * Language options with reactive labels
 */
const languageList = computed(() => [
  { value: 'en-US', label: t('language.english'), icon: '/flag-usa.svg' },
  { value: 'pt-BR', label: t('language.portuguese'), icon: '/flag-br.svg' },
]);

/**
 * Current locale
 */
const currentLocale = computed(() => locale.value);

/** FUNCTIONS */

/**
 * Change application language and persist to localStorage
 *
 * @param {string} lang - Locale code (e.g. 'en-US', 'pt-BR')
 */
function changeLanguage(lang: string): void {
  locale.value = lang;
  localStorage.setItem('user-locale', lang);
  const label = languageList.value.find(l => l.value === lang)?.label;
  notifyInfo({ message: `${t('language.changed')} ${label}` });
}
</script>

<template>
  <q-header
    elevated
    class="header-container"
  >
    <q-toolbar class="q-px-lg q-py-sm">
      <!-- Menu Toggle and Separator -->
      <div class="row items-center">
        <q-btn
          flat
          dense
          round
          no-wrap
          class="q-mr-sm"
          icon="menu"
          color="primary"
          aria-label="Menu"
          @click="emit('toggle-drawer')"
        />
        <q-separator
          vertical
          class="q-mx-sm"
        />
      </div>

      <!-- Breadcrumbs -->
      <q-breadcrumbs
        id="header-breadcrumbs"
        class="text-primary"
        active-color="primary"
        separator-color="grey-6"
      >
        <q-breadcrumbs-el
          icon="dashboard"
          to="/"
          :label="breadcrumbs.length === 0 ? t('layout.home') : undefined"
        />
        <q-breadcrumbs-el
          v-for="(crumb, index) in breadcrumbs"
          :key="index"
          :label="crumb.label"
          :icon="crumb.icon"
          :to="crumb.to"
        />
      </q-breadcrumbs>

      <q-space/>

      <!-- Action Buttons -->
      <div class="row items-center q-gutter-sm">

        <!-- Theme Toggle -->
        <q-btn
          id="header-theme-toggle"
          flat
          round
          :icon="isDark ? 'light_mode' : 'dark_mode'"
          :color="isDark ? 'amber' : 'primary'"
          @click="toggleTheme"
        >
          <q-tooltip>{{ isDark ? t('layout.toLight') : t('layout.toDark') }}</q-tooltip>
        </q-btn>

        <!-- Language Selector -->
        <q-btn
          id="header-lang-selector"
          flat
          round
          color="primary"
        >
          <q-icon name="language"/>
          <q-menu
            :class="isDark ? 'bg-dark' : 'bg-white'"
            anchor="bottom right"
            self="top right"
          >
            <q-list style="min-width: 150px">
              <q-item
                v-for="lang in languageList"
                v-close-popup
                clickable
                active-class="bg-primary text-white"
                :key="lang.value"
                :active="currentLocale === lang.value"
                @click="changeLanguage(lang.value)"
              >
                <q-item-section avatar>
                  <q-img
                    width="24px"
                    :src="lang.icon"
                  />
                </q-item-section>
                <q-item-section>{{ lang.label }}</q-item-section>
              </q-item>
            </q-list>
          </q-menu>
        </q-btn>
      </div>
    </q-toolbar>
  </q-header>
</template>

<style lang="scss" scoped>
.header-container {
  background: var(--mapex-header-bg);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid var(--mapex-header-border);
}
</style>
