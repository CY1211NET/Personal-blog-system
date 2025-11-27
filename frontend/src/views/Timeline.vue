<template>
  <div class="timeline-layout">
    <div class="content-area">
      <h1 class="page-title text-center">{{ $t('timeline.title') }}</h1>
      
      <!-- Active Filters Display -->
      <div v-if="hasActiveFilters" class="active-filters fade-in">
        <span v-if="route.query.search" class="filter-badge">
          🔍 {{ route.query.search }}
          <button @click="clearFilter('search')">×</button>
        </span>
        <span v-if="activeCategory" class="filter-badge">
          📂 {{ activeCategory.name }}
          <button @click="clearFilter('category_id')">×</button>
        </span>
        <span v-if="activeTag" class="filter-badge">
          🏷️ {{ activeTag.name }}
          <button @click="clearFilter('tag_id')">×</button>
        </span>
        <button v-if="hasActiveFilters" @click="clearAllFilters" class="clear-all-btn">{{ $t('common.cancel') }}</button>
      </div>

      <div v-if="loading" class="loading">{{ $t('common.loading') }}</div>
      
      <div v-else class="timeline">
        <div v-if="articles.length === 0" class="no-articles">
          No articles found matching filters.
        </div>
        <div v-for="(group, year) in groupedArticles" :key="year" class="timeline-year">
          <h2 class="year-title">{{ year }}</h2>
          <div class="timeline-items">
            <div v-for="article in group" :key="article.id" class="timeline-item">
              <div class="timeline-dot"></div>
              <div class="timeline-date">{{ formatDate(article.created_at) }}</div>
              <div class="timeline-content card">
                <router-link :to="'/articles/' + article.id" class="timeline-link">
                  <h3>{{ article.title }}</h3>
                </router-link>
                <div class="timeline-meta">
                  <span v-if="article.category" class="category-tag">{{ article.category.name }}</span>
                  <span v-for="tag in article.tags" :key="tag.id" class="tag-pill">#{{ tag.name }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <Sidebar @filter="handleFilter" class="sidebar-area fade-in" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import api from '../api/axios';
import Sidebar from '../components/Sidebar.vue';

const route = useRoute();
const router = useRouter();
const articles = ref([]);
const loading = ref(true);
const categories = ref([]);
const tags = ref([]);

const fetchArticles = async () => {
  loading.value = true;
  try {
    const params = {
      limit: 100,
      search: route.query.search,
      category_id: route.query.category_id,
      tag_id: route.query.tag_id
    };
    const response = await api.get('/articles', { params }); 
    articles.value = response.data;
  } catch (err) {
    console.error('Failed to fetch articles for timeline');
  } finally {
    loading.value = false;
  }
};

// Fetch categories and tags for displaying active filter names
const fetchMetadata = async () => {
  try {
    const [catRes, tagRes] = await Promise.all([
      api.get('/categories'),
      api.get('/tags')
    ]);
    categories.value = catRes.data;
    tags.value = tagRes.data;
  } catch (err) {
    console.error('Failed to fetch metadata');
  }
};

const activeCategory = computed(() => {
  if (!route.query.category_id) return null;
  return categories.value.find(c => c.id == route.query.category_id);
});

const activeTag = computed(() => {
  if (!route.query.tag_id) return null;
  return tags.value.find(t => t.id == route.query.tag_id);
});

const hasActiveFilters = computed(() => {
  return route.query.search || route.query.category_id || route.query.tag_id;
});

const handleFilter = (filters) => {
  // Merge new filters with existing query params
  const newQuery = { ...route.query, ...filters };
  
  // If setting category, maybe we don't clear tag? User asked for simultaneous.
  // But if user clicks a different category, it replaces the old one.
  
  router.push({ query: newQuery });
};

const clearFilter = (key) => {
  const newQuery = { ...route.query };
  delete newQuery[key];
  router.push({ query: newQuery });
};

const clearAllFilters = () => {
  router.push({ query: {} });
};

// Watch for any query changes
watch(() => route.query, () => {
  fetchArticles();
}, { deep: true });

const groupedArticles = computed(() => {
  const groups = {};
  const sorted = [...articles.value].sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
  
  sorted.forEach(article => {
    const date = new Date(article.created_at);
    const year = date.getFullYear();
    if (!groups[year]) {
      groups[year] = [];
    }
    groups[year].push(article);
  });
  return groups;
});

const formatDate = (dateString) => {
  const date = new Date(dateString);
  return `${date.getMonth() + 1}/${date.getDate()}`;
};

onMounted(() => {
  fetchMetadata();
  fetchArticles();
});
</script>

<style scoped>
.timeline-layout {
  display: grid;
  grid-template-columns: 1fr 300px;
  gap: var(--spacing-xl);
}

@media (max-width: 768px) {
  .timeline-layout {
    grid-template-columns: 1fr;
  }
}

.page-title {
  color: var(--primary-color);
  text-shadow: var(--neon-text-shadow);
  margin-bottom: var(--spacing-xl);
}

.active-filters {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
  justify-content: center;
  margin-bottom: var(--spacing-lg);
}

.filter-badge {
  background: rgba(0, 243, 255, 0.1);
  border: 1px solid var(--primary-color);
  color: var(--primary-color);
  padding: 0.25rem 0.75rem;
  border-radius: 20px;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.9rem;
}

.filter-badge button {
  background: none;
  border: none;
  color: var(--primary-color);
  cursor: pointer;
  font-size: 1.1rem;
  padding: 0;
  line-height: 1;
}

.filter-badge button:hover {
  color: #fff;
}

.clear-all-btn {
  background: none;
  border: 1px solid var(--secondary-color);
  color: var(--secondary-color);
  padding: 0.25rem 0.75rem;
  cursor: pointer;
  font-size: 0.9rem;
  transition: all 0.3s;
}

.clear-all-btn:hover {
  background: var(--secondary-color);
  color: #000;
}

.timeline {
  position: relative;
  border-left: 2px solid var(--surface-border);
  margin-left: 20px;
  padding-left: 30px;
}

.no-articles {
  text-align: center;
  color: var(--text-secondary);
  padding: 2rem;
  font-style: italic;
}

.timeline-year {
  margin-bottom: var(--spacing-xl);
}

.year-title {
  position: relative;
  left: -45px;
  color: var(--primary-color);
  font-size: 2rem;
  font-family: var(--font-heading);
  text-shadow: var(--neon-text-shadow);
  margin-bottom: var(--spacing-lg);
  background: var(--background-color);
  display: inline-block;
  padding: 0 10px;
}

.timeline-item {
  position: relative;
  margin-bottom: var(--spacing-lg);
}

.timeline-dot {
  position: absolute;
  left: -36px;
  top: 5px;
  width: 12px;
  height: 12px;
  background: var(--primary-color);
  border-radius: 50%;
  box-shadow: 0 0 10px var(--primary-color);
}

.timeline-date {
  position: absolute;
  left: -80px;
  top: 2px;
  color: var(--text-secondary);
  font-family: var(--font-heading);
  font-size: 0.9rem;
  width: 40px;
  text-align: right;
}

.timeline-content {
  padding: var(--spacing-md);
  margin-left: 0;
  border: 1px solid var(--surface-border);
  transition: all 0.3s ease;
}

.timeline-content:hover {
  border-color: var(--primary-color);
  transform: translateX(5px);
  box-shadow: 0 0 15px rgba(0, 243, 255, 0.1);
}

.timeline-link {
  text-decoration: none;
}

.timeline-link h3 {
  margin: 0;
  color: var(--text-primary);
  font-size: 1.2rem;
  transition: color 0.3s;
}

.timeline-link:hover h3 {
  color: var(--primary-color);
}

.timeline-meta {
  margin-top: var(--spacing-sm);
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.category-tag {
  font-size: 0.8rem;
  color: var(--secondary-color);
  border: 1px solid var(--secondary-color);
  padding: 2px 6px;
  border-radius: 2px;
  text-transform: uppercase;
}

.tag-pill {
  font-size: 0.8rem;
  color: var(--text-secondary);
  background: rgba(255, 255, 255, 0.05);
  padding: 2px 6px;
  border-radius: 2px;
}
</style>
