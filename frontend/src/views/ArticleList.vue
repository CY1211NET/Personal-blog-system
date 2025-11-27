<template>
  <div class="article-list-layout">
    <div class="main-content">
      <h2>Articles</h2>
      <div class="search-bar">
        <input type="text" v-model="searchQuery" @input="handleSearch" placeholder="Search articles..." />
      </div>
      <div v-if="loading">Loading articles...</div>
      <div v-else-if="error" class="error">{{ error }}</div>
      <div v-else>
        <div v-for="article in articles" :key="article.id" class="article-card">
          <h3>
            <router-link :to="'/articles/' + article.id">{{ article.title }}</router-link>
          </h3>
          <p class="meta">
            By {{ article.author?.username }} on {{ formatDate(article.created_at) }}
            <span v-if="article.category"> | In <strong>{{ article.category.name }}</strong></span>
          </p>
          <div class="tags" v-if="article.tags && article.tags.length">
            <span v-for="tag in article.tags" :key="tag.id" class="tag">#{{ tag.name }}</span>
          </div>
          <p class="excerpt">{{ article.content.substring(0, 100) }}...</p>
        </div>
      </div>
      <router-link to="/articles/new" class="create-btn" v-if="authStore.isAuthenticated">Write Article</router-link>
    </div>
    <div class="sidebar-container">
      <Sidebar @filter-category="filterByCategory" @filter-tag="filterByTag" />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import api from '../api/axios';
import { useAuthStore } from '../stores/auth';
import Sidebar from '../components/Sidebar.vue';

const articles = ref([]);
const loading = ref(true);
const error = ref('');
const searchQuery = ref('');
const authStore = useAuthStore();
let timeout = null;

const fetchArticles = async (params = {}) => {
  loading.value = true;
  try {
    const response = await api.get('/articles', { params });
    articles.value = response.data;
  } catch (err) {
    error.value = 'Failed to load articles';
    console.error(err);
  } finally {
    loading.value = false;
  }
};

const handleSearch = () => {
  clearTimeout(timeout);
  timeout = setTimeout(() => {
    fetchArticles({ search: searchQuery.value });
  }, 300);
};

const filterByCategory = (categoryId) => {
  fetchArticles({ category_id: categoryId });
};

const filterByTag = (tagName) => {
  // Backend support for tag filtering needs to be implemented if not already
  // For now, let's just search by tag name as a workaround or implement tag filtering in backend
  fetchArticles({ search: tagName }); 
};

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString();
};

onMounted(() => {
  fetchArticles();
});
</script>

<style scoped>
.article-list-layout {
  display: flex;
  gap: 30px;
}
.main-content {
  flex: 3;
}
.sidebar-container {
  flex: 1;
}
.article-card {
  border-bottom: 1px solid #eee;
  padding: 20px 0;
}
.meta {
  color: #666;
  font-size: 0.9em;
}
.create-btn {
  display: inline-block;
  margin-top: 20px;
  padding: 10px 20px;
  background-color: #42b983;
  color: white;
  text-decoration: none;
  border-radius: 4px;
}
.tags {
  margin-bottom: 10px;
}
.tag {
  background-color: #eee;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 0.8em;
  margin-right: 5px;
  color: #666;
}
.error {
  color: red;
}

@media (max-width: 768px) {
  .article-list-layout {
    flex-direction: column;
  }
}
</style>
