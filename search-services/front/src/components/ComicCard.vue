<script setup>
defineProps({
  comic: Object
})
const emit = defineEmits(['click-image'])
</script>

<template>
  <div class="comic-card">
    <div class="card-image" @click="emit('click-image', comic)">
      <img :src="comic.url" loading="lazy" />
    </div>
    <div class="card-content">
      <div class="card-meta">ID #{{ comic.id }}</div>
      <h3 class="card-title">XKCD #{{ comic.id }}</h3>
      <a :href="`https://xkcd.com/${comic.id}`" target="_blank" class="card-link">View Source &rarr;</a>
    </div>
  </div>
</template>

<style scoped>
.comic-card {
    background: white;
    border-radius: 12px;
    overflow: hidden;
    box-shadow: 0 4px 15px rgba(0,0,0,0.05);
    transition: transform 0.2s, box-shadow 0.2s;
    border: 1px solid transparent;
}
.comic-card:hover {
    transform: translateY(-5px);
    box-shadow: 0 8px 25px rgba(0,0,0,0.1);
    border-color: #e0e0e0;
}
.card-image {
    background: #fff;
    height: 250px;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 1rem;
    border-bottom: 1px solid #f0f0f0;
    cursor: pointer;
}
.card-image img {
    max-width: 100%;
    max-height: 100%;
    object-fit: contain;
}
.card-content {
    padding: 1.5rem;
}
.card-meta {
    font-size: 0.75rem;
    color: #95a5a6;
    text-transform: uppercase;
    font-weight: 700;
    margin-bottom: 0.5rem;
}
.card-title {
    font-size: 1.2rem;
    color: #2c3e50;
    margin: 0 0 1rem;
    font-weight: 600;
}
.card-link {
    font-size: 0.9rem;
    color: #3498db;
    font-weight: 600;
}
.card-link:hover { text-decoration: underline; }
</style>
