<template>
  <div class="products-page">
    <h1>Our Wines</h1>
    
    <div v-if="productStore.loading" class="loading">Loading wines...</div>
    
    <div v-else class="products-grid">
      <div 
        v-for="product in productStore.products" 
        :key="product.ID" 
        class="product-card"
        @click="$router.push(`/products/${product.ID}`)"
      >
        <div class="product-image">🍷</div>
        <div class="product-info">
          <h3>{{ product.name }}</h3>
          <p class="category">{{ product.category }}</p>
          <p class="price">${{ product.price.toFixed(2) }}</p>
          <button 
            class="btn btn-primary" 
            @click.stop="addToCart(product.ID)"
          >
            Add to Cart
          </button>
        </div>
      </div>
    </div>
    
    <div v-if="productStore.products.length === 0 && !productStore.loading" class="empty">
      No wines available at the moment.
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useProductStore } from '../stores/products'
import { useCartStore } from '../stores/cart'
import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'

const productStore = useProductStore()
const cartStore = useCartStore()
const authStore = useAuthStore()
const router = useRouter()

onMounted(() => {
  productStore.fetchProducts()
})

const addToCart = async (productId) => {
  if (!authStore.isLoggedIn) {
    router.push('/login')
    return
  }
  await cartStore.addToCart(productId, 1)
  alert('Added to cart!')
}
</script>

<style scoped>
.products-page {
  padding: 40px 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.products-page h1 {
  text-align: center;
  margin-bottom: 40px;
  color: var(--primary);
}

.products-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 30px;
}

.product-card {
  background: var(--card-bg);
  border-radius: 16px;
  overflow: hidden;
  cursor: pointer;
  transition: transform 0.3s, box-shadow 0.3s;
}

.product-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 30px rgba(139, 69, 102, 0.3);
}

.product-image {
  height: 180px;
  background: linear-gradient(135deg, #3d2a4a 0%, #5a3a6a 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 5rem;
}

.product-info {
  padding: 20px;
}

.product-info h3 {
  margin-bottom: 5px;
  color: #fff;
}

.category {
  color: #999;
  font-size: 0.9rem;
  margin-bottom: 10px;
}

.price {
  font-size: 1.5rem;
  font-weight: bold;
  color: var(--primary);
  margin-bottom: 15px;
}

.loading, .empty {
  text-align: center;
  padding: 60px;
  color: #888;
  font-size: 1.2rem;
}
</style>
