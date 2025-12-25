<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const username = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

const login = async () => {
    loading.value = true
    error.value = ''
    try {
        const res = await fetch('/api/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name: username.value, password: password.value })
        })
        
        if (!res.ok) {
             throw new Error("Invalid credentials")
        }

        const token = await res.text()
        if (token) {
            localStorage.setItem('token', token)
            router.push('/admin')
        } else {
            throw new Error("No token received")
        }
    } catch (e) {
        error.value = e.message
    } finally {
        loading.value = false
    }
}
</script>

<template>
  <div class="login-container">
    <div class="login-card">
        <h2>Admin Access</h2>
        <form @submit.prevent="login">
            <div class="input-group">
                <label>Username</label>
                <input v-model="username" type="text" required />
            </div>
            <div class="input-group">
                <label>Password</label>
                <input v-model="password" type="password" required />
            </div>
            <button type="submit" :disabled="loading">
                {{ loading ? 'Signing in...' : 'Sign In' }}
            </button>
            <p v-if="error" class="error">{{ error }}</p>
        </form>
    </div>
  </div>
</template>

<style scoped>
.login-container {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 60vh;
}
.login-card {
    background: #1e1e1e;
    padding: 2rem;
    border-radius: 1rem;
    width: 100%;
    max-width: 400px;
    border: 1px solid #333;
    box-shadow: 0 10px 30px rgba(0,0,0,0.5);
}
h2 {
    text-align: center;
    margin-bottom: 2rem;
    color: #fff;
}
.input-group {
    margin-bottom: 1.5rem;
}
label {
    display: block;
    margin-bottom: 0.5rem;
    color: #aaa;
}
input {
    width: 100%;
    padding: 0.8rem;
    border-radius: 0.5rem;
    border: 1px solid #444;
    background: #2a2a2a;
    color: white;
    box-sizing: border-box; /* Fix padding issue */
}
input:focus {
    border-color: #646cff;
    outline: none;
}
button {
    width: 100%;
    padding: 1rem;
    background: #646cff;
    color: white;
    border: none;
    border-radius: 0.5rem;
    font-weight: bold;
    cursor: pointer;
    transition: background 0.2s;
}
button:hover:not(:disabled) {
    background: #535bf2;
}
button:disabled {
    opacity: 0.7;
}
.error {
    color: #ff4444;
    text-align: center;
    margin-top: 1rem;
}
</style>
