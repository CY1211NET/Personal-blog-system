<template>
  <transition-group name="toast" tag="div" class="toast-container">
    <div 
      v-for="toast in toasts" 
      :key="toast.id" 
      class="toast" 
      :class="toast.type"
    >
      {{ toast.message }}
    </div>
  </transition-group>
</template>

<script setup>
import { useToast } from '../composables/useToast';

const { toasts } = useToast();
</script>

<style scoped>
.toast-container {
  position: fixed;
  bottom: 20px;
  right: 20px;
  z-index: 9999;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.toast {
  padding: 12px 20px;
  border-radius: var(--radius-md);
  color: white;
  font-weight: 500;
  box-shadow: var(--shadow-lg);
  min-width: 250px;
}

.toast.success {
  background-color: var(--success-color);
}

.toast.error {
  background-color: var(--danger-color);
}

.toast.info {
  background-color: var(--primary-color);
}

/* Transitions */
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from {
  opacity: 0;
  transform: translateX(30px);
}

.toast-leave-to {
  opacity: 0;
  transform: translateY(30px);
}
</style>
