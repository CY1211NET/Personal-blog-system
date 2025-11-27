<template>
  <div class="about-page fade-in">
    <div class="about-container">
      <h1 class="page-title text-center">{{ $t('nav.about') }}</h1>
      <AboutWidget class="about-widget-large" />
    </div>
  </div>
</template>

<script setup>
import AboutWidget from '../components/AboutWidget.vue';
</script>

<style scoped>
.about-page {
  display: flex;
  justify-content: center;
  padding: var(--spacing-xl) 0;
}

.about-container {
  width: 100%;
  max-width: 800px;
}

.page-title {
  color: var(--primary-color);
  text-shadow: var(--neon-text-shadow);
  margin-bottom: var(--spacing-xl);
}

/* Override widget styles for full page */
:deep(.about-widget) {
  padding: 3rem;
  background: rgba(5, 5, 16, 0.8);
  backdrop-filter: blur(10px);
  border: 1px solid var(--primary-color);
  box-shadow: 0 0 30px rgba(0, 243, 255, 0.1);
}

:deep(.avatar-container) {
  width: 150px;
  height: 150px;
}

:deep(.username) {
  font-size: 2rem;
  margin-top: 1rem;
}

:deep(.bio) {
  font-size: 1.1rem;
  max-width: 600px;
  margin-left: auto;
  margin-right: auto;
}
</style>
