<template>
  <div class="product-detail">
    <div v-if="productStore.loading" class="loading">Loading...</div>
    
    <div v-else-if="product" class="detail-content">
      <div class="product-hero">🍷</div>
      
      <div class="product-info">
        <span class="category-badge">{{ product.category }}</span>
        <h1>{{ product.name }}</h1>
        <p class="description">{{ product.description }}</p>
        
        <div class="meta">
          <span class="price">${{ product.price.toFixed(2) }}</span>
          <span class="stock">{{ product.stock }} in stock</span>
        </div>
        
        <div class="actions">
          <input v-model.number="quantity" type="number" min="1" :max="product.stock" />
          <button class="btn btn-primary" @click="handleAddToCart">Add to Cart</button>
        </div>
        
        <router-link to="/products" class="back-link">← Back to Wines</router-link>
      </div>
    </div>
    
    <div v-else class="not-found">Wine not found.</div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useProductStore } from '../stores/products'
import { useCartStore } from '../stores/cart'
import { useAuthStore } from '../stores/auth'

const route = useRoute()
const router = useRouter()
const productStore = useProductStore()
const cartStore = useCartStore()
const authStore = useAuthStore()

const quantity = ref(1)
const product = computed(() => productStore.currentProduct)

onMounted(() => {
  productStore.fetchProduct(route.params.id)
})

const handleAddToCart = async () => {
  if (!authStore.isLoggedIn) {
    router.push('/login')
    return
  }
  await cartStore.addToCart(product.value.ID, quantity.value)
  alert(`Added ${quantity.value} to cart!`)
}
</script>

<style scoped>
.product-detail {
  padding: 40px 20px;
  max-width: 900px;
  margin: 0 auto;
}

.detail-content {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 40px;
  align-items: start;
}

.product-hero {
  background: linear-gradient(135deg, #3d2a4a 0%, #5a3a6a 100%);
  border-radius: 20px;
  height: 400px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 8rem;
}

.category-badge {
  background: var(--primary);
  color: #fff;
  padding: 5px 15px;
  border-radius: 20px;
  font-size: 0.85rem;
  display: inline-block;
  margin-bottom: 15px;
}

.product-info h1 {
  font-size: 2.5rem;
  margin-bottom: 15px;
}

.description {
  color: #aaa;
  line-height: 1.7;
  margin-bottom: 25px;
}

.meta {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 30px;
}

.price {
  font-size: 2rem;
  font-weight: bold;
  color: var(--primary);
}

.stock {
  color: #6a6;
  font-size: 0.95rem;
}

.actions {
  display: flex;
  gap: 15px;
  margin-bottom: 30px;
}

.actions input {
  width: 80px;
  padding: 12px;
  border-radius: 8px;
  border: 1px solid #444;
  background: #222;
  color: #fff;
  text-align: center;
}

.back-link {
  color: var(--primary);
  text-decoration: none;
}

.back-link:hover {
  text-decoration: underline;
}

.loading, .not-found {
  text-align: center;
  padding: 100px;
  color: #888;
}

@media (max-width: 768px) {
  .detail-content {
    grid-template-columns: 1fr;
  }
}
</style>
