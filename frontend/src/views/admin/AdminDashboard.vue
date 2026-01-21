<template>
  <div class="dashboard">
    <h1>📊 Admin Dashboard</h1>
    
    <!-- Stats Cards -->
    <div class="stats-grid">
      <div class="stat-card">
        <span class="stat-icon">💰</span>
        <div class="stat-info">
          <span class="stat-value">${{ formatNumber(stats.total_revenue) }}</span>
          <span class="stat-label">Total Revenue</span>
        </div>
      </div>
      <div class="stat-card">
        <span class="stat-icon">📦</span>
        <div class="stat-info">
          <span class="stat-value">{{ stats.total_orders }}</span>
          <span class="stat-label">Orders</span>
        </div>
      </div>
      <div class="stat-card">
        <span class="stat-icon">🍾</span>
        <div class="stat-info">
          <span class="stat-value">{{ stats.total_products }}</span>
          <span class="stat-label">Products</span>
        </div>
      </div>
      <div class="stat-card">
        <span class="stat-icon">👥</span>
        <div class="stat-info">
          <span class="stat-value">{{ stats.total_customers }}</span>
          <span class="stat-label">Customers</span>
        </div>
      </div>
    </div>

    <!-- Charts Row -->
    <div class="charts-row">
      <!-- Sales by Category (Pie Chart) -->
      <div class="chart-card">
        <h3>Sales by Category</h3>
        <div class="chart-container" v-if="categoryData.labels.length">
          <Pie :data="categoryData" :options="pieOptions" />
        </div>
        <p v-else class="no-data">No sales data yet</p>
      </div>

      <!-- Sales Over Time (Line Chart) -->
      <div class="chart-card wide">
        <h3>Sales (Last 30 Days)</h3>
        <div class="chart-container" v-if="salesData.labels.length">
          <Line :data="salesData" :options="lineOptions" />
        </div>
        <p v-else class="no-data">No sales data yet</p>
      </div>
    </div>

    <!-- Top Products & Recent Orders -->
    <div class="tables-row">
      <!-- Top Products -->
      <div class="table-card">
        <h3>🏆 Top Selling Wines</h3>
        <table v-if="topProducts.length">
          <thead>
            <tr>
              <th>Wine</th>
              <th>Sold</th>
              <th>Revenue</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="product in topProducts" :key="product.id">
              <td>{{ product.name }}</td>
              <td>{{ product.quantity }}</td>
              <td>${{ formatNumber(product.revenue) }}</td>
            </tr>
          </tbody>
        </table>
        <p v-else class="no-data">No products sold yet</p>
      </div>

      <!-- Recent Orders -->
      <div class="table-card">
        <h3>🕐 Recent Orders</h3>
        <table v-if="recentOrders.length">
          <thead>
            <tr>
              <th>Order</th>
              <th>Customer</th>
              <th>Total</th>
              <th>Status</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="order in recentOrders" :key="order.id">
              <td>#{{ order.id }}</td>
              <td>{{ order.user_email || 'Guest' }}</td>
              <td>${{ formatNumber(order.total) }}</td>
              <td><span :class="['status', order.status]">{{ order.status }}</span></td>
            </tr>
          </tbody>
        </table>
        <p v-else class="no-data">No orders yet</p>
      </div>
    </div>
    
    <!-- Quick Actions -->
    <div class="quick-actions">
      <h3>Quick Actions</h3>
      <div class="actions-row">
        <router-link to="/admin/products/new" class="action-btn">➕ Add New Wine</router-link>
        <router-link to="/admin/products" class="action-btn">📋 Manage Products</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive } from 'vue'
import { Chart as ChartJS, ArcElement, Tooltip, Legend, CategoryScale, LinearScale, PointElement, LineElement, Title, Filler } from 'chart.js'
import { Pie, Line } from 'vue-chartjs'
import api from '../../services/api'

ChartJS.register(ArcElement, Tooltip, Legend, CategoryScale, LinearScale, PointElement, LineElement, Title, Filler)

const stats = ref({
  total_revenue: 0,
  total_orders: 0,
  total_products: 0,
  total_customers: 0
})

const topProducts = ref([])
const recentOrders = ref([])

const categoryData = reactive({
  labels: [],
  datasets: [{
    data: [],
    backgroundColor: ['#722F37', '#8B4513', '#D4A574', '#A0522D', '#CD853F']
  }]
})

const salesData = reactive({
  labels: [],
  datasets: [{
    label: 'Revenue',
    data: [],
    borderColor: '#722F37',
    backgroundColor: 'rgba(114, 47, 55, 0.2)',
    fill: true,
    tension: 0.4
  }]
})

const pieOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { position: 'bottom', labels: { color: '#ccc' } }
  }
}

const lineOptions = {
  responsive: true,
  maintainAspectRatio: false,
  scales: {
    y: { beginAtZero: true, ticks: { color: '#888' }, grid: { color: 'rgba(255,255,255,0.1)' } },
    x: { ticks: { color: '#888' }, grid: { color: 'rgba(255,255,255,0.1)' } }
  },
  plugins: {
    legend: { labels: { color: '#ccc' } }
  }
}

const formatNumber = (num) => {
  return (num || 0).toFixed(2)
}

onMounted(async () => {
  try {
    // Fetch all analytics data
    const [statsRes, categoryRes, topRes, salesRes, ordersRes] = await Promise.all([
      api.get('/admin/analytics/stats'),
      api.get('/admin/analytics/sales-by-category'),
      api.get('/admin/analytics/top-products?limit=5'),
      api.get('/admin/analytics/sales-by-day?days=30'),
      api.get('/admin/analytics/recent-orders?limit=5')
    ])

    // Stats
    stats.value = statsRes.data.data || stats.value

    // Category Pie Chart
    const categories = categoryRes.data.data || []
    categoryData.labels = categories.map(c => c.category || 'Unknown')
    categoryData.datasets[0].data = categories.map(c => c.revenue)

    // Top Products
    topProducts.value = topRes.data.data || []

    // Sales Line Chart
    const sales = salesRes.data.data || []
    salesData.labels = sales.map(s => s.date)
    salesData.datasets[0].data = sales.map(s => s.revenue)

    // Recent Orders
    recentOrders.value = ordersRes.data.data || []

  } catch (err) {
    console.error('Failed to load analytics:', err)
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
  margin-bottom: 30px;
}

.stat-card {
  background: var(--card-bg);
  border-radius: 16px;
  padding: 25px;
  display: flex;
  align-items: center;
  gap: 20px;
}

.stat-icon { font-size: 2.5rem; }

.stat-info { display: flex; flex-direction: column; }

.stat-value {
  font-size: 1.8rem;
  font-weight: bold;
  color: #fff;
}

.stat-label { color: #888; font-size: 0.9rem; }

/* Charts */
.charts-row {
  display: grid;
  grid-template-columns: 1fr 2fr;
  gap: 20px;
  margin-bottom: 30px;
}

.chart-card {
  background: var(--card-bg);
  border-radius: 16px;
  padding: 20px;
}

.chart-card h3 {
  color: #ccc;
  margin-bottom: 15px;
  font-size: 1rem;
}

.chart-container {
  height: 250px;
}

.no-data {
  color: #666;
  text-align: center;
  padding: 50px;
}

/* Tables */
.tables-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  margin-bottom: 30px;
}

.table-card {
  background: var(--card-bg);
  border-radius: 16px;
  padding: 20px;
}

.table-card h3 {
  color: #ccc;
  margin-bottom: 15px;
  font-size: 1rem;
}

table {
  width: 100%;
  border-collapse: collapse;
}

th, td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid rgba(255,255,255,0.1);
}

th { color: #888; font-weight: 600; }
td { color: #ddd; }

.status {
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 0.8rem;
}

.status.pending { background: #d4a574; color: #222; }
.status.completed { background: #4caf50; color: #fff; }
.status.cancelled { background: #f44336; color: #fff; }

/* Quick Actions */
.quick-actions h3 { color: #ccc; margin-bottom: 15px; }

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

@media (max-width: 768px) {
  .charts-row, .tables-row { grid-template-columns: 1fr; }
}
</style>
