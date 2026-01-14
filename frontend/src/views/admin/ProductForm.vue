<template>
  <div class="product-form-page">
    <h1>{{ isEdit ? 'Edit Wine' : 'Add New Wine' }}</h1>
    
    <form @submit.prevent="handleSubmit" class="product-form">
      <div class="form-group">
        <label>Name</label>
        <input v-model="form.name" type="text" placeholder="e.g. Château Margaux" required />
      </div>
      
      <div class="form-group">
        <label>Description</label>
        <textarea v-model="form.description" placeholder="Wine description..." rows="4"></textarea>
      </div>
      
      <div class="form-row">
        <div class="form-group">
          <label>Price ($)</label>
          <input v-model.number="form.price" type="number" step="0.01" min="0" required />
        </div>
        
        <div class="form-group">
          <label>Stock</label>
          <input v-model.number="form.stock" type="number" min="0" required />
        </div>
      </div>
      
      <div class="form-group">
        <label>Category</label>
        <select v-model="form.category">
          <option value="Red">Red</option>
          <option value="White">White</option>
          <option value="Rosé">Rosé</option>
          <option value="Sparkling">Sparkling</option>
          <option value="Dessert">Dessert</option>
        </select>
      </div>
      
      <div class="form-actions">
        <router-link to="/admin/products" class="btn btn-secondary">Cancel</router-link>
        <button type="submit" class="btn btn-primary" :disabled="loading">
          {{ loading ? 'Saving...' : (isEdit ? 'Update Wine' : 'Create Wine') }}
        </button>
      </div>
      
      <p v-if="error" class="error">{{ error }}</p>
    </form>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '../../services/api'

const route = useRoute()
const router = useRouter()

const isEdit = computed(() => !!route.params.id)
const loading = ref(false)
const error = ref('')

const form = ref({
  name: '',
  description: '',
  price: 0,
  stock: 0,
  category: 'Red'
})

onMounted(async () => {
  if (isEdit.value) {
    try {
      const res = await api.get(`/products/${route.params.id}`)
      const product = res.data.data
      form.value = {
        name: product.name,
        description: product.description || '',
        price: product.price,
        stock: product.stock,
        category: product.category || 'Red'
      }
    } catch (err) {
      error.value = 'Failed to load product'
    }
  }
})

const handleSubmit = async () => {
  loading.value = true
  error.value = ''
  
  try {
    if (isEdit.value) {
      await api.put(`/admin/products/${route.params.id}`, form.value)
    } else {
      await api.post('/admin/products', form.value)
    }
    router.push('/admin/products')
  } catch (err) {
    error.value = err.response?.data?.error || 'Failed to save product'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.product-form-page {
  max-width: 600px;
}

.product-form-page h1 {
  color: var(--primary);
  margin-bottom: 30px;
}

.product-form {
  background: var(--card-bg);
  padding: 30px;
  border-radius: 16px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  color: #ccc;
}

.form-group input,
.form-group textarea,
.form-group select {
  width: 100%;
  padding: 12px;
  border-radius: 10px;
  border: 1px solid #444;
  background: #1a1a2e;
  color: #fff;
  font-size: 1rem;
}

.form-group input:focus,
.form-group textarea:focus,
.form-group select:focus {
  outline: none;
  border-color: var(--primary);
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.form-actions {
  display: flex;
  gap: 15px;
  margin-top: 30px;
}

.btn-secondary {
  background: #333;
  color: #fff;
  padding: 12px 25px;
  border-radius: 10px;
  text-decoration: none;
}

.error {
  color: #e55;
  margin-top: 15px;
}
</style>
