<template>
  <div class="auth-container fade-in">
    <div class="card auth-card">
      <div class="hologram-effect"></div>
      <h2 class="text-center">{{ $t('auth.login_title') }}</h2>
      <form @submit.prevent="handleLogin">
        <div class="form-group">
          <label>{{ $t('auth.username_or_email') }}</label>
          <input type="text" v-model="loginInput" required placeholder="ENTER USERNAME OR EMAIL" />
        </div>
        <div class="form-group">
          <label>{{ $t('auth.password') }}</label>
          <input type="password" v-model="password" required placeholder="ENTER PASSWORD" />
        </div>
        <button type="submit" class="btn btn-primary w-full" :disabled="loading">
          {{ loading ? $t('common.loading') : $t('auth.login_title') }}
        </button>
      </form>
      <p class="auth-link">
        <router-link to="/register">{{ $t('auth.no_account') }}</router-link>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useAuthStore } from '../stores/auth';
import { useRouter } from 'vue-router';
import { useToast } from '../composables/useToast';

const loginInput = ref('');
const password = ref('');
const authStore = useAuthStore();
const router = useRouter();
const toast = useToast();
const loading = ref(false);

const handleLogin = async () => {
  loading.value = true;
  try {
    await authStore.login({ username: loginInput.value, password: password.value });
    toast.success('ACCESS GRANTED');
    router.push('/');
  } catch (err) {
    toast.error('ACCESS DENIED: ' + (err.response?.data?.error || err.message));
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.auth-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 80vh;
  position: relative;
}

.auth-card {
  width: 100%;
  max-width: 450px;
  border: 1px solid var(--primary-color);
  box-shadow: 0 0 20px rgba(0, 243, 255, 0.1);
}

.form-group {
  margin-bottom: var(--spacing-lg);
}

.form-group label {
  display: block;
  margin-bottom: var(--spacing-xs);
  color: var(--primary-color);
  font-family: var(--font-heading);
  font-size: 0.9rem;
  letter-spacing: 1px;
}

.auth-link {
  margin-top: var(--spacing-lg);
  text-align: center;
  font-size: 0.9rem;
  border-top: 1px solid var(--surface-border);
  padding-top: var(--spacing-md);
}
</style>
