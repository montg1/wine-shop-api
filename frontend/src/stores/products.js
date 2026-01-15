import { defineStore } from 'pinia'
import api from '../services/api'

export const useProductStore = defineStore('products', {
    state: () => ({
        products: [],
        currentProduct: null,
        loading: false,
        meta: { total: 0, page: 1, limit: 10, search: '', category: '' }
    }),

    actions: {
        async fetchProducts(page = 1, limit = 10, search = '', category = '') {
            this.loading = true
            try {
                const params = { page, limit }
                if (search) params.search = search
                if (category) params.category = category

                const response = await api.get('/products', { params })
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
