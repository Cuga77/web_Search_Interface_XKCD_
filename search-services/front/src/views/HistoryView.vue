<script setup>
import { useHistory } from '../composables/useHistory'
import ComicCard from '../components/ComicCard.vue'

const { history, clearHistory } = useHistory()
</script>

<template>
  <div class="history-view">
    <div class="header">
      <h1>Reading History</h1>
      <button @click="clearHistory" class="clear-btn" v-if="history.length">Clear History</button>
    </div>

    <div v-if="history.length === 0" class="empty-state">
      No history yet. Go explore some comics!
    </div>

    <div class="grid" v-else>
      <div v-for="item in history" :key="item.id" class="history-card">
          <!-- We reuse ComicCard structure manually since history stores partial data -->
          <div class="history-image">
             <img :src="item.url" loading="lazy">
          </div>
          <div class="meta">
              <h3>XKCD #{{ item.id }}</h3>
              <div class="date">{{ new Date(item.date).toLocaleString() }}</div>
              <a :href="`https://xkcd.com/${item.id}`" target="_blank">View Original</a>
          </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.history-view {
  max-width: 1400px;
  margin: 0 auto;
  padding: 3rem 2rem;
}
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 3rem;
  border-bottom: 2px solid #eee;
  padding-bottom: 1rem;
}
.header h1 { color: #2c3e50; }
.clear-btn {
  background: #e74c3c;
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  cursor: pointer;
}
.grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 2rem;
}
.empty-state {
  text-align: center;
  color: #7f8c8d;
  font-size: 1.2rem;
  margin-top: 4rem;
}
.history-card {
  background: white;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 4px 6px rgba(0,0,0,0.05);
}
.history-image {
    height: 200px;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 1rem;
    border-bottom: 1px solid #eee;
}
.history-image img {
    max-width: 100%;
    max-height: 100%;
}
.meta {
    padding: 1rem;
}
.meta h3 { margin: 0 0 0.5rem; }
.date { font-size: 0.8rem; color: #95a5a6; margin-bottom: 0.5rem; }
.meta a { color: #3498db; text-decoration: none; font-size: 0.9rem; }
</style>
