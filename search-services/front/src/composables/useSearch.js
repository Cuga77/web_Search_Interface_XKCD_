import { ref } from 'vue'

export function useSearch() {
    const phrase = ref('')
    const results = ref(null)
    const loading = ref(false)
    const error = ref('')

    const search = async () => {
        if (!phrase.value) return
        loading.value = true
        error.value = ''
        results.value = null

        try {
            const res = await fetch(`/api/search?phrase=${encodeURIComponent(phrase.value)}`)
            if (!res.ok) throw new Error(`Error: ${res.statusText}`)
            results.value = await res.json()
        } catch (e) {
            error.value = e.message
        } finally {
            loading.value = false
        }
    }

    return { phrase, results, loading, error, search }
}
