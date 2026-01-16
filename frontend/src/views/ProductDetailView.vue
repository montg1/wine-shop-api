<template>
  <div class="product-detail">
    <div v-if="productStore.loading" class="loading">Loading...</div>
    
    <div v-else-if="product" class="detail-content">
      <div class="product-hero">
        <img :src="getWineImage(product)" :alt="product.name" />
      </div>
      
      <div class="product-info">
        <span class="category-badge">{{ product.category }}</span>
        <h1>{{ product.name }}</h1>
        
        <!-- Rating Summary -->
        <div class="rating-summary" v-if="reviewMeta.total_reviews > 0">
          <div class="stars">
            <span v-for="i in 5" :key="i" class="star" :class="{ filled: i <= Math.round(reviewMeta.average_rating) }">★</span>
          </div>
          <span class="rating-text">{{ reviewMeta.average_rating.toFixed(1) }} ({{ reviewMeta.total_reviews }} reviews)</span>
        </div>
        
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
    
    <!-- Reviews Section -->
    <div v-if="product" class="reviews-section">
      <h2>Customer Reviews</h2>
      
      <!-- Add Review Form -->
      <div v-if="authStore.isLoggedIn" class="add-review">
        <h3>Write a Review</h3>
        <div class="star-input">
          <span 
            v-for="i in 5" 
            :key="i" 
            class="star clickable"
            :class="{ filled: i <= newReview.rating }"
            @click="newReview.rating = i"
          >★</span>
        </div>
        <textarea v-model="newReview.comment" placeholder="Share your experience..." rows="3"></textarea>
        <button class="btn btn-primary" @click="submitReview" :disabled="!newReview.rating">Submit Review</button>
      </div>
      <div v-else class="login-prompt">
        <router-link to="/login">Login</router-link> to write a review
      </div>
      
      <!-- Reviews List -->
      <div class="reviews-list">
        <div v-if="reviews.length === 0" class="no-reviews">No reviews yet. Be the first!</div>
        <div v-for="review in reviews" :key="review.ID" class="review-card">
          <div class="review-header">
            <div class="stars">
              <span v-for="i in 5" :key="i" class="star small" :class="{ filled: i <= review.rating }">★</span>
            </div>
            <span class="review-author">{{ review.user?.email || 'Anonymous' }}</span>
            <span class="review-date">{{ formatDate(review.CreatedAt) }}</span>
          </div>
          <p class="review-comment">{{ review.comment || 'No comment' }}</p>
        </div>
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
import api from '../services/api'

// Import all wine images (fallback for products without Cloudinary URL)
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

// Use Cloudinary URL if available, fallback to local mapping
const getWineImage = (product) => {
  if (product?.image_url) {
    return product.image_url
  }
  const lowerName = (product?.name || '').toLowerCase()
  return wineImages[lowerName] || defaultWineImg
}

const route = useRoute()
const router = useRouter()
const productStore = useProductStore()
const cartStore = useCartStore()
const authStore = useAuthStore()

const quantity = ref(1)
const product = computed(() => productStore.currentProduct)
const reviews = ref([])
const reviewMeta = ref({ average_rating: 0, total_reviews: 0 })
const newReview = ref({ rating: 0, comment: '' })

onMounted(async () => {
  await productStore.fetchProduct(route.params.id)
  await fetchReviews()
})

const fetchReviews = async () => {
  try {
    const response = await api.get(`/products/${route.params.id}/reviews`)
    reviews.value = response.data.data || []
    reviewMeta.value = response.data.meta
  } catch (error) {
    console.error('Failed to fetch reviews:', error)
  }
}

const submitReview = async () => {
  if (!newReview.value.rating) return
  
  try {
    await api.post(`/products/${route.params.id}/reviews`, {
      rating: newReview.value.rating,
      comment: newReview.value.comment
    })
    newReview.value = { rating: 0, comment: '' }
    await fetchReviews()
  } catch (error) {
    alert(error.response?.data?.error || 'Failed to submit review')
  }
}

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

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
  margin-bottom: 10px;
}

.rating-summary {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 20px;
}

.rating-text {
  color: var(--text-muted);
  font-size: 0.95rem;
}

.stars {
  display: flex;
  gap: 2px;
}

.star {
  color: #ddd;
  font-size: 1.2rem;
}

.star.filled {
  color: #f4c150;
}

.star.small {
  font-size: 1rem;
}

.star.clickable {
  cursor: pointer;
  font-size: 1.8rem;
  transition: transform 0.2s;
}

.star.clickable:hover {
  transform: scale(1.2);
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

/* Reviews Section */
.reviews-section {
  margin-top: 60px;
  padding-top: 40px;
  border-top: 1px solid var(--border);
}

.reviews-section h2 {
  font-size: 1.8rem;
  color: var(--text);
  margin-bottom: 30px;
}

.add-review {
  background: var(--card-bg);
  padding: 25px;
  border-radius: 12px;
  margin-bottom: 30px;
}

.add-review h3 {
  margin-bottom: 15px;
  font-size: 1.1rem;
}

.star-input {
  margin-bottom: 15px;
}

.add-review textarea {
  width: 100%;
  padding: 12px;
  border-radius: 8px;
  border: 1px solid var(--border);
  background: var(--bg);
  color: var(--text);
  resize: vertical;
  margin-bottom: 15px;
  font-family: inherit;
}

.login-prompt {
  background: var(--card-bg);
  padding: 20px;
  border-radius: 12px;
  text-align: center;
  margin-bottom: 30px;
  color: var(--text-muted);
}

.login-prompt a {
  color: var(--primary);
}

.reviews-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.no-reviews {
  text-align: center;
  padding: 40px;
  color: var(--text-muted);
}

.review-card {
  background: var(--card-bg);
  padding: 20px;
  border-radius: 12px;
}

.review-header {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 10px;
}

.review-author {
  font-weight: 500;
  color: var(--text);
}

.review-date {
  color: var(--text-muted);
  font-size: 0.85rem;
  margin-left: auto;
}

.review-comment {
  color: var(--text-muted);
  line-height: 1.6;
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
