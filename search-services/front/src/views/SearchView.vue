<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useSearch } from '../composables/useSearch'
import { useHistory } from '../composables/useHistory'
import { useAdmin } from '../composables/useAdmin'
import { useStats } from '../composables/useStats'

import ComicCard from '../components/ComicCard.vue'
import HistoryPanel from '../components/HistoryPanel.vue'
import AdminPanel from '../components/AdminPanel.vue'
import Toast from '../components/Toast.vue'

const router = useRouter()
const toastRef = ref(null)

const showToast = (msg, type) => toastRef.value?.addToast(msg, type)

const { phrase, results, loading, error, search } = useSearch()
const { history, addToHistory, clearHistory } = useHistory()
const { isAdmin, adminLoading, checkAuth, logout, updateDB, dropDB } = useAdmin(showToast)
const { stats, status } = useStats()

const handleLoginClick = () => router.push('/login')

</script>

<template>
  <div class="dashboard">
    <Toast ref="toastRef" />

    <!-- HERO / TITLE SECTION -->
    <div class="hero-section">
        <h1 class="hero-title">Comic Search & Solutions</h1>
        <p class="hero-subtitle">Comprehensive XKCD archive indexing and retrieval system</p>
        
        <!-- SEARCH BOX FLOATING -->
        <div class="search-wrapper">
             <div class="search-box">
                <input v-model="phrase" @keyup.enter="search" placeholder="Search knowledge base..." />
                <button @click="search" :disabled="loading">
                    {{ loading ? 'SEARCHING...' : 'SEARCH' }}
                </button>
            </div>
        </div>
    </div>

    <!-- MAIN CONTENT -->
    <div class="main-layout">
        
        <!-- RESULTS GRID -->
        <div class="content-section">
            <div class="section-header">
                <h2>Search Results</h2>
                <div class="system-status">
                    <span :class="['status-dot', status]"></span>
                    <span>System: {{ status }}</span>
                    <span v-if="stats" class="ml-2">({{ stats.comics_total }} items)</span>
                </div>
            </div>

            <div v-if="error" class="error-msg">{{ error }}</div>

            <!-- RESULTS (Transition Group for smooth entry) -->
            <TransitionGroup name="grid" tag="div" class="results-grid" v-if="results">
                <ComicCard 
                    v-for="comic in results.comics" 
                    :key="comic.id" 
                    :comic="comic" 
                    @click-image="addToHistory" 
                />
            </TransitionGroup>
            
             <div v-else-if="loading" class="grid-skeleton">
                <div class="skeleton-card" v-for="i in 3" :key="i"></div>
             </div>
             
             <div v-else class="empty-state">
                Enter keywords to begin searching the archive.
             </div>
        </div>

        <!-- RIGHT SIDEBAR -->
        <aside class="sidebar">
            <AdminPanel 
                :isAdmin="isAdmin" 
                :loading="adminLoading"
                @login-click="handleLoginClick"
                @logout="logout"
                @update="updateDB"
                @drop="dropDB"
            />

            <HistoryPanel 
                :history="history" 
                @clear="clearHistory" 
            />
        </aside>

    </div>
  </div>
</template>

<style scoped>
.dashboard {
    min-height: 100vh;
    background: #f5f7fa;
    color: #333;
    font-family: 'Inter', sans-serif;
}

/* HERO */
.hero-section {
    background: #ffffff;
    padding: 3rem 2rem 5rem;
    text-align: center;
    border-bottom: 1px solid #eee;
    position: relative;
    margin-bottom: 3rem;
}
.hero-title {
    font-size: 2.5rem;
    font-weight: 800;
    color: #2c3e50;
    margin-bottom: 0.5rem;
    letter-spacing: -0.5px;
}
.hero-subtitle {
    color: #666;
    font-size: 1.1rem;
    margin-bottom: 2.5rem;
}
.search-wrapper {
    position: absolute;
    bottom: -25px; 
    left: 0;
    right: 0;
    display: flex;
    justify-content: center;
}
.search-box {
    background: white;
    padding: 0.5rem;
    border-radius: 8px;
    box-shadow: 0 10px 25px rgba(0,0,0,0.1);
    display: flex;
    width: 600px;
    max-width: 90%;
    border: 1px solid #e0e0e0;
}
.search-box input {
    flex: 1;
    border: none;
    padding: 1rem;
    font-size: 1rem;
    outline: none;
    color: #333;
}
.search-box button {
    background: #2c3e50;
    color: white;
    padding: 0 2rem;
    border-radius: 6px;
    font-weight: 700;
    font-size: 0.9rem;
    transition: background 0.2s;
}
.search-box button:hover {
    background: #34495e;
}

/* LAYOUT */
.main-layout {
    max-width: 1400px;
    margin: 0 auto;
    padding: 0 2rem 4rem;
    display: grid;
    grid-template-columns: 1fr 300px;
    gap: 3rem;
}

/* SECTION HEADER */
.section-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-end;
    margin-bottom: 2rem;
    border-bottom: 2px solid #e0e0e0;
    padding-bottom: 0.5rem;
}
.section-header h2 {
    font-size: 1.8rem;
    color: #2c3e50;
    font-weight: 700;
}
.system-status {
    font-size: 0.9rem;
    color: #7f8c8d;
    display: flex;
    align-items: center;
    gap: 0.5rem;
}
.status-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: #ccc;
}
.status-dot.running { background: #27ae60; box-shadow: 0 0 5px #27ae60; }
.status-dot.idle { background: #95a5a6; }

/* GRID */
.results-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 2rem;
}
.empty-state {
    text-align: center;
    padding: 4rem;
    color: #95a5a6;
    font-size: 1.1rem;
}
.error-msg {
    color: #e74c3c;
    background: #fdf2f2;
    padding: 1rem;
    border-radius: 8px;
    margin-bottom: 2rem;
    text-align: center;
}

/* SIDEBAR */
.sidebar {
    display: flex;
    flex-direction: column;
    gap: 2rem;
}

/* SKELETON */
.grid-skeleton {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 2rem;
}
.skeleton-card {
  height: 350px;
  background: white;
  border-radius: 12px;
  animation: pulse 1.5s infinite ease-in-out;
}
@keyframes pulse {
  0% { opacity: 0.6; background: #e0e0e0; }
  50% { opacity: 1; background: #f5f5f5; }
  100% { opacity: 0.6; background: #e0e0e0; }
}

/* TRANSITIONS */
.grid-enter-active,
.grid-leave-active {
  transition: all 0.5s ease;
}
.grid-enter-from,
.grid-leave-to {
  opacity: 0;
  transform: translateY(20px);
}
</style>
