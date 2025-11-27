<template>
  <div class="timeline-container">
    <div class="timeline">
      <div v-for="(group, year) in groupedArticles" :key="year" class="timeline-year">
        <h2 class="year-title">{{ year }}</h2>
        <div class="timeline-items">
          <div v-for="article in group" :key="article.id" class="timeline-item">
            <div class="timeline-dot"></div>
            <div class="timeline-date">{{ formatDate(article.created_at) }}</div>
            <div class="timeline-content card">
              <router-link :to="'/articles/' + article.id" class="timeline-link">
                <h3>{{ article.title }}</h3>
              </router-link>
              <div class="timeline-meta">
                <span v-if="article.category" class="category-tag">{{ article.category.name }}</span>
                <span v-for="tag in article.tags" :key="tag.id" class="tag-pill">#{{ tag.name }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  articles: {
    type: Array,
    required: true
  }
});

const groupedArticles = computed(() => {
  const groups = {};
  // Sort articles by date descending
  const sorted = [...props.articles].sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
  
  sorted.forEach(article => {
    const date = new Date(article.created_at);
    const year = date.getFullYear();
    if (!groups[year]) {
      groups[year] = [];
    }
    groups[year].push(article);
  });
  return groups;
});

const formatDate = (dateString) => {
  const date = new Date(dateString);
  return `${date.getMonth() + 1}/${date.getDate()}`;
};
</script>

<style scoped>
.timeline-container {
  padding: var(--spacing-md);
}

.timeline {
  position: relative;
  border-left: 2px solid var(--surface-border);
  margin-left: 20px;
  padding-left: 30px;
}

.timeline-year {
  margin-bottom: var(--spacing-xl);
}

.year-title {
  position: relative;
  left: -45px;
  color: var(--primary-color);
  font-size: 2rem;
  font-family: var(--font-heading);
  text-shadow: var(--neon-text-shadow);
  margin-bottom: var(--spacing-lg);
  background: var(--background-color);
  display: inline-block;
  padding: 0 10px;
}

.timeline-item {
  position: relative;
  margin-bottom: var(--spacing-lg);
}

.timeline-dot {
  position: absolute;
  left: -36px;
  top: 5px;
  width: 12px;
  height: 12px;
  background: var(--primary-color);
  border-radius: 50%;
  box-shadow: 0 0 10px var(--primary-color);
}

.timeline-date {
  position: absolute;
  left: -80px;
  top: 2px;
  color: var(--text-secondary);
  font-family: var(--font-heading);
  font-size: 0.9rem;
  width: 40px;
  text-align: right;
}

.timeline-content {
  padding: var(--spacing-md);
  margin-left: 0;
  border: 1px solid var(--surface-border);
  transition: all 0.3s ease;
}

.timeline-content:hover {
  border-color: var(--primary-color);
  transform: translateX(5px);
  box-shadow: 0 0 15px rgba(0, 243, 255, 0.1);
}

.timeline-link {
  text-decoration: none;
}

.timeline-link h3 {
  margin: 0;
  color: var(--text-primary);
  font-size: 1.2rem;
  transition: color 0.3s;
}

.timeline-link:hover h3 {
  color: var(--primary-color);
}

.timeline-meta {
  margin-top: var(--spacing-sm);
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.category-tag {
  font-size: 0.8rem;
  color: var(--secondary-color);
  border: 1px solid var(--secondary-color);
  padding: 2px 6px;
  border-radius: 2px;
  text-transform: uppercase;
}

.tag-pill {
  font-size: 0.8rem;
  color: var(--text-secondary);
  background: rgba(255, 255, 255, 0.05);
  padding: 2px 6px;
  border-radius: 2px;
}
</style>
