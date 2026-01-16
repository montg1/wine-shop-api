import { defineStore } from 'pinia'
import api from '../services/api'

export const useAuthStore = defineStore('auth', {
    state: () => ({
        user: null,
        token: localStorage.getItem('token') || null
    }),

    getters: {
        isLoggedIn: (state) => !!state.token,
        isAdmin: (state) => state.user?.role === 'admin'
    },

    actions: {
        async login(email, password) {
            const response = await api.post('/login', { email, password })
            this.token = response.data.token
            localStorage.setItem('token', this.token)
            // Fetch user info after login
            await this.fetchUser()
            return response.data
        },

        async register(email, password) {
            const response = await api.post('/register', { email, password })
            return response.data
        },

        async fetchUser() {
            if (!this.token) return
            try {
                const response = await api.get('/me')
                this.user = response.data.data
            } catch (error) {
                console.error('Failed to fetch user:', error)
                // Token might be invalid, clear it
                if (error.response?.status === 401) {
                    this.logout()
                }
            }
        },

        logout() {
            this.token = null
            this.user = null
            localStorage.removeItem('token')
        },

        // Initialize auth state on app load
        async init() {
            if (this.token) {
                await this.fetchUser()
            }
        }
    }
})
