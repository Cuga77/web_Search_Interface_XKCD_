import { ref, onMounted } from 'vue'

export function useAdmin(showToast) {
    const isAdmin = ref(false)
    const adminLoading = ref(false)

    const getToken = () => localStorage.getItem('token')

    const checkAuth = () => {
        isAdmin.value = !!getToken()
    }

    const apiCall = async (url, options = {}) => {
        const token = getToken()
        const headers = { ...options.headers }
        if (token) headers['Authorization'] = `Bearer ${token}`

        const res = await fetch(url, { ...options, headers })
        if (res.status === 401) {
            logout()
            return null
        }
        return res
    }

    const logout = () => {
        localStorage.removeItem('token')
        isAdmin.value = false
        if (showToast) showToast("Logged out successfully", "info")
    }

    const updateDB = async (callback) => {
        adminLoading.value = true
        try {
            const res = await apiCall('/api/db/update', { method: 'POST' })
            if (res && res.ok) {
                if (showToast) showToast("Database Update Triggered", "success")
                if (callback) callback()
            } else if (res) {
                throw new Error("Update failed")
            }
        } catch (e) {
            if (showToast) showToast(`Error: ${e.message}`, "error")
        } finally {
            adminLoading.value = false
        }
    }

    const dropDB = async (callback) => {
        if (!confirm("Are you sure? This will delete all data.")) return
        adminLoading.value = true
        try {
            const res = await apiCall('/api/db', { method: 'DELETE' })
            if (res && res.ok) {
                if (showToast) showToast("Database Dropped", "success")
                if (callback) callback()
            } else if (res) {
                throw new Error("Drop failed")
            }
        } catch (e) {
            if (showToast) showToast(`Error: ${e.message}`, "error")
        } finally {
            adminLoading.value = false
        }
    }

    onMounted(() => {
        checkAuth()
    })

    return { isAdmin, adminLoading, checkAuth, logout, updateDB, dropDB }
}
