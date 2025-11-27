<template>
  <div class="profile-container fade-in">
    <div class="card profile-card">
      <h2 class="text-center">{{ $t('profile.title') }}</h2>
      <div class="profile-content">
        <div class="avatar-section">
          <div class="avatar-wrapper">
            <img :src="avatarUrl || 'https://via.placeholder.com/150'" alt="Avatar" class="avatar" />
            <div class="avatar-overlay">
              <label for="avatar-upload" class="upload-label">
                <span v-if="uploading">{{ $t('profile.uploading') }}</span>
                <span v-else>📷</span>
              </label>
              <input type="file" id="avatar-upload" @change="uploadAvatar" accept="image/*" class="hidden-input" />
            </div>
          </div>
        </div>
        <div class="info-section">
          <div class="info-item">
            <label>{{ $t('auth.username') }}</label>
            <p>{{ user?.username }}</p>
          </div>
          <div class="info-item">
            <label>{{ $t('auth.email') }}</label>
            <p>{{ user?.email }}</p>
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

const authStore = useAuthStore();
const toast = useToast();
const user = ref(null);
const avatarUrl = ref('');
const uploading = ref(false);

const fetchProfile = async () => {
  try {
    const response = await api.get('/user/profile');
    user.value = response.data;
    avatarUrl.value = response.data.avatar_url;
  } catch (err) {
    toast.error('Failed to load profile');
  }
};

const uploadAvatar = async (event) => {
  const file = event.target.files[0];
  if (!file) return;

  uploading.value = true;
  const formData = new FormData();
  formData.append('file', file);

  try {
    const uploadResp = await api.post('/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    });
    const newAvatarUrl = uploadResp.data.url;

    await api.put('/user/profile', { avatar_url: newAvatarUrl });
    
    avatarUrl.value = newAvatarUrl;
    toast.success('Avatar updated successfully!');
  } catch (err) {
    toast.error('Failed to upload avatar');
  } finally {
    uploading.value = false;
  }
};

onMounted(() => {
  fetchProfile();
});
</script>

<style scoped>
.profile-container {
  max-width: 600px;
  margin: 0 auto;
}

.profile-card {
  padding: var(--spacing-xl);
  border: 1px solid var(--primary-color);
  box-shadow: 0 0 30px rgba(0, 243, 255, 0.1);
}

.text-center {
  text-align: center;
  margin-bottom: var(--spacing-xl);
  color: var(--primary-color);
  text-shadow: var(--neon-text-shadow);
}

.profile-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--spacing-xl);
}

.avatar-wrapper {
  position: relative;
  width: 150px;
  height: 150px;
  border-radius: 50%;
  overflow: hidden;
  border: 2px solid var(--primary-color);
  box-shadow: 0 0 20px var(--primary-glow);
}

.avatar {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  justify-content: center;
  align-items: center;
  opacity: 0;
  transition: opacity var(--transition-fast);
}

.avatar-wrapper:hover .avatar-overlay {
  opacity: 1;
}

.upload-label {
  color: var(--primary-color);
  cursor: pointer;
  font-size: 1.5rem;
  text-shadow: var(--neon-text-shadow);
}

.hidden-input {
  display: none;
}

.info-section {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.info-item {
  background-color: rgba(0, 0, 0, 0.3);
  padding: var(--spacing-md);
  border: 1px solid var(--surface-border);
  position: relative;
  overflow: hidden;
}

.info-item::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  height: 100%;
  width: 3px;
  background: var(--secondary-color);
}

.info-item label {
  display: block;
  font-size: 0.8rem;
  color: var(--secondary-color);
  margin-bottom: var(--spacing-xs);
  font-weight: 500;
  font-family: var(--font-heading);
  letter-spacing: 1px;
  text-transform: uppercase;
}

.info-item p {
  font-size: 1.2rem;
  font-weight: 600;
  color: var(--text-primary);
  font-family: var(--font-body);
}
</style>
