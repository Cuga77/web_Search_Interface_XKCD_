import { ref, onMounted } from 'vue'

export function useHistory() {
    const history = ref([])

    const loadHistory = () => {
        try {
            history.value = JSON.parse(localStorage.getItem('history') || '[]')
        } catch (e) { history.value = [] }
    }

    const addToHistory = (comic) => {
        let hist = [...history.value]
        hist = hist.filter(c => c.id !== comic.id)
        hist.unshift({ id: comic.id, url: comic.url, date: new Date().toISOString() })
        if (hist.length > 20) hist.pop()
        localStorage.setItem('history', JSON.stringify(hist))
        history.value = hist
    }

    const clearHistory = () => {
        localStorage.removeItem('history')
        history.value = []
    }

    onMounted(() => {
        loadHistory()
    })

    return { history, addToHistory, clearHistory }
}
