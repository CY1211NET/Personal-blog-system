<template>
  <div class="sidebar">
    <div class="widget">
      <h3>Categories</h3>
      <ul>
        <li v-for="cat in categories" :key="cat.id">
          <a href="#" @click.prevent="$emit('filter-category', cat.id)">{{ cat.name }}</a>
        </li>
      </ul>
    </div>
    <div class="widget">
      <h3>Tags</h3>
      <div class="tag-cloud">
        <span v-for="tag in tags" :key="tag.id" class="tag-item">
          <a href="#" @click.prevent="$emit('filter-tag', tag.name)">#{{ tag.name }}</a>
        </span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import api from '../api/axios';

const categories = ref([]);
const tags = ref([]);

const fetchData = async () => {
  try {
    const [catResp, tagResp] = await Promise.all([
      api.get('/categories'),
      api.get('/tags')
    ]);
    categories.value = catResp.data;
    tags.value = tagResp.data;
  } catch (err) {
    console.error('Failed to load sidebar data');
  }
};

onMounted(() => {
  fetchData();
});
</script>

<style scoped>
.sidebar {
  padding: 20px;
  background-color: #f9f9f9;
  border-radius: 8px;
}
.widget {
  margin-bottom: 30px;
}
h3 {
  margin-top: 0;
  border-bottom: 2px solid #ddd;
  padding-bottom: 10px;
  margin-bottom: 15px;
}
ul {
  list-style: none;
  padding: 0;
}
li {
  margin-bottom: 8px;
}
.tag-cloud {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}
.tag-item a {
  background-color: #e0e0e0;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 0.9em;
  color: #333;
  text-decoration: none;
}
.tag-item a:hover {
  background-color: #d0d0d0;
}
</style>
