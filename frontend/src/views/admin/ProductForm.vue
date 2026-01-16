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

      <!-- Image Upload -->
      <div class="form-group">
        <label>Product Image</label>
        <div class="image-upload-area">
          <div v-if="form.image_url || imagePreview" class="image-preview">
            <img :src="imagePreview || form.image_url" alt="Preview" />
            <button type="button" class="remove-btn" @click="removeImage">×</button>
          </div>
          <div v-else class="upload-placeholder">
            <input 
              type="file" 
              accept="image/*" 
              @change="handleImageSelect" 
              ref="fileInput"
              id="image-input"
            />
            <label for="image-input" class="upload-label">
              📷 Click to upload image
            </label>
          </div>
        </div>
        <p v-if="uploading" class="upload-status">Uploading image...</p>
      </div>
      
      <div class="form-actions">
        <router-link to="/admin/products" class="btn btn-secondary">Cancel</router-link>
        <button type="submit" class="btn btn-primary" :disabled="loading || uploading">
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
const uploading = ref(false)
const error = ref('')
const fileInput = ref(null)
const imagePreview = ref(null)

const form = ref({
  name: '',
  description: '',
  price: 0,
  stock: 0,
  category: 'Red',
  image_url: ''
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
        category: product.category || 'Red',
        image_url: product.image_url || ''
      }
    } catch (err) {
      error.value = 'Failed to load product'
    }
  }
})

const handleImageSelect = async (event) => {
  const file = event.target.files[0]
  if (!file) return

  // Show preview
  const reader = new FileReader()
  reader.onload = (e) => {
    imagePreview.value = e.target.result
  }
  reader.readAsDataURL(file)

  // Upload to Cloudinary
  uploading.value = true
  try {
    const formData = new FormData()
    formData.append('file', file)
    
    const res = await api.post('/admin/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    
    form.value.image_url = res.data.url
    imagePreview.value = null
  } catch (err) {
    error.value = 'Failed to upload image'
    imagePreview.value = null
  } finally {
    uploading.value = false
  }
}

const removeImage = () => {
  form.value.image_url = ''
  imagePreview.value = null
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}

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
  color: var(--text-muted);
}

.form-group input,
.form-group textarea,
.form-group select {
  width: 100%;
  padding: 12px;
  border-radius: 10px;
  border: 1px solid var(--border);
  background: var(--bg);
  color: var(--text);
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

/* Image Upload Styles */
.image-upload-area {
  border: 2px dashed var(--border);
  border-radius: 12px;
  padding: 20px;
  text-align: center;
  position: relative;
}

.upload-placeholder input[type="file"] {
  position: absolute;
  opacity: 0;
  width: 100%;
  height: 100%;
  top: 0;
  left: 0;
  cursor: pointer;
}

.upload-label {
  display: block;
  padding: 40px;
  cursor: pointer;
  color: var(--text-muted);
  font-size: 1.1rem;
}

.upload-label:hover {
  color: var(--primary);
}

.image-preview {
  position: relative;
  display: inline-block;
}

.image-preview img {
  max-width: 200px;
  max-height: 200px;
  border-radius: 8px;
  object-fit: cover;
}

.remove-btn {
  position: absolute;
  top: -10px;
  right: -10px;
  width: 30px;
  height: 30px;
  border-radius: 50%;
  background: #e55;
  color: white;
  border: none;
  cursor: pointer;
  font-size: 1.2rem;
  line-height: 1;
}

.upload-status {
  color: var(--primary);
  margin-top: 10px;
  font-size: 0.9rem;
}

.form-actions {
  display: flex;
  gap: 15px;
  margin-top: 30px;
}

.btn-secondary {
  background: var(--border);
  color: var(--text);
  padding: 12px 25px;
  border-radius: 10px;
  text-decoration: none;
}

.error {
  color: #e55;
  margin-top: 15px;
}
</style>
