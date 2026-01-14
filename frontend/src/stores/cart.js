import { defineStore } from 'pinia'
import api from '../services/api'

export const useCartStore = defineStore('cart', {
    state: () => ({
        items: [],
        loading: false
    }),

    getters: {
        totalItems: (state) => state.items.reduce((sum, item) => sum + item.quantity, 0),
        totalPrice: (state) => state.items.reduce((sum, item) => sum + (item.product.price * item.quantity), 0)
    },

    actions: {
        async fetchCart() {
            this.loading = true
            try {
                const response = await api.get('/cart')
                this.items = response.data.data?.items || []
            } catch (error) {
                console.error('Failed to fetch cart:', error)
            } finally {
                this.loading = false
            }
        },

        async addToCart(productId, quantity = 1) {
            await api.post('/cart', { product_id: productId, quantity })
            await this.fetchCart()
        },

        async checkout() {
            const response = await api.post('/orders')
            this.items = []
            return response.data
        }
    }
})
