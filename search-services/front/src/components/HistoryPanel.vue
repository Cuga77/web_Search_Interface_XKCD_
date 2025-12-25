<script setup>
defineProps({
  history: Array
})
const emit = defineEmits(['clear'])
</script>

<template>
  <div class="sidebar-block history-block">
    <div class="block-header">
        <h3>Recent Views</h3>
        <button v-if="history.length" @click="emit('clear')" class="link-btn">Clear</button>
    </div>
    
    <TransitionGroup name="list" tag="div" class="history-list">
        <div v-for="item in history" :key="item.id" class="hist-item">
            <img :src="item.url" class="hist-img" />
            <div class="hist-details">
                <span class="hist-id">#{{ item.id }}</span>
                <a :href="`https://xkcd.com/${item.id}`" target="_blank">Open Concept</a>
            </div>
        </div>
        <div v-if="history.length === 0" class="no-history" key="empty">No items viewed</div>
    </TransitionGroup>
  </div>
</template>

<style scoped>
.sidebar-block {
    background: white;
    padding: 1.5rem;
    border-radius: 12px;
    box-shadow: 0 4px 15px rgba(0,0,0,0.05);
}
.block-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    border-bottom: 1px solid #eee;
    margin-bottom: 1rem;
    padding-bottom: 0.5rem;
}
.block-header h3 { font-size: 1.1rem; color: #2c3e50; margin: 0; }
.link-btn {
    background: none;
    color: #e74c3c;
    font-size: 0.8rem;
    padding: 0;
}
.history-list {
    max-height: 400px;
    overflow-y: auto;
    display: flex;
    flex-direction: column;
    gap: 1rem;
}
.hist-item {
    display: flex;
    gap: 1rem;
    align-items: center;
}
.hist-img {
    width: 50px;
    height: 50px;
    background: #f9f9f9;
    border-radius: 6px;
    object-fit: cover;
    border: 1px solid #eee;
}
.hist-details {
    display: flex;
    flex-direction: column;
    font-size: 0.9rem;
}
.hist-id { font-weight: 700; color: #2c3e50; }
.hist-details a { font-size: 0.8rem; color: #95a5a6; }
.no-history { color: #999; text-align: center; padding: 1rem; }

/* Transitions */
.list-move,
.list-enter-active,
.list-leave-active {
  transition: all 0.3s ease;
}
.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: translateX(20px);
}
.list-leave-active {
  position: absolute;
}
</style>
