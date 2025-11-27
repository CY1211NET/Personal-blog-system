<template>
  <nav class="nav-bar">
    <router-link to="/" style="font-size: 1.5rem; font-weight: bold; color: var(--text-color);">My Blog</router-link>
    <div class="nav-links">
      <router-link to="/">Home</router-link>
      <template v-if="authStore.isAuthenticated">
        <router-link to="/articles/new">Write</router-link>
        <a href="#" @click.prevent="logout">Logout</a>
      </template>
      <template v-else>
        <router-link to="/login">Login</router-link>
        <router-link to="/register">Register</router-link>
      </template>
    </div>
  </nav>
  <div class="container">
    <router-view></router-view>
  </div>
</template>

<script setup>
import { useAuthStore } from './stores/auth';
import { useRouter } from 'vue-router';

const authStore = useAuthStore();
const router = useRouter();

const logout = () => {
  authStore.logout();
  router.push('/login');
};
</script>

<style scoped>
.logo {
  height: 6em;
  padding: 1.5em;
  will-change: filter;
  transition: filter 300ms;
}
.logo:hover {
  filter: drop-shadow(0 0 2em #646cffaa);
}
.logo.vue:hover {
  filter: drop-shadow(0 0 2em #42b883aa);
}
</style>
