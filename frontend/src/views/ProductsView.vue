<template>
  <div class="products-page">
    <div class="page-header">
      <h1>Our Collection</h1>
      <p>Carefully selected wines from exceptional vineyards</p>
    </div>
    
    <div v-if="productStore.loading" class="loading">Loading wines...</div>
    
    <div v-else class="products-grid">
      <div 
        v-for="product in productStore.products" 
        :key="product.ID" 
        class="product-card"
        @click="$router.push(`/products/${product.ID}`)"
      >
        <div class="product-image">
          <img :src="wineBottleImg" :alt="product.name" />
        </div>
        <div class="product-info">
          <span class="category">{{ product.category }}</span>
          <h3>{{ product.name }}</h3>
          <p class="price">${{ product.price.toFixed(2) }}</p>
          <button 
            class="btn btn-outline" 
            @click.stop="addToCart(product.ID)"
          >
            Add to Cart
          </button>
        </div>
      </div>
    </div>
    
    <div v-if="productStore.products.length === 0 && !productStore.loading" class="empty">
      <p>No wines available at the moment.</p>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useProductStore } from '../stores/products'
import { useCartStore } from '../stores/cart'
import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'
import wineBottleImg from '../assets/images/wine-bottle.png'

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
  padding: 60px 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  text-align: center;
  margin-bottom: 60px;
}

.page-header h1 {
  font-size: 2.8rem;
  color: var(--primary);
  margin-bottom: 10px;
}

.page-header p {
  color: var(--text-muted);
  font-size: 1.1rem;
}

.products-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 40px;
}

.product-card {
  background: var(--card-bg);
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  transition: transform 0.3s, box-shadow 0.3s;
  border: 1px solid var(--border);
}

.product-card:hover {
  transform: translateY(-5px);
  box-shadow: var(--shadow-lg);
}

.product-image {
  height: 280px;
  background: var(--bg-warm);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 30px;
}

.product-image img {
  max-height: 100%;
  max-width: 100%;
  object-fit: contain;
}

.product-info {
  padding: 25px;
  text-align: center;
}

.category {
  font-size: 0.8rem;
  text-transform: uppercase;
  letter-spacing: 2px;
  color: var(--text-muted);
}

.product-info h3 {
  font-size: 1.4rem;
  color: var(--text);
  margin: 8px 0;
}

.price {
  font-size: 1.3rem;
  color: var(--primary);
  margin-bottom: 20px;
}

.loading, .empty {
  text-align: center;
  padding: 60px;
  color: var(--text-muted);
  font-size: 1.1rem;
}
</style>
