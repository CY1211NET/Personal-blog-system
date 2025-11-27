<template>
  <div class="home-layout">
    <div class="content-area">
      <div v-if="loading" class="loading">{{ $t('common.loading') }}</div>
      <div v-else class="article-list">
        <div v-for="article in articles" :key="article.id" class="card article-card fade-in">
          <router-link :to="'/articles/' + article.id" class="article-title">
            <h2>{{ article.title }}</h2>
          </router-link>
          <div class="article-meta text-sm text-gray">
            <span>{{ $t('article.by') }} {{ article.author?.username }}</span>
            <span> • {{ formatDate(article.created_at) }}</span>
            <span v-if="article.category"> • {{ $t('article.in') }} {{ article.category.name }}</span>
            <span> • {{ article.views }} {{ $t('article.views') }}</span>
          </div>
          <p class="article-excerpt">{{ getExcerpt(article.content) }}</p>
          <div class="article-tags" v-if="article.tags && article.tags.length">
            <span v-for="tag in article.tags" :key="tag.id" class="tag">#{{ tag.name }}</span>
          </div>
          <div class="read-more">
            <router-link :to="'/articles/' + article.id" class="btn btn-ghost">{{ $t('article.read_more') }} →</router-link>
          </div>
        </div>
      </div>
    </div>
    <Sidebar @filter="handleFilter" class="sidebar-area fade-in" />
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue';
import { useRoute } from 'vue-router';
import api from '../api/axios';
import Sidebar from '../components/Sidebar.vue';
import { marked } from 'marked';

const route = useRoute();
const articles = ref([]);
const loading = ref(true);

const fetchArticles = async (params = {}) => {
  loading.value = true;
  try {
    const response = await api.get('/articles', { params });
    articles.value = response.data;
  } catch (err) {
    console.error('Failed to fetch articles');
  } finally {
    loading.value = false;
  }
};

const handleFilter = (filters) => {
  fetchArticles(filters);
};

// Watch for search query changes from App.vue
watch(() => route.query.search, (newSearch) => {
  fetchArticles({ search: newSearch });
});

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString();
};

const getExcerpt = (content) => {
  if (!content) return '';
  const text = content.replace(/[#*`]/g, ''); // Simple markdown strip
  return text.length > 150 ? text.substring(0, 150) + '...' : text;
};

onMounted(() => {
  // Initial fetch with query params if any
  fetchArticles(route.query.search ? { search: route.query.search } : {});
});
</script>

<style scoped>
.home-layout {
  display: grid;
  grid-template-columns: 1fr 300px;
  gap: var(--spacing-xl);
}

@media (max-width: 768px) {
  .home-layout {
    grid-template-columns: 1fr;
  }
}

.article-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.article-card {
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.article-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.5);
  border-color: var(--primary-color);
}

.article-title {
  text-decoration: none;
  color: var(--text-primary);
}

.article-title h2 {
  margin-bottom: var(--spacing-sm);
  transition: color 0.3s ease;
}

.article-title:hover h2 {
  color: var(--primary-color);
  text-shadow: var(--neon-text-shadow);
}

.article-meta {
  margin-bottom: var(--spacing-md);
}

.article-excerpt {
  margin-bottom: var(--spacing-md);
  color: var(--text-secondary);
  line-height: 1.6;
}

.article-tags {
  margin-bottom: var(--spacing-md);
}

.tag {
  margin-right: 0.5rem;
  font-size: 0.8rem;
  color: var(--secondary-color);
}

.read-more {
  text-align: right;
}
</style>
