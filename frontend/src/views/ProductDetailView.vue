<template>
  <div class="product-detail">
    <div v-if="productStore.loading" class="loading">Loading...</div>
    
    <div v-else-if="product" class="detail-content">
      <div class="product-hero">
        <img :src="getWineImage(product.name)" :alt="product.name" />
      </div>
      
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

// Import all wine images
import pinotNoirImg from '../assets/images/pinot-noir.png'
import cabernetImg from '../assets/images/cabernet.png'
import merlotImg from '../assets/images/merlot.png'
import chardonnayImg from '../assets/images/chardonnay.png'
import sauvignonBlancImg from '../assets/images/sauvignon-blanc.png'
import roseImg from '../assets/images/rose.png'
import defaultWineImg from '../assets/images/wine-bottle.png'

const wineImages = {
  'pinot noir': pinotNoirImg,
  'cabernet sauvignon': cabernetImg,
  'cabernet': cabernetImg,
  'merlot': merlotImg,
  'chardonnay': chardonnayImg,
  'sauvignon blanc': sauvignonBlancImg,
  'rosé': roseImg,
  'rose': roseImg
}

const getWineImage = (name) => {
  const lowerName = name.toLowerCase()
  return wineImages[lowerName] || defaultWineImg
}

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
  padding: 60px 20px;
  max-width: 1000px;
  margin: 0 auto;
}

.detail-content {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 60px;
  align-items: start;
}

.product-hero {
  background: var(--bg-warm);
  border-radius: 12px;
  height: 450px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
}

.product-hero img {
  max-height: 100%;
  max-width: 100%;
  object-fit: contain;
}

.category-badge {
  background: var(--primary);
  color: #fff;
  padding: 6px 16px;
  border-radius: 20px;
  font-size: 0.8rem;
  text-transform: uppercase;
  letter-spacing: 1px;
  display: inline-block;
  margin-bottom: 15px;
}

.product-info h1 {
  font-size: 2.5rem;
  color: var(--text);
  margin-bottom: 20px;
}

.description {
  color: var(--text-muted);
  line-height: 1.8;
  margin-bottom: 30px;
  font-size: 1.05rem;
}

.meta {
  display: flex;
  align-items: center;
  gap: 25px;
  margin-bottom: 35px;
}

.price {
  font-size: 2.2rem;
  font-weight: bold;
  color: var(--primary);
}

.stock {
  color: #5a8a5a;
  font-size: 0.95rem;
}

.actions {
  display: flex;
  gap: 15px;
  margin-bottom: 35px;
}

.actions input {
  width: 80px;
  padding: 14px;
  border-radius: 8px;
  border: 1px solid var(--border);
  background: var(--bg);
  color: var(--text);
  text-align: center;
  font-size: 1rem;
}

.back-link {
  color: var(--primary);
  text-decoration: none;
  font-size: 0.95rem;
}

.back-link:hover {
  text-decoration: underline;
}

.loading, .not-found {
  text-align: center;
  padding: 100px;
  color: var(--text-muted);
}

@media (max-width: 768px) {
  .detail-content {
    grid-template-columns: 1fr;
  }
  .product-hero {
    height: 300px;
  }
}
</style>
