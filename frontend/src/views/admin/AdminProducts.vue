<template>
  <div class="admin-products">
    <div class="page-header">
      <h1>Products</h1>
      <router-link to="/admin/products/new" class="btn btn-primary">
        ➕ Add Wine
      </router-link>
    </div>
    
    <div v-if="loading" class="loading">Loading products...</div>
    
    <table v-else class="products-table">
      <thead>
        <tr>
          <th>ID</th>
          <th>Name</th>
          <th>Category</th>
          <th>Price</th>
          <th>Stock</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="product in products" :key="product.ID">
          <td>{{ product.ID }}</td>
          <td>{{ product.name }}</td>
          <td>{{ product.category }}</td>
          <td>${{ product.price.toFixed(2) }}</td>
          <td>{{ product.stock }}</td>
          <td class="actions">
            <router-link :to="`/admin/products/${product.ID}/edit`" class="btn-icon edit">
              ✏️
            </router-link>
            <button @click="deleteProduct(product.ID)" class="btn-icon delete">
              🗑️
            </button>
          </td>
        </tr>
      </tbody>
    </table>
    
    <div v-if="!loading && products.length === 0" class="empty">
      No products yet. <router-link to="/admin/products/new">Add one!</router-link>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../../services/api'

const products = ref([])
const loading = ref(true)

const fetchProducts = async () => {
  loading.value = true
  try {
    const res = await api.get('/products?limit=100')
    products.value = res.data.data || []
  } catch (err) {
    console.error('Failed to fetch products:', err)
  } finally {
    loading.value = false
  }
}

const deleteProduct = async (id) => {
  if (!confirm('Are you sure you want to delete this product?')) return
  
  try {
    await api.delete(`/admin/products/${id}`)
    products.value = products.value.filter(p => p.ID !== id)
  } catch (err) {
    alert('Failed to delete product: ' + (err.response?.data?.error || 'Unknown error'))
  }
}

onMounted(fetchProducts)
</script>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.page-header h1 {
  color: var(--primary);
  margin: 0;
}

.products-table {
  width: 100%;
  border-collapse: collapse;
  background: var(--card-bg);
  border-radius: 12px;
  overflow: hidden;
}

.products-table th,
.products-table td {
  padding: 15px;
  text-align: left;
  border-bottom: 1px solid #333;
}

.products-table th {
  background: rgba(139, 69, 102, 0.2);
  color: var(--primary);
  font-weight: 600;
}

.products-table tr:hover {
  background: rgba(139, 69, 102, 0.1);
}

.actions {
  display: flex;
  gap: 10px;
}

.btn-icon {
  background: none;
  border: none;
  font-size: 1.2rem;
  cursor: pointer;
  padding: 5px 10px;
  border-radius: 8px;
  transition: background 0.3s;
  text-decoration: none;
}

.btn-icon.edit:hover {
  background: rgba(100, 100, 255, 0.2);
}

.btn-icon.delete:hover {
  background: rgba(255, 100, 100, 0.2);
}

.loading, .empty {
  text-align: center;
  padding: 60px;
  color: #888;
}

.empty a {
  color: var(--primary);
}
</style>
