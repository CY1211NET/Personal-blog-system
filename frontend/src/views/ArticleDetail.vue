<template>
  <div class="article-detail-container" v-if="article">
    <h1>{{ article.title }}</h1>
    <p class="meta">
      By {{ article.author?.username }} on {{ formatDate(article.created_at) }}
      <span v-if="article.category"> | Category: <strong>{{ article.category.name }}</strong></span>
      <span> | Views: {{ article.views }}</span>
    </p>
    <div class="tags" v-if="article.tags && article.tags.length">
      <span v-for="tag in article.tags" :key="tag.id" class="tag">#{{ tag.name }}</span>
    </div>
    <div class="content" v-html="renderedContent"></div>
    
    <div class="interactions">
      <button @click="likeArticle" :disabled="liked" class="like-btn">
        {{ liked ? 'Liked' : 'Like' }} ({{ article.likes }})
      </button>
    </div>

    <div class="actions" v-if="isAuthor">
      <router-link :to="'/articles/' + article.id + '/edit'" class="edit-btn">Edit</router-link>
      <button @click="deleteArticle" class="delete-btn">Delete</button>
    </div>

    <div class="comments-section">
      <h3>Comments</h3>
      <div v-for="comment in comments" :key="comment.id" class="comment">
        <p><strong>{{ comment.user?.username }}</strong>: {{ comment.content }}</p>
      </div>
      
      <form @submit.prevent="addComment" v-if="authStore.isAuthenticated">
        <textarea v-model="newComment" placeholder="Add a comment..." required></textarea>
        <button type="submit">Post Comment</button>
      </form>
    </div>
  </div>
  <div v-else-if="loading">Loading...</div>
  <div v-else class="error">{{ error }}</div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import api from '../api/axios';
import { useAuthStore } from '../stores/auth';
import { marked } from 'marked';
import DOMPurify from 'dompurify';

const route = useRoute();
const router = useRouter();
const authStore = useAuthStore();
  }
};

const addComment = async () => {
  try {
    await api.post(`/articles/${route.params.id}/comments`, { content: newComment.value });
    newComment.value = '';
    fetchComments();
  } catch (err) {
    alert('Failed to post comment');
  }
};

const deleteArticle = async () => {
  if (!confirm('Are you sure?')) return;
  try {
    await api.delete(`/articles/${route.params.id}`);
    router.push('/');
  } catch (err) {
    alert('Failed to delete article');
  }
};

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString();
};

onMounted(() => {
  fetchArticle();
});
</script>

<style scoped>
.article-detail-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}
.meta {
  color: #666;
  margin-bottom: 20px;
}
.content {
  line-height: 1.6;
  margin-bottom: 40px;
}
.actions {
  margin-bottom: 20px;
}
.edit-btn, .delete-btn {
  margin-right: 10px;
  padding: 5px 10px;
  cursor: pointer;
}
.delete-btn {
  background-color: #ff4444;
  color: white;
  border: none;
  border-radius: 4px;
}
.comments-section {
  border-top: 1px solid #eee;
  padding-top: 20px;
}
.comment {
  margin-bottom: 10px;
  padding: 10px;
  background-color: #f9f9f9;
  border-radius: 4px;
}
textarea {
  width: 100%;
  height: 60px;
  margin-bottom: 10px;
}
.tags {
  margin-bottom: 20px;
}
.tag {
  background-color: #eee;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 0.9em;
  margin-right: 8px;
  color: #555;
}
.interactions {
  margin: 20px 0;
  border-top: 1px solid #eee;
  border-bottom: 1px solid #eee;
  padding: 10px 0;
}
.like-btn {
  background-color: #ff4757;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 20px;
  cursor: pointer;
  font-size: 1em;
  transition: background-color 0.2s;
}
.like-btn:disabled {
  background-color: #ff6b81;
  cursor: default;
}
.like-btn:hover:not(:disabled) {
  background-color: #e84118;
}
</style>
