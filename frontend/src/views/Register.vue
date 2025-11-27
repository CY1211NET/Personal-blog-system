<template>
  <div class="auth-container fade-in">
    <div class="card auth-card">
      <h2 class="text-center">{{ $t('auth.register_title') }}</h2>
      
      <!-- Registration Closed Message -->
      <div v-if="!registrationAllowed" class="registration-closed">
        <div class="closed-icon">🔒</div>
        <h3>注册已关闭</h3>
        <p>这是一个单用户博客系统，目前已有管理员账户。</p>
        <p>如需访问，请联系管理员。</p>
        <router-link to="/login" class="btn btn-primary">
          前往登录
        </router-link>
      </div>

      <!-- Registration Form -->
      <form v-else @submit.prevent="handleRegister">
        <div class="form-group">
          <label>{{ $t('auth.username') }}</label>
          <input type="text" v-model="username" required placeholder="CODENAME" />
        </div>
        <div class="form-group">
          <label>{{ $t('auth.email') }}</label>
          <input type="email" v-model="email" required placeholder="COMM LINK" />
        </div>
        <div class="form-group">
          <label>{{ $t('auth.password') }}</label>
          <input type="password" v-model="password" required placeholder="ACCESS CODE" />
        </div>
        <button type="submit" class="btn btn-primary w-full" :disabled="loading">
          {{ loading ? $t('common.loading') : $t('auth.register_title') }}
        </button>
      </form>
      
      <p class="auth-link">
        <router-link to="/login">{{ $t('auth.has_account') }}</router-link>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import api from '../api/axios';
import { useRouter } from 'vue-router';
import { useToast } from '../composables/useToast';

const username = ref('');
const email = ref('');
const password = ref('');
const router = useRouter();
const toast = useToast();
const loading = ref(false);
const registrationAllowed = ref(true); // Default to true, will be updated

const checkRegistrationStatus = async () => {
  try {
    const response = await api.get('/registration-status');
    registrationAllowed.value = response.data.registration_allowed;
    console.log('Registration status:', response.data);
  } catch (err) {
    console.error('Failed to check registration status:', err);
    // If check fails, allow registration (fail-open)
    registrationAllowed.value = true;
  }
};

const handleRegister = async () => {
  loading.value = true;
  try {
    await api.post('/auth/register', {
      username: username.value,
      email: email.value,
      password: password.value
    });
    toast.success('IDENTITY CREATED. PLEASE AUTHENTICATE.');
    router.push('/login');
  } catch (err) {
    toast.error('REGISTRATION FAILED: ' + (err.response?.data?.error || err.message));
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  checkRegistrationStatus();
});
</script>

<style scoped>
.auth-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 80vh;
}

.auth-card {
  width: 100%;
  max-width: 450px;
  border: 1px solid var(--secondary-color);
}

.auth-card h2 {
  color: var(--secondary-color);
  text-shadow: 0 0 5px var(--secondary-glow);
}

.form-group {
  margin-bottom: var(--spacing-lg);
}

.form-group label {
  display: block;
  margin-bottom: var(--spacing-xs);
  color: var(--secondary-color);
  font-family: var(--font-heading);
  font-size: 0.9rem;
  letter-spacing: 1px;
}

.btn-primary {
  border-color: var(--secondary-color);
  color: var(--secondary-color);
  background: rgba(188, 19, 254, 0.1);
}

.btn-primary:hover {
  background: var(--secondary-color);
  color: #000;
  box-shadow: 0 0 20px var(--secondary-glow);
}

.auth-link {
  margin-top: var(--spacing-lg);
  text-align: center;
  font-size: 0.9rem;
  border-top: 1px solid var(--surface-border);
  padding-top: var(--spacing-md);
}

.registration-closed {
  text-align: center;
  padding: 2rem;
  animation: fadeIn 0.5s;
}

.closed-icon {
  font-size: 4rem;
  margin-bottom: 1rem;
  animation: pulse 2s infinite;
}

.registration-closed h3 {
  color: var(--secondary-color);
  margin-bottom: 1rem;
  font-size: 1.5rem;
}

.registration-closed p {
  color: var(--text-secondary);
  margin-bottom: 0.5rem;
  line-height: 1.6;
}

.registration-closed .btn {
  margin-top: 2rem;
  min-width: 200px;
}

@keyframes pulse {
  0%, 100% {
    transform: scale(1);
    opacity: 1;
  }
  50% {
    transform: scale(1.1);
    opacity: 0.8;
  }
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
