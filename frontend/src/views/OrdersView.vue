<template>
  <div class="orders-page">
    <h1>📦 Order History</h1>
    
    <div v-if="loading" class="loading">Loading orders...</div>
    
    <div v-else-if="orders.length === 0" class="empty">
      <p>No orders yet</p>
      <router-link to="/products" class="btn btn-primary">Start Shopping</router-link>
    </div>
    
    <div v-else class="orders-list">
      <div v-for="order in orders" :key="order.ID" class="order-card">
        <div class="order-header">
          <span class="order-id">Order #{{ order.ID }}</span>
          <span class="order-date">{{ formatDate(order.CreatedAt) }}</span>
          <span class="order-status" :class="order.status.toLowerCase()">{{ order.status }}</span>
        </div>
        
        <div class="order-items">
          <div v-for="item in order.items" :key="item.ID" class="order-item">
            <span class="item-name">{{ item.product?.name || 'Product' }}</span>
            <span class="item-qty">× {{ item.quantity }}</span>
            <span class="item-price">${{ item.price.toFixed(2) }}</span>
          </div>
        </div>
        
        <div class="order-footer">
          <span class="total-label">Total:</span>
          <span class="total-amount">${{ order.total.toFixed(2) }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../services/api'

const orders = ref([])
const loading = ref(true)

onMounted(async () => {
  try {
    const response = await api.get('/orders')
    orders.value = response.data.data || []
  } catch (error) {
    console.error('Failed to fetch orders:', error)
  } finally {
    loading.value = false
  }
})

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}
</script>

<style scoped>
.orders-page {
  padding: 40px 20px;
  max-width: 800px;
  margin: 0 auto;
}

.orders-page h1 {
  text-align: center;
  margin-bottom: 40px;
}

.orders-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.order-card {
  background: var(--card-bg);
  border-radius: 16px;
  overflow: hidden;
}

.order-header {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 20px;
  background: rgba(139, 69, 102, 0.2);
}

.order-id {
  font-weight: bold;
  color: var(--primary);
}

.order-date {
  flex: 1;
  color: #888;
  font-size: 0.9rem;
}

.order-status {
  padding: 5px 12px;
  border-radius: 20px;
  font-size: 0.8rem;
  font-weight: bold;
  text-transform: uppercase;
}

.order-status.paid {
  background: #2a4a2a;
  color: #6f6;
}

.order-status.pending {
  background: #4a4a2a;
  color: #ff6;
}

.order-items {
  padding: 20px;
}

.order-item {
  display: flex;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid #333;
}

.order-item:last-child {
  border-bottom: none;
}

.item-name {
  flex: 1;
}

.item-qty {
  color: #888;
  margin: 0 20px;
}

.item-price {
  font-weight: bold;
  color: var(--primary);
}

.order-footer {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 10px;
  padding: 20px;
  background: rgba(0, 0, 0, 0.2);
}

.total-label {
  color: #888;
}

.total-amount {
  font-size: 1.3rem;
  font-weight: bold;
  color: var(--primary);
}

.loading, .empty {
  text-align: center;
  padding: 60px;
  color: #888;
}

.empty p {
  margin-bottom: 20px;
}
</style>
