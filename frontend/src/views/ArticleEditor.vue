<template>
  <div class="editor-container">
    <h2>{{ isEditing ? 'Edit Article' : 'New Article' }}</h2>
    <form @submit.prevent="saveArticle">
      <div class="form-group">
        <label for="title">Title</label>
        <input type="text" id="title" v-model="title" required />
      </div>
      
      <div class="form-group">
        <label for="category">Category</label>
        <select id="category" v-model="selectedCategory">
          <option value="">Select Category</option>
          <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
        </select>
        <button type="button" @click="createCategory" style="margin-top: 5px; font-size: 0.8em;">+ New Category</button>
      </div>

      <div class="form-group">
        <label for="tags">Tags (comma separated)</label>
        <input type="text" id="tags" v-model="tagsInput" placeholder="tech, coding, life" />
      </div>

      <div class="form-group">
        <label for="content">Content (Markdown supported)</label>
        <div class="editor-layout">
          <textarea id="content" v-model="content" required class="editor-textarea"></textarea>
          <div class="preview" v-html="renderedContent"></div>
        </div>
      </div>
      <button type="submit" :disabled="loading">
        {{ loading ? 'Saving...' : 'Save Article' }}
      </button>
    </form>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import api from '../api/axios';
import { marked } from 'marked';
import DOMPurify from 'dompurify';

const route = useRoute();
const router = useRouter();

const title = ref('');
const content = ref('');
const selectedCategory = ref('');
const tagsInput = ref('');
const categories = ref([]);
const loading = ref(false);

const isEditing = computed(() => route.params.id !== undefined);

const renderedContent = computed(() => {
  return DOMPurify.sanitize(marked(content.value || ''));
});

const fetchCategories = async () => {
  try {
    const response = await api.get('/categories');
    categories.value = response.data;
  } catch (err) {
    console.error('Failed to load categories');
  }
};

const createCategory = async () => {
  const name = prompt('Enter category name:');
  if (!name) return;
  try {
    const response = await api.post('/categories', { name });
    categories.value.push(response.data);
    selectedCategory.value = response.data.id;
  } catch (err) {
    alert('Failed to create category');
  }
};

const fetchArticle = async () => {
  if (!isEditing.value) return;
  try {
    const response = await api.get(`/articles/${route.params.id}`);
    const article = response.data;
    title.value = article.title;
    content.value = article.content;
    selectedCategory.value = article.category_id || '';
    tagsInput.value = article.tags ? article.tags.map(t => t.name).join(', ') : '';
  } catch (err) {
    alert('Failed to load article');
    router.push('/');
  }
};

const saveArticle = async () => {
  loading.value = true;
  const tags = tagsInput.value.split(',').map(t => t.trim()).filter(t => t);
  const payload = {
    title: title.value,
    content: content.value,
    category_id: selectedCategory.value ? parseInt(selectedCategory.value) : null,
    tags: tags
  };

  try {
    if (isEditing.value) {
      await api.put(`/articles/${route.params.id}`, payload);
    } else {
      await api.post('/articles/', payload);
    }
    router.push('/');
  } catch (err) {
    alert('Failed to save article');
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  fetchCategories();
  fetchArticle();
});
</script>

<style scoped>
.editor-container {
  max-width: 800px;
  margin: 50px auto;
  padding: 20px;
}
.form-group {
  margin-bottom: 20px;
}
label {
  display: block;
  margin-bottom: 5px;
}
input, textarea, select {
  width: 100%;
  padding: 10px;
  box-sizing: border-box;
}
.editor-layout {
  display: flex;
  gap: 20px;
}
.editor-textarea {
  height: 400px;
  width: 50%;
  font-family: monospace;
}
.preview {
  width: 50%;
  height: 400px;
  overflow-y: auto;
  border: 1px solid #ccc;
  padding: 10px;
  border-radius: 4px;
  background-color: #f9f9f9;
}
button {
  padding: 10px 20px;
  background-color: #42b983;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}
</style>
