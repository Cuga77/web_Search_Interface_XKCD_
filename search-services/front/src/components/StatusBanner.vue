<script setup>
import { useStats } from '../composables/useStats'

const { stats, status } = useStats()
</script>

<template>
  <div class="status-banner" v-if="stats">
    <div class="wrapper">
      <div class="stat-item">
        <span class="label">Status:</span>
        <span :class="['value', 'status-dot', status]">{{ status }}</span>
      </div>
      <div class="stat-item">
        <span class="label">Comics Fetched:</span>
        <span class="value">{{ stats.comics_fetched }}</span>
      </div>
      <div class="stat-item">
        <span class="label">Total Words:</span>
        <span class="value">{{ stats.words_total }}</span>
      </div>
      <div class="stat-item">
        <span class="label">Indexed Words:</span>
        <span class="value">{{ stats.words_unique }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.status-banner {
  background: #2c3e50;
  color: white;
  padding: 0.5rem 0;
  font-size: 0.85rem;
}
.wrapper {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 2rem;
  display: flex;
  justify-content: flex-end;
  gap: 2rem;
}
.stat-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}
.label {
  color: #bdc3c7;
}
.value {
  font-weight: 700;
  color: white;
}
.status-dot {
  text-transform: uppercase;
}
.status-dot.running { color: #2ecc71; }
.status-dot.idle { color: #95a5a6; }
</style>
