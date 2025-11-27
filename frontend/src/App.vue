<template>
  <nav class="navbar">
    <div class="container navbar-content">
      <router-link to="/" class="brand">{{ $t('nav.blog_name') }}</router-link>
      <div class="nav-links">
        <div class="nav-search" ref="searchContainer">
          <input 
            v-model="searchQuery" 
            @input="onSearchInput"
            @focus="showResults = true"
            @keyup.enter="handleSearch"
            type="text" 
            :placeholder="$t('common.search_placeholder')"
          />
          <!-- Search Results Dropdown -->
          <div v-if="showResults && searchQuery" class="search-results fade-in">
            <div v-if="isSearching" class="search-loading">{{ $t('common.loading') }}</div>
            <div v-else-if="searchResults.length === 0" class="search-no-results">
              No results found
            </div>
            <div v-else class="search-list">
              <div 
                v-for="article in searchResults" 
                :key="article.id" 
                class="search-item"
                @click="goToArticle(article.id)"
              >
                <div class="search-item-header">
                  <div class="search-item-title">{{ article.title }}</div>
                  <div class="search-item-date">{{ formatDate(article.created_at) }}</div>
                </div>
                <div class="search-item-meta">
                  <span v-if="article.category" class="search-category">{{ article.category.name }}</span>
                  <span v-for="tag in article.tags" :key="tag.id" class="search-tag">#{{ tag.name }}</span>
                </div>
                <div class="search-item-excerpt" v-html="getExcerpt(article.content)"></div>
              </div>
            </div>
          </div>
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
import { ref, onMounted, onUnmounted } from 'vue';
import { useAuthStore } from './stores/auth';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import Toast from './components/Toast.vue';
import api from './api/axios';
import DOMPurify from 'dompurify';

const authStore = useAuthStore();
const router = useRouter();
const { locale } = useI18n();
const searchQuery = ref('');
const searchResults = ref([]);
const showResults = ref(false);
const isSearching = ref(false);
const searchContainer = ref(null);

let searchTimeout;

const logout = () => {
  authStore.logout();
  router.push('/login');
};

const handleSearch = () => {
  showResults.value = false;
  if (searchQuery.value.trim()) {
    // Preserve existing filters (category_id, tag_id) if we are on a page that supports them
    const currentQuery = { ...router.currentRoute.value.query };
    currentQuery.search = searchQuery.value;
    
    // If we are not on home or timeline, go to home with query
    if (router.currentRoute.value.path !== '/' && router.currentRoute.value.path !== '/timeline') {
      router.push({ path: '/', query: { search: searchQuery.value } });
    } else {
      // Otherwise update current route query
      router.push({ query: currentQuery });
    }
  } else {
    // Clear search but keep other filters
    const currentQuery = { ...router.currentRoute.value.query };
    delete currentQuery.search;
    router.push({ query: currentQuery });
  }
};

const onSearchInput = () => {
  showResults.value = true;
  if (!searchQuery.value.trim()) {
    searchResults.value = [];
    return;
  }

  isSearching.value = true;
  clearTimeout(searchTimeout);
  
  searchTimeout = setTimeout(async () => {
    try {
      // Pass current route filters to search context if needed, 
      // but for global search dropdown usually we want global results unless specified.
      // However, user asked for "search only in this category/tag" context.
      // Let's respect the current route filters for the dropdown too if we are on a filtered page.
      
      const params = { 
        search: searchQuery.value, 
        limit: 5 
      };

      // If we are on home or timeline, respect the current filters
      if (router.currentRoute.value.path === '/' || router.currentRoute.value.path === '/timeline') {
        if (router.currentRoute.value.query.category_id) {
          params.category_id = router.currentRoute.value.query.category_id;
        }
        if (router.currentRoute.value.query.tag_id) {
          params.tag_id = router.currentRoute.value.query.tag_id;
        }
      }

      const response = await api.get('/articles', { params });
      searchResults.value = response.data;
    } catch (err) {
      console.error('Search failed', err);
    } finally {
      isSearching.value = false;
    }
  }, 300);
};

const goToArticle = (id) => {
  router.push(`/articles/${id}`);
  showResults.value = false;
  searchQuery.value = '';
};

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString();
};

const getExcerpt = (content) => {
  if (!content) return '';
  
  // Strip markdown/HTML roughly to get plain text
  let text = content.replace(/[#*`]/g, '')
                    .replace(/!\[.*?\]\(.*?\)/g, '') // Remove images
                    .replace(/\[.*?\]\(.*?\)/g, '$1') // Keep link text
                    .replace(/<[^>]*>/g, ''); // Remove HTML tags

  const query = searchQuery.value;

  if (query && query.trim()) {
    const lowerText = text.toLowerCase();
    const lowerQuery = query.toLowerCase();
    const index = lowerText.indexOf(lowerQuery);

    if (index !== -1) {
      // Found keyword
      const start = Math.max(0, index - 30);
      const end = Math.min(text.length, index + query.length + 60);
      let snippet = text.substring(start, end);

      // Add ellipses
      if (start > 0) snippet = '...' + snippet;
      if (end < text.length) snippet = snippet + '...';

      // Escape HTML before adding highlight tags to prevent XSS
      snippet = DOMPurify.sanitize(snippet, { ALLOWED_TAGS: [] });

      // Highlight keyword (case-insensitive replacement)
      const regex = new RegExp(`(${query.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')})`, 'gi');
      snippet = snippet.replace(regex, '<mark class="highlight">$1</mark>');
      
      return snippet;
    }
  }

  // Default behavior if no keyword found in first part (or empty query)
  const snippet = text.length > 80 ? text.substring(0, 80) + '...' : text;
  return DOMPurify.sanitize(snippet, { ALLOWED_TAGS: [] });
};

const switchLang = (lang) => {
  locale.value = lang;
};

// Close search results when clicking outside
const handleClickOutside = (event) => {
  if (searchContainer.value && !searchContainer.value.contains(event.target)) {
    showResults.value = false;
  }
};

onMounted(() => {
  document.addEventListener('click', handleClickOutside);
});

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside);
});
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

.nav-search {
  position: relative;
}

.nav-search input {
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid var(--surface-border);
  color: var(--text-primary);
  padding: 0.4rem 0.8rem;
  border-radius: 4px;
  font-size: 0.9rem;
  width: 250px;
  transition: all 0.3s ease;
}

.nav-search input:focus {
  border-color: var(--primary-color);
  box-shadow: 0 0 10px rgba(0, 243, 255, 0.1);
  outline: none;
}

.search-results {
  position: absolute;
  top: 100%;
  left: 0;
  width: 400px; /* Wider to accommodate excerpt */
  background: rgba(5, 5, 16, 0.95);
  border: 1px solid var(--primary-color);
  box-shadow: 0 0 20px rgba(0, 243, 255, 0.2);
  border-radius: 4px;
  margin-top: 5px;
  max-height: 500px;
  overflow-y: auto;
  z-index: 1000;
}

.search-loading,
.search-no-results {
  padding: 1rem;
  text-align: center;
  color: var(--text-secondary);
  font-size: 0.9rem;
}

.search-item {
  padding: 0.75rem 1rem;
  border-bottom: 1px solid var(--surface-border);
  cursor: pointer;
  transition: background 0.2s;
}

.search-item:last-child {
  border-bottom: none;
}

.search-item:hover {
  background: rgba(0, 243, 255, 0.1);
}

.search-item-header {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  margin-bottom: 0.25rem;
}

.search-item-title {
  color: var(--text-primary);
  font-family: var(--font-heading);
  font-size: 0.95rem;
  font-weight: bold;
}

.search-item-date {
  color: var(--text-secondary);
  font-size: 0.75rem;
  white-space: nowrap;
  margin-left: 0.5rem;
}

.search-item-meta {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
  margin-bottom: 0.5rem;
  font-size: 0.75rem;
}

.search-category {
  color: var(--secondary-color);
  border: 1px solid var(--secondary-color);
  padding: 0 4px;
  border-radius: 2px;
}

.search-tag {
  color: var(--text-secondary);
  background: rgba(255, 255, 255, 0.1);
  padding: 0 4px;
  border-radius: 2px;
}

.search-item-excerpt {
  color: var(--text-secondary);
  font-size: 0.85rem;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* Highlight style */
:deep(mark.highlight) {
  background-color: rgba(0, 243, 255, 0.2);
  color: var(--primary-color);
  padding: 0 2px;
  border-radius: 2px;
  font-weight: bold;
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
