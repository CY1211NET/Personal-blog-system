<template>
  <div class="article-list-container">
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
</template>

<script setup>
import { ref, onMounted } from 'vue';
import api from '../api/axios';
import { useAuthStore } from '../stores/auth';

const articles = ref([]);
const loading = ref(true);
const error = ref('');
const searchQuery = ref('');
const authStore = useAuthStore();
let timeout = null;

const fetchArticles = async () => {
  loading.value = true;
  try {
    const response = await api.get('/articles', {
      params: { search: searchQuery.value }
    });
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
    fetchArticles();
  }, 300);
};

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString();
};

onMounted(() => {
  fetchArticles();
});
</script>

<style scoped>
.article-list-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
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
</style>
