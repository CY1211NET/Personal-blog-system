<template>
  <div class="home-layout">
    <div class="content-area">
      <!-- Toolbar for Batch Actions -->
      <div v-if="authStore.isAuthenticated" class="toolbar fade-in">
         <button v-if="selectedArticles.length > 0" @click="batchDelete" class="btn btn-danger">
           {{ $t('common.delete_selected') }} ({{ selectedArticles.length }})
         </button>
      </div>

      <div v-if="loading" class="loading">{{ $t('common.loading') }}</div>
      <div v-else class="article-list">
        <div v-for="article in articles" :key="article.id" class="card article-card fade-in">
          <div class="card-header-actions">
            <!-- Checkbox for batch delete -->
            <div v-if="authStore.isAuthenticated" class="checkbox-wrapper">
               <input type="checkbox" :value="article.id" v-model="selectedArticles" class="custom-checkbox" />
            </div>
            
            <router-link :to="'/articles/' + article.id" class="article-title">
              <h2>{{ article.title }}</h2>
            </router-link>
          </div>

          <div class="article-meta text-sm text-gray">
            <span>{{ $t('article.by') }} {{ article.author?.username }}</span>
            <span> • {{ formatDate(article.created_at) }}</span>
            <span v-if="article.category"> • {{ $t('article.in') }} {{ article.category.name }}</span>
            <span> • {{ article.views }} {{ $t('article.views') }}</span>
          </div>
          <p class="article-excerpt" v-html="getExcerpt(article.content)"></p>
          <div class="article-tags" v-if="article.tags && article.tags.length">
            <span v-for="tag in article.tags" :key="tag.id" class="tag">#{{ tag.name }}</span>
          </div>
          
          <div class="card-footer">
            <div class="read-more">
              <router-link :to="'/articles/' + article.id" class="btn btn-ghost">{{ $t('article.read_more') }} →</router-link>
            </div>
            <!-- Edit/Delete Actions -->
            <div v-if="authStore.isAuthenticated" class="article-actions">
               <button @click.prevent="editArticle(article.id)" class="btn btn-sm btn-primary">{{ $t('common.edit') }}</button>
               <button @click.prevent="deleteArticle(article.id)" class="btn btn-sm btn-danger">{{ $t('common.delete') }}</button>
            </div>
          </div>
        </div>
      </div>
    </div>
    <Sidebar @filter="handleFilter" class="sidebar-area fade-in" />
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import api from '../api/axios';
import Sidebar from '../components/Sidebar.vue';
import DOMPurify from 'dompurify';
import { useAuthStore } from '../stores/auth';
import { useToast } from '../composables/useToast';

const route = useRoute();
const router = useRouter();
const authStore = useAuthStore();
const toast = useToast();

const articles = ref([]);
const loading = ref(true);
const selectedArticles = ref([]);

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
  
  // Strip markdown/HTML roughly to get plain text
  let text = content.replace(/[#*`]/g, '')
                    .replace(/!\[.*?\]\(.*?\)/g, '') // Remove images
                    .replace(/\[.*?\]\(.*?\)/g, '$1') // Keep link text
                    .replace(/<[^>]*>/g, ''); // Remove HTML tags

  const query = route.query.search;

  if (query && query.trim()) {
    const lowerText = text.toLowerCase();
    const lowerQuery = query.toLowerCase();
    const index = lowerText.indexOf(lowerQuery);

    if (index !== -1) {
      // Found keyword
      const start = Math.max(0, index - 50);
      const end = Math.min(text.length, index + query.length + 100);
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

  // Default behavior
  const snippet = text.length > 150 ? text.substring(0, 150) + '...' : text;
  return DOMPurify.sanitize(snippet, { ALLOWED_TAGS: [] });
};

const editArticle = (id) => {
  router.push(`/articles/${id}/edit`);
};

const deleteArticle = async (id) => {
  if (!confirm('Are you sure you want to delete this article?')) return;
  try {
    await api.delete(`/articles/${id}`);
    toast.success('Article deleted successfully');
    fetchArticles();
  } catch (err) {
    toast.error('Failed to delete article');
  }
};

const batchDelete = async () => {
  if (!confirm(`Are you sure you want to delete ${selectedArticles.value.length} articles?`)) return;
  try {
    await api.post('/articles/batch-delete', { ids: selectedArticles.value });
    toast.success('Articles deleted successfully');
    selectedArticles.value = [];
    fetchArticles();
  } catch (err) {
    toast.error('Failed to delete articles');
  }
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
  position: relative;
}

.article-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.5);
  border-color: var(--primary-color);
}

.card-header-actions {
  display: flex;
  align-items: flex-start;
  gap: 1rem;
}

.checkbox-wrapper {
  padding-top: 0.5rem;
}

.custom-checkbox {
  width: 18px;
  height: 18px;
  cursor: pointer;
  accent-color: var(--primary-color);
}

.article-title {
  text-decoration: none;
  color: var(--text-primary);
  flex-grow: 1;
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

/* Highlight style */
:deep(mark.highlight) {
  background-color: rgba(0, 243, 255, 0.2);
  color: var(--primary-color);
  padding: 0 2px;
  border-radius: 2px;
  font-weight: bold;
  box-shadow: 0 0 5px rgba(0, 243, 255, 0.3);
}

.article-tags {
  margin-bottom: var(--spacing-md);
}

.tag {
  margin-right: 0.5rem;
  font-size: 0.8rem;
  color: var(--secondary-color);
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: var(--spacing-md);
}

.article-actions {
  display: flex;
  gap: 0.5rem;
}

.btn-sm {
  padding: 0.25rem 0.5rem;
  font-size: 0.8rem;
}

.toolbar {
  margin-bottom: var(--spacing-md);
  display: flex;
  justify-content: flex-end;
}
</style>
