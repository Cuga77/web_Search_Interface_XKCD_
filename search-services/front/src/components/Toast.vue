<script setup>
import { ref } from 'vue'

const toasts = ref([])
let idCounter = 0

const addToast = (message, type = 'info') => {
  const id = idCounter++
  toasts.value.push({ id, message, type })
  setTimeout(() => removeToast(id), 3000)
}

const removeToast = (id) => {
  toasts.value = toasts.value.filter(t => t.id !== id)
}

defineExpose({ addToast })
</script>

<template>
  <div class="toast-container">
    <TransitionGroup name="toast">
      <div v-for="toast in toasts" :key="toast.id" :class="['toast', toast.type]">
        {{ toast.message }}
      </div>
    </TransitionGroup>
  </div>
</template>

<style scoped>
.toast-container {
  position: fixed;
  bottom: 2rem;
  right: 2rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  z-index: 2000;
}
.toast {
  padding: 1rem 1.5rem;
  border-radius: 8px;
  background: white;
  box-shadow: 0 5px 15px rgba(0,0,0,0.1);
  font-weight: 500;
  min-width: 250px;
}
.toast.info { border-left: 4px solid #3498db; color: #333; }
.toast.success { border-left: 4px solid #27ae60; color: #2ecc71; }
.toast.error { border-left: 4px solid #e74c3c; color: #e74c3c; }

.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}
.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateY(20px);
}
</style>
