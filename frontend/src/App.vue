<template>
  <nav class="navbar">
    <div class="container navbar-content">
      <router-link to="/" class="brand">{{ $t('nav.blog_name') }}</router-link>
      <div class="nav-links">
        <div class="nav-search">
          <input 
            v-model="searchQuery" 
            @keyup.enter="handleSearch"
            type="text" 
            :placeholder="$t('common.search_placeholder')"
          />
        </div>
        <router-link to="/">{{ $t('nav.home') }}</router-link>
        <router-link to="/timeline">{{ $t('timeline.title') }}</router-link>
        <template v-if="authStore.isAuthenticated">
          <router-link to="/articles/new">{{ $t('nav.write') }}</router-link>
          <router-link to="/profile">{{ $t('nav.profile') }}</router-link>
          <a href="#" @click.prevent="logout">{{ $t('nav.logout') }}</a>
        </template>
        <template v-else>
          <router-link to="/login">{{ $t('nav.login') }}</router-link>
          <router-link to="/register">{{ $t('nav.register') }}</router-link>
        </template>
        <div class="lang-switch">
          <button @click="switchLang('en')" :class="{ active: locale === 'en' }">EN</button>
          <span class="divider">/</span>
          <button @click="switchLang('zh')" :class="{ active: locale === 'zh' }">中文</button>
        </div>
      </div>
    </div>
  </nav>
  <main class="container main-content">
    <router-view v-slot="{ Component }">
      <transition name="fade" mode="out-in">
        <component :is="Component" />
      </transition>
    </router-view>
    <Toast />
  </main>
</template>

<script setup>
import { ref } from 'vue';
import { useAuthStore } from './stores/auth';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import Toast from './components/Toast.vue';

const authStore = useAuthStore();
const router = useRouter();
const { locale } = useI18n();
const searchQuery = ref('');

const logout = () => {
  authStore.logout();
  router.push('/login');
};

const handleSearch = () => {
  if (searchQuery.value.trim()) {
    router.push({ path: '/', query: { search: searchQuery.value } });
  } else {
    router.push('/');
  }
};

const switchLang = (lang) => {
  locale.value = lang;
};
</script>

<style scoped>
.navbar {
  background: rgba(5, 5, 16, 0.8);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid var(--surface-border);
  padding: 1rem 0;
  position: sticky;
  top: 0;
  z-index: 100;
}

.navbar-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.brand {
  font-family: var(--font-heading);
  font-size: 1.8rem;
  font-weight: 900;
  color: var(--primary-color);
  text-transform: uppercase;
  letter-spacing: 2px;
  text-shadow: var(--neon-text-shadow);
}

.nav-links {
  display: flex;
  gap: 1.5rem;
  align-items: center;
}

.nav-search input {
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid var(--surface-border);
  color: var(--text-primary);
  padding: 0.4rem 0.8rem;
  border-radius: 4px;
  font-size: 0.9rem;
  width: 200px;
  transition: all 0.3s ease;
}

.nav-search input:focus {
  border-color: var(--primary-color);
  box-shadow: 0 0 10px rgba(0, 243, 255, 0.1);
  outline: none;
}

.nav-links a {
  color: var(--text-secondary);
  font-family: var(--font-heading);
  font-size: 0.9rem;
  letter-spacing: 1px;
  text-transform: uppercase;
  white-space: nowrap;
}

.nav-links a:hover, .nav-links a.router-link-active {
  color: var(--primary-color);
  text-shadow: var(--neon-text-shadow);
}

.lang-switch {
  display: flex;
  align-items: center;
  margin-left: 0.5rem;
  padding-left: 1rem;
  border-left: 1px solid var(--surface-border);
}

.lang-switch button {
  background: none;
  border: none;
  padding: 0;
  font-family: var(--font-heading);
  font-size: 0.8rem;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.3s ease;
}

.lang-switch button.active {
  color: var(--secondary-color);
  text-shadow: 0 0 5px var(--secondary-glow);
}

.divider {
  margin: 0 0.5rem;
  color: var(--surface-border);
}

.main-content {
  padding-top: 2rem;
  padding-bottom: 2rem;
  flex: 1;
}

/* Route Transition */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(10px);
}
</style>
