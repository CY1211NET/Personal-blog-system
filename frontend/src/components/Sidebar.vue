<template>
  <div class="sidebar card">
    <div class="sidebar-section">
      <h3>{{ $t('sidebar.categories') }}</h3>
      <ul class="category-list">
        <li v-for="category in categories" :key="category.id">
          <a href="#" @click.prevent="$emit('filter', { category_id: category.id })">
            {{ category.name }}
          </a>
        </li>
      </ul>
    </div>

    <div class="sidebar-section">
      <h3>{{ $t('sidebar.tags') }}</h3>
      <div class="tag-cloud">
        <a 
          href="#" 
          v-for="tag in tags" 
          :key="tag.id" 
          class="tag-pill"
          @click.prevent="$emit('filter', { search: '#' + tag.name })"
        >
          {{ tag.name }}
        </a>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import api from '../api/axios';

const categories = ref([]);
const tags = ref([]);

onMounted(async () => {
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

.nav-list {
  list-style: none;
  padding: 0;
}

.nav-link {
  color: var(--text-secondary);
  display: block;
  padding: var(--spacing-xs) 0;
  transition: all 0.3s ease;
  font-family: var(--font-heading);
  letter-spacing: 1px;
  text-transform: uppercase;
  font-size: 1rem;
  text-decoration: none;
}

.nav-link:hover {
  color: var(--primary-color);
  transform: translateX(10px);
  text-shadow: var(--neon-text-shadow);
}

.category-list {
  list-style: none;
  padding: 0;
}

.category-list li {
  margin-bottom: var(--spacing-sm);
}

.category-list a {
  color: var(--text-secondary);
  display: block;
  padding: var(--spacing-xs) 0;
  transition: all 0.3s ease;
  font-family: var(--font-heading);
  letter-spacing: 1px;
  text-transform: uppercase;
  font-size: 0.9rem;
}

.category-list a:hover {
  color: var(--primary-color);
  transform: translateX(10px);
  text-shadow: var(--neon-text-shadow);
}

.tag-cloud {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-sm);
}

.tag-pill {
  background-color: rgba(188, 19, 254, 0.1);
  color: var(--secondary-color);
  padding: 0.25rem 0.75rem;
  border: 1px solid var(--secondary-color);
  border-radius: 0;
  font-family: var(--font-heading);
  font-size: 0.8rem;
  transition: all 0.3s ease;
  text-transform: uppercase;
}

.tag-pill:hover {
  background-color: var(--secondary-color);
  color: #000;
  box-shadow: 0 0 10px var(--secondary-glow);
  transform: scale(1.05);
}
</style>
