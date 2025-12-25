import { createRouter, createWebHistory } from 'vue-router'
import SearchView from '../views/SearchView.vue'
import LoginView from '../views/LoginView.vue'
import HistoryView from '../views/HistoryView.vue'

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/',
            name: 'home',
            component: SearchView
        },
        {
            path: '/login',
            name: 'login',
            component: LoginView
        },
        {
            path: '/history',
            name: 'history',
            component: HistoryView
        }
    ]
})

export default router
