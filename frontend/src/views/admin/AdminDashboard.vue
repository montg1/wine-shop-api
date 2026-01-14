<template>
  <div class="dashboard">
    <h1>Dashboard</h1>
    
    <div class="stats-grid">
      <div class="stat-card">
        <span class="stat-icon">🍾</span>
        <div class="stat-info">
          <span class="stat-value">{{ stats.totalProducts }}</span>
          <span class="stat-label">Products</span>
        </div>
      </div>
      <div class="stat-card">
        <span class="stat-icon">📦</span>
        <div class="stat-info">
          <span class="stat-value">{{ stats.totalOrders }}</span>
          <span class="stat-label">Orders</span>
        </div>
      </div>
      <div class="stat-card">
        <span class="stat-icon">💰</span>
        <div class="stat-info">
          <span class="stat-value">${{ stats.totalRevenue.toFixed(2) }}</span>
          <span class="stat-label">Revenue</span>
        </div>
      </div>
    </div>
    
    <div class="quick-actions">
      <h2>Quick Actions</h2>
      <div class="actions-row">
        <router-link to="/admin/products/new" class="action-btn">
          ➕ Add New Wine
        </router-link>
        <router-link to="/admin/products" class="action-btn">
          📋 Manage Products
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../../services/api'

const stats = ref({
  totalProducts: 0,
  totalOrders: 0,
  totalRevenue: 0
})

onMounted(async () => {
  try {
    const productsRes = await api.get('/products')
    stats.value.totalProducts = productsRes.data.meta?.total || productsRes.data.data?.length || 0
    
    const ordersRes = await api.get('/orders')
    const orders = ordersRes.data.data || []
    stats.value.totalOrders = orders.length
    stats.value.totalRevenue = orders.reduce((sum, o) => sum + (o.total || 0), 0)
  } catch (err) {
    console.error('Failed to load stats:', err)
  }
})
</script>

<style scoped>
.dashboard h1 {
  margin-bottom: 30px;
  color: var(--primary);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  margin-bottom: 40px;
}

.stat-card {
  background: var(--card-bg);
  border-radius: 16px;
  padding: 25px;
  display: flex;
  align-items: center;
  gap: 20px;
}

.stat-icon {
  font-size: 2.5rem;
}

.stat-info {
  display: flex;
  flex-direction: column;
}

.stat-value {
  font-size: 1.8rem;
  font-weight: bold;
  color: #fff;
}

.stat-label {
  color: #888;
  font-size: 0.9rem;
}

.quick-actions h2 {
  margin-bottom: 20px;
  color: #ccc;
}

.actions-row {
  display: flex;
  gap: 15px;
  flex-wrap: wrap;
}

.action-btn {
  background: var(--card-bg);
  padding: 15px 25px;
  border-radius: 12px;
  color: var(--primary);
  text-decoration: none;
  transition: all 0.3s;
}

.action-btn:hover {
  background: rgba(139, 69, 102, 0.3);
  transform: translateY(-2px);
}
</style>
