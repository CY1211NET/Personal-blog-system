<template>
  <div class="profile-container">
    <h2>User Profile</h2>
    <div class="profile-card">
      <div class="avatar-section">
        <img :src="avatarUrl || 'https://via.placeholder.com/150'" alt="Avatar" class="avatar" />
        <input type="file" @change="uploadAvatar" accept="image/*" />
        <p v-if="uploading">Uploading...</p>
      </div>
      <div class="info-section">
        <p><strong>Username:</strong> {{ user?.username }}</p>
        <p><strong>Email:</strong> {{ user?.email }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import api from '../api/axios';
import { useAuthStore } from '../stores/auth';

const authStore = useAuthStore();
const user = ref(null);
const avatarUrl = ref('');
const uploading = ref(false);

const fetchProfile = async () => {
  // Assuming we have a GetProfile endpoint or store has user info
  // For now, let's assume authStore has basic info, but we might need to fetch fresh
  // Since we didn't implement GetProfile fully in frontend store, let's fetch it
  try {
    const response = await api.get('/user/profile');
    user.value = response.data;
    avatarUrl.value = response.data.avatar_url;
  } catch (err) {
    console.error('Failed to load profile');
  }
};

const uploadAvatar = async (event) => {
  const file = event.target.files[0];
  if (!file) return;

  uploading.value = true;
  const formData = new FormData();
  formData.append('file', file);

  try {
    // 1. Upload file
    const uploadResp = await api.post('/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    });
    const newAvatarUrl = uploadResp.data.url;

    // 2. Update user profile
    await api.put('/user/profile', { avatar_url: newAvatarUrl });
    
    avatarUrl.value = newAvatarUrl;
    alert('Avatar updated successfully!');
  } catch (err) {
    alert('Failed to upload avatar');
    console.error(err);
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
  padding: 20px;
}
.profile-card {
  display: flex;
  gap: 30px;
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}
.avatar {
  width: 150px;
  height: 150px;
  border-radius: 50%;
  object-fit: cover;
  margin-bottom: 10px;
}
.avatar-section {
  display: flex;
  flex-direction: column;
  align-items: center;
}
</style>
