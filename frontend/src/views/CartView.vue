<template>
  <div class="cart-page">
    <h1>🛒 Your Cart</h1>
    
    <div v-if="cartStore.loading" class="loading">Loading cart...</div>
    
    <div v-else-if="cartStore.items.length === 0" class="empty-cart">
      <p>Your cart is empty</p>
      <router-link to="/products" class="btn btn-primary">Browse Wines</router-link>
    </div>
    
    <div v-else class="cart-content">
      <div class="cart-items">
        <div v-for="item in cartStore.items" :key="item.ID" class="cart-item">
          <div class="item-image">🍷</div>
          <div class="item-details">
            <h3>{{ item.product.name }}</h3>
            <p class="item-price">${{ item.product.price.toFixed(2) }} × {{ item.quantity }}</p>
          </div>
          <div class="item-total">
            ${{ (item.product.price * item.quantity).toFixed(2) }}
          </div>
        </div>
      </div>
      
      <div class="cart-summary">
        <div class="summary-row">
          <span>Items:</span>
          <span>{{ cartStore.totalItems }}</span>
        </div>
        <div class="summary-row total">
          <span>Total:</span>
          <span>${{ cartStore.totalPrice.toFixed(2) }}</span>
        </div>
        <button class="btn btn-primary btn-block" @click="handleCheckout" :disabled="checkingOut">
          {{ checkingOut ? 'Processing...' : 'Checkout' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useCartStore } from '../stores/cart'

const router = useRouter()
const cartStore = useCartStore()
const checkingOut = ref(false)

onMounted(() => {
  cartStore.fetchCart()
})

const handleCheckout = async () => {
  checkingOut.value = true
  try {
    await cartStore.checkout()
    alert('Order placed successfully!')
    router.push('/orders')
  } catch (error) {
    alert('Checkout failed: ' + (error.response?.data?.error || 'Unknown error'))
  } finally {
    checkingOut.value = false
  }
}
</script>

<style scoped>
.cart-page {
  padding: 40px 20px;
  max-width: 900px;
  margin: 0 auto;
}

.cart-page h1 {
  text-align: center;
  margin-bottom: 40px;
}

.empty-cart {
  text-align: center;
  padding: 60px;
}

.empty-cart p {
  font-size: 1.2rem;
  color: #888;
  margin-bottom: 20px;
}

.cart-content {
  display: grid;
  grid-template-columns: 1fr 300px;
  gap: 30px;
}

.cart-items {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.cart-item {
  background: var(--card-bg);
  border-radius: 12px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 20px;
}

.item-image {
  font-size: 2.5rem;
  width: 60px;
  text-align: center;
}

.item-details {
  flex: 1;
}

.item-details h3 {
  margin-bottom: 5px;
}

.item-price {
  color: #888;
  font-size: 0.9rem;
}

.item-total {
  font-size: 1.2rem;
  font-weight: bold;
  color: var(--primary);
}

.cart-summary {
  background: var(--card-bg);
  border-radius: 16px;
  padding: 25px;
  height: fit-content;
  position: sticky;
  top: 100px;
}

.summary-row {
  display: flex;
  justify-content: space-between;
  padding: 10px 0;
  color: #aaa;
}

.summary-row.total {
  font-size: 1.4rem;
  font-weight: bold;
  color: #fff;
  border-top: 1px solid #444;
  margin-top: 10px;
  padding-top: 20px;
}

.btn-block {
  width: 100%;
  margin-top: 20px;
}

.loading {
  text-align: center;
  padding: 60px;
  color: #888;
}

@media (max-width: 768px) {
  .cart-content {
    grid-template-columns: 1fr;
  }
}
</style>
