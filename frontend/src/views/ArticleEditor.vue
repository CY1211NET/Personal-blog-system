<template>
  <div class="editor-container">
    <h2>{{ isEditing ? 'Edit Article' : 'New Article' }}</h2>
    <form @submit.prevent="saveArticle">
      <div class="form-group">
        <label for="title">Title</label>
        <input type="text" id="title" v-model="title" required />
      </div>
      <div class="form-group">
        <label for="content">Content</label>
        <textarea id="content" v-model="content" required></textarea>
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

const route = useRoute();
const router = useRouter();

const title = ref('');
const content = ref('');
const loading = ref(false);

const isEditing = computed(() => route.params.id !== undefined);

const fetchArticle = async () => {
  if (!isEditing.value) return;
  try {
    const response = await api.get(`/articles/${route.params.id}`);
    title.value = response.data.title;
    content.value = response.data.content;
  } catch (err) {
    alert('Failed to load article');
    router.push('/');
  }
};

const saveArticle = async () => {
  loading.value = true;
  try {
    if (isEditing.value) {
      await api.put(`/articles/${route.params.id}`, {
        title: title.value,
        content: content.value,
      });
    } else {
      await api.post('/articles/', {
        title: title.value,
        content: content.value,
      });
    }
    router.push('/');
  } catch (err) {
    alert('Failed to save article');
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
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
input, textarea {
  width: 100%;
  padding: 10px;
  box-sizing: border-box;
}
textarea {
  height: 300px;
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
