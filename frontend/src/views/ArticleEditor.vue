<template>
  <div class="editor-container fade-in">
    <div class="card">
      <h2 class="text-center">{{ isEditing ? $t('common.edit') : $t('article.write_new') }}</h2>
      <form @submit.prevent="saveArticle">
        <div class="form-group">
          <label for="title">{{ $t('article.title') }}</label>
          <input type="text" id="title" v-model="title" required placeholder="ENTER TITLE" />
        </div>
        
        <div class="form-group">
          <label for="category">{{ $t('article.category') }}</label>
          <div class="flex gap-2">
            <select id="category" v-model="selectedCategory">
              <option value="">{{ $t('article.select_category') }}</option>
              <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
            </select>
            <button type="button" @click="createCategory" class="btn btn-ghost">{{ $t('article.new_category') }}</button>
          </div>
        </div>

        <div class="form-group">
          <label for="tags">{{ $t('article.tags') }}</label>
          <input type="text" id="tags" v-model="tagsInput" :placeholder="$t('article.tags_placeholder')" />
        </div>

        <div class="form-group">
          <label>{{ $t('article.content') }}</label>
          <div id="vditor" class="vditor-container"></div>
        </div>
        
        <div class="form-actions">
          <button type="button" @click="$router.back()" class="btn btn-ghost">{{ $t('common.cancel') }}</button>
          <button type="submit" :disabled="loading" class="btn btn-primary">
            {{ loading ? $t('common.loading') : $t('common.save') }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import api from '../api/axios';
import { useToast } from '../composables/useToast';
import Vditor from 'vditor';
import 'vditor/dist/index.css';

const route = useRoute();
const router = useRouter();
const toast = useToast();

const title = ref('');
const content = ref('');
const selectedCategory = ref('');
const tagsInput = ref('');
const categories = ref([]);
const loading = ref(false);
const vditor = ref(null);

const isEditing = computed(() => route.params.id !== undefined);

const fetchCategories = async () => {
  try {
    const response = await api.get('/categories');
    categories.value = response.data;
  } catch (err) {
    toast.error('Failed to load categories');
  }
};

const createCategory = async () => {
  const name = prompt('Enter category name:');
  if (!name) return;
  try {
    const response = await api.post('/categories', { name });
    categories.value.push(response.data);
    selectedCategory.value = response.data.id;
    toast.success('Category created');
  } catch (err) {
    toast.error('Failed to create category');
  }
};

const vditorReady = ref(false);

const fetchArticle = async () => {
  if (!isEditing.value) return;
  try {
    const response = await api.get(`/articles/${route.params.id}`);
    const article = response.data;
    title.value = article.title;
    content.value = article.content;
    selectedCategory.value = article.category_id || '';
    tagsInput.value = article.tags ? article.tags.map(t => t.name).join(', ') : '';
    
    if (vditorReady.value && vditor.value) {
      vditor.value.setValue(content.value);
    }
  } catch (err) {
    toast.error('Failed to load article');
    router.push('/');
  }
};

// ... saveArticle ...
const saveArticle = async () => {
  loading.value = true;
  const currentContent = vditor.value ? vditor.value.getValue() : content.value;
  const tags = tagsInput.value.split(',').map(t => t.trim()).filter(t => t);
  const payload = {
    title: title.value,
    content: currentContent,
    category_id: selectedCategory.value ? parseInt(selectedCategory.value) : null,
    tags: tags
  };

  try {
    if (isEditing.value) {
      await api.put(`/articles/${route.params.id}`, payload);
      toast.success('Article updated');
    } else {
      await api.post('/articles/', payload);
      toast.success('Article created');
    }
    router.push('/');
  } catch (err) {
    toast.error('Failed to save article');
  } finally {
    loading.value = false;
  }
};

const initVditor = () => {
  vditor.value = new Vditor('vditor', {
    height: 600,
    mode: 'ir', // Instant Rendering mode (like Obsidian)
    theme: 'dark',
    placeholder: 'INITIATE DATA STREAM...',
    toolbar: [
      'emoji', 'headings', 'bold', 'italic', 'strike', 'link', '|',
      'list', 'ordered-list', 'check', 'outdent', 'indent', '|',
      'quote', 'line', 'code', 'inline-code', 'insert-before', 'insert-after', '|',
      'upload', 'table', '|',
      'undo', 'redo', '|',
      'fullscreen', 'edit-mode',
      {
        name: 'more',
        toolbar: [
          'both',
          'code-theme',
          'content-theme',
          'export',
          'outline',
          'preview',
        ],
      },
    ],
    cache: {
      enable: false,
    },
    after: () => {
      vditorReady.value = true;
      if (content.value) {
        vditor.value.setValue(content.value);
      }
    },
    input: (value) => {
      content.value = value;
    },
    preview: {
      theme: {
        current: 'dark',
        path: 'https://cdn.jsdelivr.net/npm/vditor/dist/css/content-theme',
      },
      hljs: {
        style: 'atom-one-dark',
      }
    }
  });
};

onMounted(() => {
  fetchCategories();
  initVditor();
  fetchArticle();
});
</script>

<style scoped>
.editor-container {
  max-width: 1200px;
  margin: 0 auto;
}

.text-center {
  text-align: center;
  margin-bottom: var(--spacing-xl);
  color: var(--primary-color);
  text-shadow: var(--neon-text-shadow);
}

.form-group {
  margin-bottom: var(--spacing-lg);
}

label {
  display: block;
  margin-bottom: var(--spacing-xs);
  color: var(--primary-color);
  font-family: var(--font-heading);
  font-size: 0.9rem;
  letter-spacing: 1px;
}

/* Vditor Customization for Sci-Fi Theme */
:deep(.vditor) {
  border: 1px solid var(--primary-color) !important;
  background-color: rgba(5, 5, 16, 0.9) !important;
  --vditor-border-color: var(--surface-border);
  --vditor-tool-icon-hover-color: var(--primary-color);
}

:deep(.vditor-toolbar) {
  background-color: rgba(0, 0, 0, 0.5) !important;
  border-bottom: 1px solid var(--surface-border) !important;
}

:deep(.vditor-content) {
  background-color: transparent !important;
}

:deep(.vditor-reset) {
  color: var(--text-primary) !important;
  font-family: var(--font-body) !important;
}

:deep(.vditor-textarea) {
  background-color: transparent !important;
  color: var(--text-primary) !important;
}
</style>
