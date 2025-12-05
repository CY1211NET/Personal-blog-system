<template>
  <div class="sidebar card">
    <div class="sidebar-section">
      <h3>{{ $t('sidebar.categories') }}</h3>
      <ul class="category-list">
        <li v-for="category in categories" :key="category.id" class="category-item">
          <a href="#" @click.prevent="$emit('filter', { category_id: category.id })">
            {{ category.name }}
          </a>
          <div v-if="authStore.isAuthenticated" class="item-actions">
            <button @click.prevent="editCategory(category)" class="btn-icon" title="Edit">✏️</button>
            <button @click.prevent="deleteCategory(category.id)" class="btn-icon" title="Delete">🗑️</button>
          </div>
        </li>
      </ul>
    </div>

    <div class="sidebar-section">
      <h3>{{ $t('sidebar.tags') }}</h3>
      <div class="tag-cloud">
        <div v-for="tag in tags" :key="tag.id" class="tag-wrapper">
          <a 
            href="#" 
            class="tag-pill"
            @click.prevent="$emit('filter', { tag_id: tag.id })"
          >
            {{ tag.name }}
          </a>
          <div v-if="authStore.isAuthenticated" class="tag-actions">
             <span @click.prevent="editTag(tag)" class="action-icon" title="Edit">✏️</span>
             <span @click.prevent="deleteTag(tag.id)" class="action-icon" title="Delete">×</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import api from '../api/axios';
import { useAuthStore } from '../stores/auth';
import { useToast } from '../composables/useToast';

const categories = ref([]);
const tags = ref([]);
const authStore = useAuthStore();
const toast = useToast();

const fetchSidebarData = async () => {
  try {
    const [catRes, tagRes] = await Promise.all([
      api.get('/categories'),
      api.get('/tags')
    ]);
    categories.value = catRes.data;
    tags.value = tagRes.data;
  } catch (err) {
    console.error('Failed to load sidebar data');
  }
};

const editCategory = async (category) => {
  const newName = prompt('Edit category name:', category.name);
  if (newName && newName !== category.name) {
    try {
      await api.put(`/categories/${category.id}`, { name: newName });
      toast.success('Category updated');
      fetchSidebarData();
    } catch (err) {
      toast.error('Failed to update category');
    }
  }
};

const deleteCategory = async (id) => {
  if (confirm('Delete this category? ALL associated articles will also be DELETED!')) {
    try {
      await api.delete(`/categories/${id}`);
      toast.success('Category and articles deleted');
      fetchSidebarData();
      // Emit event to refresh article list if needed, or just let user navigate
      window.location.reload(); // Simple reload to refresh main list
    } catch (err) {
      toast.error('Failed to delete category');
    }
  }
};

const editTag = async (tag) => {
  const newName = prompt('Edit tag name:', tag.name);
  if (newName && newName !== tag.name) {
    try {
      await api.put(`/tags/${tag.id}`, { name: newName });
      toast.success('Tag updated');
      fetchSidebarData();
    } catch (err) {
      toast.error('Failed to update tag');
    }
  }
};

const deleteTag = async (id) => {
  if (confirm('Delete this tag? ALL associated articles will also be DELETED!')) {
    try {
      await api.delete(`/tags/${id}`);
      toast.success('Tag and articles deleted');
      fetchSidebarData();
      window.location.reload();
    } catch (err) {
      toast.error('Failed to delete tag');
    }
  }
};

onMounted(() => {
  fetchSidebarData();
});
</script>

<style scoped>
.sidebar {
  height: fit-content;
  position: sticky;
  top: 100px;
  border: 1px solid var(--surface-border);
}

.sidebar-section {
  margin-bottom: var(--spacing-xl);
}

.sidebar-section:last-child {
  margin-bottom: 0;
}

.sidebar-section h3 {
  font-size: 1.2rem;
  margin-bottom: var(--spacing-md);
  padding-bottom: var(--spacing-xs);
  border-bottom: 1px solid var(--primary-color);
  color: var(--primary-color);
  text-shadow: var(--neon-text-shadow);
}

.category-list {
  list-style: none;
  padding: 0;
}

.category-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-sm);
}

.category-item a {
  color: var(--text-secondary);
  text-decoration: none;
  transition: all 0.3s ease;
  font-family: var(--font-heading);
  letter-spacing: 1px;
  text-transform: uppercase;
  font-size: 0.9rem;
}

.category-item a:hover {
  color: var(--primary-color);
  transform: translateX(5px);
  text-shadow: var(--neon-text-shadow);
}

.item-actions {
  display: flex;
  gap: 5px;
  opacity: 0; /* Hidden by default */
  transition: opacity 0.3s;
}

.category-item:hover .item-actions {
  opacity: 1;
}

.btn-icon {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 0.8rem;
  padding: 0;
  opacity: 0.7;
}

.btn-icon:hover {
  opacity: 1;
  transform: scale(1.2);
}

.tag-cloud {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-sm);
}

.tag-wrapper {
  position: relative;
  display: inline-block;
}

.tag-pill {
  background-color: rgba(188, 19, 254, 0.1);
  color: var(--secondary-color);
  padding: 0.25rem 0.75rem;
  border: 1px solid var(--secondary-color);
  font-family: var(--font-heading);
  font-size: 0.8rem;
  transition: all 0.3s ease;
  text-transform: uppercase;
  text-decoration: none;
  display: block;
}

.tag-pill:hover {
  background-color: var(--secondary-color);
  color: #000;
  box-shadow: 0 0 10px var(--secondary-glow);
}

.tag-actions {
  position: absolute;
  top: -8px;
  right: -8px;
  display: flex;
  gap: 2px;
  background: rgba(0,0,0,0.8);
  border-radius: 4px;
  padding: 2px;
  opacity: 0;
  transition: opacity 0.3s;
  z-index: 10;
}

.tag-wrapper:hover .tag-actions {
  opacity: 1;
}

.action-icon {
  cursor: pointer;
  font-size: 0.7rem;
  color: #fff;
  padding: 0 2px;
}

.action-icon:hover {
  color: var(--primary-color);
}
</style>
