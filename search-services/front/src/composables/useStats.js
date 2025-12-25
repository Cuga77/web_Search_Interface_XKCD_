import { ref, onMounted, onUnmounted } from 'vue'

export function useStats() {
    const stats = ref(null)
    const status = ref('unknown')
    let intervalId = null

    const fetchStats = async () => {
        try {
            const res1 = await fetch('/api/db/stats')
            if (res1.ok) stats.value = await res1.json()
            const res2 = await fetch('/api/db/status')
            if (res2.ok) {
                const d = await res2.json()
                status.value = d.status
            }
        } catch (_) { }
    }

    onMounted(() => {
        fetchStats()
        intervalId = setInterval(fetchStats, 5000)
    })

    onUnmounted(() => {
        if (intervalId) clearInterval(intervalId)
    })

    return { stats, status, fetchStats }
}
