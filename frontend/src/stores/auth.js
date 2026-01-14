import { defineStore } from 'pinia'
import api from '../services/api'

export const useAuthStore = defineStore('auth', {
    state: () => ({
        user: null,
        token: localStorage.getItem('token') || null
    }),

    getters: {
        isLoggedIn: (state) => !!state.token
    },

    actions: {
        async login(email, password) {
            const response = await api.post('/login', { email, password })
            this.token = response.data.token
            localStorage.setItem('token', this.token)
            return response.data
        },

        async register(email, password) {
            const response = await api.post('/register', { email, password })
            return response.data
        },

        logout() {
            this.token = null
            this.user = null
            localStorage.removeItem('token')
        }
    }
})
