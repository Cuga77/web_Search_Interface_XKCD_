<script setup>
defineProps({
  isAdmin: Boolean,
  loading: Boolean
})
const emit = defineEmits(['login-click', 'logout', 'update', 'drop'])
</script>

<template>
  <div class="sidebar-block auth-block">
     <h3>Administration</h3>
     <template v-if="isAdmin">
        <div class="admin-actions">
            <button @click="emit('update')" :disabled="loading" class="btn-outline update">Update Database</button>
            <button @click="emit('drop')" :disabled="loading" class="btn-outline drop">Drop Database</button>
            <button @click="emit('logout')" class="btn-text">Logout</button>
        </div>
     </template>
     <button v-else @click="emit('login-click')" class="btn-primary full-width">Admin Login</button>
  </div>
</template>

<style scoped>
.sidebar-block {
    background: white;
    padding: 1.5rem;
    border-radius: 12px;
    box-shadow: 0 4px 15px rgba(0,0,0,0.05);
}
.sidebar-block h3 {
    font-size: 1.1rem;
    color: #2c3e50;
    margin-bottom: 1rem;
    border-bottom: 1px solid #eee;
    padding-bottom: 0.5rem;
}
.btn-primary {
    background: #3f51b5;
    color: white;
    width: 100%;
    padding: 0.8rem;
    border-radius: 6px;
    font-weight: 600;
    transition: background 0.2s;
}
.btn-primary:hover { background: #303f9f; }
.admin-actions {
    display: flex;
    flex-direction: column;
    gap: 0.8rem;
}
.btn-outline {
    background: white;
    border: 1px solid #ddd;
    padding: 0.6rem;
    border-radius: 6px;
    color: #555;
    font-weight: 500;
    transition: all 0.2s;
}
.btn-outline.update:hover { border-color: #3498db; color: #3498db; background: #f0f7fb; }
.btn-outline.drop:hover { border-color: #e74c3c; color: #e74c3c; background: #fdf2f2; }
.btn-text {
    background: none;
    color: #7f8c8d;
    padding: 0.5rem;
    font-size: 0.9rem;
}
.btn-text:hover { color: #333; }
</style>
