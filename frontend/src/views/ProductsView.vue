<template>
  <div class="products-page">
    <div class="page-header">
      <h1>Our Collection</h1>
      <p>Carefully selected wines from exceptional vineyards</p>
    </div>
    
    <!-- Search and Filter -->
    <div class="filters">
      <div class="search-box">
        <input 
          v-model="searchQuery" 
          type="text" 
          placeholder="Search wines..." 
          @input="handleSearch"
        />
      </div>
      <div class="category-filters">
        <button 
          :class="{ active: selectedCategory === '' }" 
          @click="filterByCategory('')"
        >
          All
        </button>
        <button 
          :class="{ active: selectedCategory === 'Red' }" 
          @click="filterByCategory('Red')"
        >
          Red
        </button>
        <button 
          :class="{ active: selectedCategory === 'White' }" 
          @click="filterByCategory('White')"
        >
          White
        </button>
        <button 
          :class="{ active: selectedCategory === 'Rosé' }" 
          @click="filterByCategory('Rosé')"
        >
          Rosé
        </button>
      </div>
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
          <img :src="getWineImage(product.name)" :alt="product.name" />
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
      <p>No wines found{{ searchQuery ? ' for "' + searchQuery + '"' : '' }}.</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useProductStore } from '../stores/products'
import { useCartStore } from '../stores/cart'
import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'

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

const productStore = useProductStore()
const cartStore = useCartStore()
const authStore = useAuthStore()
const router = useRouter()

const searchQuery = ref('')
const selectedCategory = ref('')
let searchTimeout = null

onMounted(() => {
  productStore.fetchProducts()
})

const handleSearch = () => {
  // Debounce search
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    productStore.fetchProducts(1, 50, searchQuery.value, selectedCategory.value)
  }, 300)
}

const filterByCategory = (category) => {
  selectedCategory.value = category
  productStore.fetchProducts(1, 50, searchQuery.value, category)
}

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
  margin-bottom: 40px;
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

/* Filters */
.filters {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
  justify-content: center;
  align-items: center;
  margin-bottom: 50px;
}

.search-box input {
  padding: 14px 20px;
  border-radius: 30px;
  border: 1px solid var(--border);
  background: var(--card-bg);
  color: var(--text);
  font-size: 1rem;
  width: 280px;
  transition: border-color 0.3s;
}

.search-box input:focus {
  outline: none;
  border-color: var(--primary);
}

.category-filters {
  display: flex;
  gap: 10px;
}

.category-filters button {
  padding: 10px 20px;
  border-radius: 25px;
  border: 1px solid var(--border);
  background: var(--card-bg);
  color: var(--text-muted);
  cursor: pointer;
  font-size: 0.9rem;
  transition: all 0.3s;
}

.category-filters button:hover {
  border-color: var(--primary);
  color: var(--primary);
}

.category-filters button.active {
  background: var(--primary);
  border-color: var(--primary);
  color: #fff;
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

@media (max-width: 600px) {
  .filters {
    flex-direction: column;
  }
  .search-box input {
    width: 100%;
  }
}
</style>
