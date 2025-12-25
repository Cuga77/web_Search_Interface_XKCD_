<script setup>
import { RouterLink, RouterView } from 'vue-router'
import StatusBanner from './components/StatusBanner.vue'
</script>

<template>
  <header>
    <div class="wrapper">
      <nav>
        <div class="nav-left">
             <RouterLink to="/" class="logo">SearchEngine</RouterLink>
             <StatusBanner />
        </div>
        <div class="links">
             <RouterLink to="/history" class="nav-link">History</RouterLink>
        </div>
      </nav>
    </div>
  </header>

  <main>
    <RouterView />
  </main>
</template>

<style scoped>
header {
  background-color: #ffffff;
  padding: 1.5rem 0;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05);
  position: sticky;
  top: 0;
  z-index: 1000;
}
.wrapper {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 2rem;
}
nav {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.logo {
  font-weight: 900;
  font-size: 1.5rem;
  color: #2c3e50;
  text-decoration: none;
  letter-spacing: -0.5px;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}
/* Unified Links area handled by dashboard logic mostly */
.links {
  display: flex;
  gap: 2rem;
}
main {
    width: 100%;
}
.nav-left {
    display: flex;
    align-items: center;
    gap: 2rem;
}
.nav-link {
    text-decoration: none;
    color: #666;
    font-weight: 600;
    transition: color 0.2s;
}
.nav-link:hover, .nav-link.router-link-active {
    color: #2c3e50;
}
</style>
