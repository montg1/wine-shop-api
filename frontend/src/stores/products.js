import { defineStore } from 'pinia'
import api from '../services/api'

export const useProductStore = defineStore('products', {
    state: () => ({
        products: [],
        currentProduct: null,
        loading: false,
        meta: { total: 0, page: 1, limit: 10 }
    }),

    actions: {
        async fetchProducts(page = 1, limit = 10) {
            this.loading = true
            try {
                const response = await api.get('/products', { params: { page, limit } })
                this.products = response.data.data
                this.meta = response.data.meta
            } catch (error) {
                console.error('Failed to fetch products:', error)
            } finally {
                this.loading = false
            }
        },

        async fetchProduct(id) {
            this.loading = true
            try {
                const response = await api.get(`/products/${id}`)
                this.currentProduct = response.data.data
            } catch (error) {
                console.error('Failed to fetch product:', error)
            } finally {
                this.loading = false
            }
        }
    }
})
